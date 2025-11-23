package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dokhoyan/avito-pr-test/internal/http-server/handlers/testutil"
	userhandler "github.com/Dokhoyan/avito-pr-test/internal/http-server/handlers/user"
	"github.com/Dokhoyan/avito-pr-test/internal/model"
	"github.com/Dokhoyan/avito-pr-test/internal/service"
	servicemocks "github.com/Dokhoyan/avito-pr-test/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	testutil.InitTestLogger()
}

func TestSetUserStatus_Success(t *testing.T) {
	userID := "user1"
	isActive := true

	expectedUser := &model.User{
		UserID:   userID,
		Username: "user1",
		IsActive: isActive,
	}

	mockService := servicemocks.NewMockUserService(t)
	mockService.EXPECT().UpdateUserStatus(mock.Anything, userID, isActive).Return(expectedUser, nil)

	handler := userhandler.NewImplementation(mockService)

	reqBody := map[string]interface{}{
		"user_id":   userID,
		"is_active": isActive,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.SetUserStatus(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	mockService.AssertExpectations(t)
}

func TestSetUserStatus_InvalidMethod(t *testing.T) {
	mockService := servicemocks.NewMockUserService(t)
	handler := userhandler.NewImplementation(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users/setIsActive", nil)
	w := httptest.NewRecorder()

	handler.SetUserStatus(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
}

func TestSetUserStatus_InvalidPayload(t *testing.T) {
	mockService := servicemocks.NewMockUserService(t)
	handler := userhandler.NewImplementation(mockService)

	req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	handler.SetUserStatus(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestSetUserStatus_UserNotFound(t *testing.T) {
	userID := "non-existent-user"
	isActive := true

	mockService := servicemocks.NewMockUserService(t)
	mockService.EXPECT().UpdateUserStatus(mock.Anything, userID, isActive).Return(nil, service.ErrUserNotFound)

	handler := userhandler.NewImplementation(mockService)

	reqBody := map[string]interface{}{
		"user_id":   userID,
		"is_active": isActive,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.SetUserStatus(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	mockService.AssertExpectations(t)
}
