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

func TestCreatePR_Success(t *testing.T) {
	prID := "pr1"
	prName := "Test PR"
	authorID := "author1"

	expectedPR := &model.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   prName,
		AuthorID:          authorID,
		Status:            model.PRStatusOpen,
		AssignedReviewers: []string{"reviewer1", "reviewer2"},
	}

	mockService := servicemocks.NewMockPRService(t)
	mockService.EXPECT().CreatePR(mock.Anything, prID, prName, authorID).Return(expectedPR, nil)

	handler := prhandler.NewImplementation(mockService)

	reqBody := map[string]interface{}{
		"pull_request_id":   prID,
		"pull_request_name": prName,
		"author_id":         authorID,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	mockService.AssertExpectations(t)
}

func TestCreatePR_InvalidMethod(t *testing.T) {
	mockService := servicemocks.NewMockPRService(t)
	handler := prhandler.NewImplementation(mockService)

	req := httptest.NewRequest(http.MethodGet, "/pullRequest/create", nil)
	w := httptest.NewRecorder()

	handler.Create(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
}

func TestCreatePR_InvalidPayload(t *testing.T) {
	mockService := servicemocks.NewMockPRService(t)
	handler := prhandler.NewImplementation(mockService)

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestCreatePR_PRExists(t *testing.T) {
	prID := "existing-pr"
	prName := "Test PR"
	authorID := "author1"

	mockService := servicemocks.NewMockPRService(t)
	mockService.EXPECT().CreatePR(mock.Anything, prID, prName, authorID).Return(nil, service.ErrPRExists)

	handler := prhandler.NewImplementation(mockService)

	reqBody := map[string]interface{}{
		"pull_request_id":   prID,
		"pull_request_name": prName,
		"author_id":         authorID,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusConflict, res.StatusCode)
	mockService.AssertExpectations(t)
}

func TestCreatePR_UserNotFound(t *testing.T) {
	prID := "pr1"
	prName := "Test PR"
	authorID := "non-existent-author"

	mockService := servicemocks.NewMockPRService(t)
	mockService.EXPECT().CreatePR(mock.Anything, prID, prName, authorID).Return(nil, service.ErrUserNotFound)

	handler := prhandler.NewImplementation(mockService)

	reqBody := map[string]interface{}{
		"pull_request_id":   prID,
		"pull_request_name": prName,
		"author_id":         authorID,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	mockService.AssertExpectations(t)
}

func TestCreatePR_TeamNotFound(t *testing.T) {
	prID := "pr1"
	prName := "Test PR"
	authorID := "author1"

	mockService := servicemocks.NewMockPRService(t)
	mockService.EXPECT().CreatePR(mock.Anything, prID, prName, authorID).Return(nil, service.ErrTeamNotFound)

	handler := prhandler.NewImplementation(mockService)

	reqBody := map[string]interface{}{
		"pull_request_id":   prID,
		"pull_request_name": prName,
		"author_id":         authorID,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	mockService.AssertExpectations(t)
}

func TestCreatePR_UserInactive(t *testing.T) {
	prID := "pr1"
	prName := "Test PR"
	authorID := "author1"

	mockService := servicemocks.NewMockPRService(t)
	mockService.EXPECT().CreatePR(mock.Anything, prID, prName, authorID).Return(nil, service.ErrUserInactive)

	handler := prhandler.NewImplementation(mockService)

	reqBody := map[string]interface{}{
		"pull_request_id":   prID,
		"pull_request_name": prName,
		"author_id":         authorID,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusConflict, res.StatusCode)
	mockService.AssertExpectations(t)
}

func TestCreatePR_InternalError(t *testing.T) {
	prID := "pr1"
	prName := "Test PR"
	authorID := "author1"

	mockService := servicemocks.NewMockPRService(t)
	mockService.EXPECT().CreatePR(mock.Anything, prID, prName, authorID).Return(nil, assert.AnError)

	handler := prhandler.NewImplementation(mockService)

	reqBody := map[string]interface{}{
		"pull_request_id":   prID,
		"pull_request_name": prName,
		"author_id":         authorID,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	mockService.AssertExpectations(t)
}
