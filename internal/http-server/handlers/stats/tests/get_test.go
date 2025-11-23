package stats_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	statshandler "github.com/Dokhoyan/avito-pr-test/internal/http-server/handlers/stats"
	"github.com/Dokhoyan/avito-pr-test/internal/http-server/handlers/testutil"
	"github.com/Dokhoyan/avito-pr-test/internal/model"
	servicemocks "github.com/Dokhoyan/avito-pr-test/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	testutil.InitTestLogger()
}

func TestGetStats_Success(t *testing.T) {
	expectedStats := &model.Stats{
		TotalTeams:       2,
		TotalUsers:       5,
		ActiveUsers:      4,
		TotalPRs:         7,
		OpenPRs:          5,
		MergedPRs:        2,
		TotalAssignments: 9,
	}

	mockService := servicemocks.NewMockStatsService(t)
	mockService.EXPECT().GetStats(mock.Anything).Return(expectedStats, nil)

	handler := statshandler.NewImplementation(mockService)

	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	w := httptest.NewRecorder()

	handler.GetStats(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	mockService.AssertExpectations(t)
}

func TestGetStats_InvalidMethod(t *testing.T) {
	mockService := servicemocks.NewMockStatsService(t)
	handler := statshandler.NewImplementation(mockService)

	req := httptest.NewRequest(http.MethodPost, "/stats", nil)
	w := httptest.NewRecorder()

	handler.GetStats(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
}
