package pr_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	prhandler "github.com/Dokhoyan/avito-pr-test/internal/http-server/handlers/pr"
	"github.com/Dokhoyan/avito-pr-test/internal/logger"
	"github.com/Dokhoyan/avito-pr-test/internal/model"
	"github.com/Dokhoyan/avito-pr-test/internal/service"
	servicemocks "github.com/Dokhoyan/avito-pr-test/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(io.Discard),
		zapcore.DebugLevel,
	)
	logger.Init(core)
}

func TestReassignPR_Success(t *testing.T) {
	prID := "pr1"
	oldReviewerID := "reviewer1"
	newReviewerID := "reviewer2"

	expectedPR := &model.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   "Test PR",
		AuthorID:          "author1",
		Status:            model.PRStatusOpen,
		AssignedReviewers: []string{newReviewerID, "reviewer3"},
	}

	mockService := servicemocks.NewMockPRService(t)
	mockService.EXPECT().ReassignReviewer(mock.Anything, prID, oldReviewerID).Return(expectedPR, newReviewerID, nil)

	handler := prhandler.NewImplementation(mockService)

	reqBody := map[string]interface{}{
		"pull_request_id": prID,
		"old_reviewer_id": oldReviewerID,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Reassign(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	mockService.AssertExpectations(t)
}

func TestReassignPR_InvalidMethod(t *testing.T) {
	mockService := servicemocks.NewMockPRService(t)
	handler := prhandler.NewImplementation(mockService)

	req := httptest.NewRequest(http.MethodGet, "/pullRequest/reassign", nil)
	w := httptest.NewRecorder()

	handler.Reassign(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
}

func TestReassignPR_InvalidPayload(t *testing.T) {
	mockService := servicemocks.NewMockPRService(t)
	handler := prhandler.NewImplementation(mockService)

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	handler.Reassign(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestReassignPR_NotFound(t *testing.T) {
	prID := "non-existent-pr"
	oldReviewerID := "reviewer1"

	mockService := servicemocks.NewMockPRService(t)
	mockService.EXPECT().ReassignReviewer(mock.Anything, prID, oldReviewerID).Return(nil, "", service.ErrPRNotFound)

	handler := prhandler.NewImplementation(mockService)

	reqBody := map[string]interface{}{
		"pull_request_id": prID,
		"old_reviewer_id": oldReviewerID,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Reassign(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	mockService.AssertExpectations(t)
}

func TestReassignPR_NoReviewers(t *testing.T) {
	prID := "pr1"
	oldReviewerID := "reviewer1"

	mockService := servicemocks.NewMockPRService(t)
	mockService.EXPECT().ReassignReviewer(mock.Anything, prID, oldReviewerID).Return(nil, "", service.ErrNoReviewers)

	handler := prhandler.NewImplementation(mockService)

	reqBody := map[string]interface{}{
		"pull_request_id": prID,
		"old_reviewer_id": oldReviewerID,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Reassign(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusConflict, res.StatusCode)
	mockService.AssertExpectations(t)
}

func TestReassignPR_InternalError(t *testing.T) {
	prID := "pr1"
	oldReviewerID := "reviewer1"

	mockService := servicemocks.NewMockPRService(t)
	mockService.EXPECT().ReassignReviewer(mock.Anything, prID, oldReviewerID).Return(nil, "", assert.AnError)

	handler := prhandler.NewImplementation(mockService)

	reqBody := map[string]interface{}{
		"pull_request_id": prID,
		"old_reviewer_id": oldReviewerID,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Reassign(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	mockService.AssertExpectations(t)
}

func TestReassignPR_UserNotFound(t *testing.T) {
	prID := "pr1"
	oldReviewerID := "reviewer1"

	mockService := servicemocks.NewMockPRService(t)
	mockService.EXPECT().ReassignReviewer(mock.Anything, prID, oldReviewerID).Return(nil, "", service.ErrUserNotFound)

	handler := prhandler.NewImplementation(mockService)

	reqBody := map[string]interface{}{
		"pull_request_id": prID,
		"old_reviewer_id": oldReviewerID,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Reassign(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	mockService.AssertExpectations(t)
}

func TestReassignPR_PRMerged(t *testing.T) {
	prID := "pr1"
	oldReviewerID := "reviewer1"

	mockService := servicemocks.NewMockPRService(t)
	mockService.EXPECT().ReassignReviewer(mock.Anything, prID, oldReviewerID).Return(nil, "", service.ErrPRMerged)

	handler := prhandler.NewImplementation(mockService)

	reqBody := map[string]interface{}{
		"pull_request_id": prID,
		"old_reviewer_id": oldReviewerID,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Reassign(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusConflict, res.StatusCode)
	mockService.AssertExpectations(t)
}

func TestReassignPR_NotAssigned(t *testing.T) {
	prID := "pr1"
	oldReviewerID := "reviewer1"

	mockService := servicemocks.NewMockPRService(t)
	mockService.EXPECT().ReassignReviewer(mock.Anything, prID, oldReviewerID).Return(nil, "", service.ErrNotAssigned)

	handler := prhandler.NewImplementation(mockService)

	reqBody := map[string]interface{}{
		"pull_request_id": prID,
		"old_reviewer_id": oldReviewerID,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Reassign(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusConflict, res.StatusCode)
	mockService.AssertExpectations(t)
}

func TestReassignPR_NoCandidate(t *testing.T) {
	prID := "pr1"
	oldReviewerID := "reviewer1"

	mockService := servicemocks.NewMockPRService(t)
	mockService.EXPECT().ReassignReviewer(mock.Anything, prID, oldReviewerID).Return(nil, "", service.ErrNoCandidate)

	handler := prhandler.NewImplementation(mockService)

	reqBody := map[string]interface{}{
		"pull_request_id": prID,
		"old_reviewer_id": oldReviewerID,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Reassign(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusConflict, res.StatusCode)
	mockService.AssertExpectations(t)
}
