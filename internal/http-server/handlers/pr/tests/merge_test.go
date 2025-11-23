package pr_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	prhandler "github.com/Dokhoyan/avito-pr-test/internal/http-server/handlers/pr"
	"github.com/Dokhoyan/avito-pr-test/internal/http-server/handlers/testutil"
	"github.com/Dokhoyan/avito-pr-test/internal/model"
	"github.com/Dokhoyan/avito-pr-test/internal/service"
	servicemocks "github.com/Dokhoyan/avito-pr-test/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	testutil.InitTestLogger()
}

func TestMergePR_Success(t *testing.T) {
	prID := "pr1"

	expectedPR := &model.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   "Test PR",
		AuthorID:          "author1",
		Status:            model.PRStatusMerged,
		AssignedReviewers: []string{"reviewer1"},
	}

	mockService := servicemocks.NewMockPRService(t)
	mockService.EXPECT().MergePR(mock.Anything, prID).Return(expectedPR, nil)

	handler := prhandler.NewImplementation(mockService)

	reqBody := map[string]interface{}{
		"pull_request_id": prID,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Merge(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	mockService.AssertExpectations(t)
}

func TestMergePR_InvalidMethod(t *testing.T) {
	mockService := servicemocks.NewMockPRService(t)
	handler := prhandler.NewImplementation(mockService)

	req := httptest.NewRequest(http.MethodGet, "/pullRequest/merge", nil)
	w := httptest.NewRecorder()

	handler.Merge(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
}

func TestMergePR_InvalidPayload(t *testing.T) {
	mockService := servicemocks.NewMockPRService(t)
	handler := prhandler.NewImplementation(mockService)

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	handler.Merge(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestMergePR_NotFound(t *testing.T) {
	prID := "non-existent-pr"

	mockService := servicemocks.NewMockPRService(t)
	mockService.EXPECT().MergePR(mock.Anything, prID).Return(nil, service.ErrPRNotFound)

	handler := prhandler.NewImplementation(mockService)

	reqBody := map[string]interface{}{
		"pull_request_id": prID,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Merge(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	mockService.AssertExpectations(t)
}
