package user_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	userhandler "github.com/Dokhoyan/avito-pr-test/internal/http-server/handlers/user"
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

func TestGetUserReviews_Success(t *testing.T) {
	userID := "user1"

	expectedReviews := []*model.PullRequestShort{
		{
			PullRequestID:   "pr1",
			PullRequestName: "PR 1",
			AuthorID:        "author1",
			Status:          model.PRStatusOpen,
		},
	}

	mockService := servicemocks.NewMockUserService(t)
	mockService.EXPECT().GetUserReviews(mock.Anything, userID).Return(expectedReviews, nil)

	handler := userhandler.NewImplementation(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users/getReview?user_id="+userID, nil)
	w := httptest.NewRecorder()

	handler.GetReview(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	mockService.AssertExpectations(t)
}

func TestGetUserReviews_InvalidMethod(t *testing.T) {
	mockService := servicemocks.NewMockUserService(t)
	handler := userhandler.NewImplementation(mockService)

	req := httptest.NewRequest(http.MethodPost, "/users/getReview", nil)
	w := httptest.NewRecorder()

	handler.GetReview(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
}

func TestGetUserReviews_MissingUserID(t *testing.T) {
	mockService := servicemocks.NewMockUserService(t)
	handler := userhandler.NewImplementation(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users/getReview", nil)
	w := httptest.NewRecorder()

	handler.GetReview(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestGetUserReviews_UserNotFound(t *testing.T) {
	userID := "non-existent-user"

	mockService := servicemocks.NewMockUserService(t)
	mockService.EXPECT().GetUserReviews(mock.Anything, userID).Return(nil, service.ErrUserNotFound)

	handler := userhandler.NewImplementation(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users/getReview?user_id="+userID, nil)
	w := httptest.NewRecorder()

	handler.GetReview(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	mockService.AssertExpectations(t)
}

func TestGetUserReviews_InternalError(t *testing.T) {
	userID := "user1"

	mockService := servicemocks.NewMockUserService(t)
	mockService.EXPECT().GetUserReviews(mock.Anything, userID).Return(nil, assert.AnError)

	handler := userhandler.NewImplementation(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users/getReview?user_id="+userID, nil)
	w := httptest.NewRecorder()

	handler.GetReview(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	mockService.AssertExpectations(t)
}
