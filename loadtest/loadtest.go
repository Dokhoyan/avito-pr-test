package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func main() {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	var targets []vegeta.Target

	rand.Seed(time.Now().UnixNano())
	teamName := fmt.Sprintf("team_%d", rand.Intn(1000000))
	userID1 := fmt.Sprintf("user_%d", rand.Intn(1000000))
	userID2 := fmt.Sprintf("user_%d", rand.Intn(1000000))
	userID3 := fmt.Sprintf("user_%d", rand.Intn(1000000))

	teamPayload := map[string]interface{}{
		"team_name": teamName,
		"members": []map[string]interface{}{
			{"user_id": userID1, "username": userID1, "is_active": true},
			{"user_id": userID2, "username": userID2, "is_active": true},
			{"user_id": userID3, "username": userID3, "is_active": true},
		},
	}
	teamBody, _ := json.Marshal(teamPayload)
	targets = append(targets, vegeta.Target{
		Method: "POST",
		URL:    fmt.Sprintf("%s/team/add", baseURL),
		Body:   teamBody,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
	})

	targets = append(targets, vegeta.Target{
		Method: "GET",
		URL:    fmt.Sprintf("%s/team/get?team_name=%s", baseURL, teamName),
	})

	userStatusPayload := map[string]interface{}{
		"user_id":   userID1,
		"is_active": true,
	}
	userStatusBody, _ := json.Marshal(userStatusPayload)
	targets = append(targets, vegeta.Target{
		Method: "POST",
		URL:    fmt.Sprintf("%s/users/setIsActive", baseURL),
		Body:   userStatusBody,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
	})

	prID := fmt.Sprintf("pr_%d", rand.Intn(1000000))
	prPayload := map[string]interface{}{
		"pull_request_id":   prID,
		"pull_request_name": fmt.Sprintf("Test PR %d", rand.Intn(1000000)),
		"author_id":         userID1,
	}
	prBody, _ := json.Marshal(prPayload)
	targets = append(targets, vegeta.Target{
		Method: "POST",
		URL:    fmt.Sprintf("%s/pullRequest/create", baseURL),
		Body:   prBody,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
	})

	targets = append(targets, vegeta.Target{
		Method: "GET",
		URL:    fmt.Sprintf("%s/stats", baseURL),
	})

	rate := vegeta.Rate{Freq: 5, Per: time.Second}
	duration := 5 * time.Minute
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	targeter := vegeta.NewStaticTargeter(targets...)
	for res := range attacker.Attack(targeter, rate, duration, "Load Test") {
		metrics.Add(res)
	}
	metrics.Close()

	generateReports(&metrics)
}

func generateReports(metrics *vegeta.Metrics) {
	sliResponseTimeOK := metrics.Latencies.P95.Seconds()*1000 < 300
	sliSuccessRateOK := (1-metrics.Success)*100 >= 99.9

	jsonReport := map[string]interface{}{
		"total_requests":    metrics.Requests,
		"request_rate":      metrics.Rate,
		"success_rate":      (1 - metrics.Success) * 100,
		"error_rate":        metrics.Success * 100,
		"avg_response_time": metrics.Latencies.Mean.Seconds() * 1000,
		"p95":               metrics.Latencies.P95.Seconds() * 1000,
		"p99":               metrics.Latencies.P99.Seconds() * 1000,
		"min":               metrics.Latencies.Min.Seconds() * 1000,
		"max":               metrics.Latencies.Max.Seconds() * 1000,
		"status_codes":      metrics.StatusCodes,
		"sli": map[string]interface{}{
			"response_time_300ms": map[string]interface{}{
				"threshold": 300,
				"actual":    metrics.Latencies.P95.Seconds() * 1000,
				"ok":        sliResponseTimeOK,
			},
			"success_rate_99_9": map[string]interface{}{
				"threshold": 99.9,
				"actual":    (1 - metrics.Success) * 100,
				"ok":        sliSuccessRateOK,
			},
		},
	}

	jsonData, _ := json.MarshalIndent(jsonReport, "", "  ")
	os.WriteFile("loadtest/summary.json", jsonData, 0644)

	htmlReport := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
  <title>Load Test Report</title>
  <style>
    body { font-family: Arial, sans-serif; margin: 20px; background: #f5f5f5; }
    .container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; }
    h1 { color: #333; }
    .metric { margin: 15px 0; padding: 10px; background: #f9f9f9; border-radius: 4px; }
    .metric-name { font-weight: bold; color: #555; }
    .metric-value { font-size: 18px; color: #2196F3; }
    table { width: 100%%; border-collapse: collapse; margin: 20px 0; }
    th, td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
    th { background: #2196F3; color: white; }
    .success { color: #4CAF50; }
    .error { color: #f44336; }
    .status-pass { color: #4CAF50; font-weight: bold; }
    .status-fail { color: #f44336; font-weight: bold; }
  </style>
</head>
<body>
  <div class="container">
    <h1>Load Test Report</h1>
    
    <div class="metric">
      <div class="metric-name">Total Requests</div>
      <div class="metric-value">%d</div>
    </div>
    
    <div class="metric">
      <div class="metric-name">Request Rate</div>
      <div class="metric-value">%.2f req/s</div>
    </div>
    
    <div class="metric">
      <div class="metric-name">Success Rate</div>
      <div class="metric-value success">%.2f%%</div>
    </div>
    
    <div class="metric">
      <div class="metric-name">Error Rate</div>
      <div class="metric-value error">%.2f%%</div>
    </div>
    
    <h2>SLI Compliance</h2>
    <table>
      <tr>
        <th>Metric</th>
        <th>Threshold</th>
        <th>Actual</th>
        <th>Status</th>
      </tr>
      <tr>
        <td>Response Time (P95)</td>
        <td>&lt; 300 ms</td>
        <td>%.2f ms</td>
        <td class="%s">%s</td>
      </tr>
      <tr>
        <td>Success Rate</td>
        <td>&gt;= 99.9%%</td>
        <td>%.2f%%</td>
        <td class="%s">%s</td>
      </tr>
    </table>
    
    <h2>Response Times</h2>
    <table>
      <tr>
        <th>Metric</th>
        <th>Value (ms)</th>
      </tr>
      <tr>
        <td>Average</td>
        <td>%.2f</td>
      </tr>
      <tr>
        <td>Min</td>
        <td>%.2f</td>
      </tr>
      <tr>
        <td>Max</td>
        <td>%.2f</td>
      </tr>
      <tr>
        <td>P(95)</td>
        <td>%.2f</td>
      </tr>
      <tr>
        <td>P(99)</td>
        <td>%.2f</td>
      </tr>
    </table>
    
    <h2>Status Codes</h2>
    <table>
      <tr>
        <th>Status Code</th>
        <th>Count</th>
      </tr>`,
		metrics.Requests,
		metrics.Rate,
		(1-metrics.Success)*100,
		metrics.Success*100,
		metrics.Latencies.P95.Seconds()*1000,
		getStatusClass(sliResponseTimeOK),
		getStatus(sliResponseTimeOK),
		(1-metrics.Success)*100,
		getStatusClass(sliSuccessRateOK),
		getStatus(sliSuccessRateOK),
		metrics.Latencies.Mean.Seconds()*1000,
		metrics.Latencies.Min.Seconds()*1000,
		metrics.Latencies.Max.Seconds()*1000,
		metrics.Latencies.P95.Seconds()*1000,
		metrics.Latencies.P99.Seconds()*1000,
	)

	for code, count := range metrics.StatusCodes {
		htmlReport += fmt.Sprintf(`
      <tr>
        <td>%s</td>
        <td>%d</td>
      </tr>`, code, count)
	}

	htmlReport += `
    </table>
  </div>
</body>
</html>`

	os.WriteFile("loadtest/summary.html", []byte(htmlReport), 0644)

	fmt.Printf(`
Load Test Summary
=================
Total Requests: %d
Request Rate: %.2f req/s
Success Rate: %.2f%%
Error Rate: %.2f%%
Avg Response Time: %.2f ms
P(95): %.2f ms
P(99): %.2f ms

SLI Compliance:
  Response Time (P95 < 300ms): %s (%.2f ms)
  Success Rate (>= 99.9%%): %s (%.2f%%)
`,
		metrics.Requests,
		metrics.Rate,
		(1-metrics.Success)*100,
		metrics.Success*100,
		metrics.Latencies.Mean.Seconds()*1000,
		metrics.Latencies.P95.Seconds()*1000,
		metrics.Latencies.P99.Seconds()*1000,
		getStatus(sliResponseTimeOK),
		metrics.Latencies.P95.Seconds()*1000,
		getStatus(sliSuccessRateOK),
		(1-metrics.Success)*100,
	)
}

func getStatus(ok bool) string {
	if ok {
		return "✓ PASS"
	}
	return "✗ FAIL"
}

func getStatusClass(ok bool) string {
	if ok {
		return "status-pass"
	}
	return "status-fail"
}
