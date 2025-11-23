package team_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	teamhandler "github.com/Dokhoyan/avito-pr-test/internal/http-server/handlers/team"
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

func TestGetTeam_Success(t *testing.T) {
	teamName := "test-team"
	expectedTeam := &model.Team{
		TeamName: teamName,
		Members: []model.TeamMember{
			{UserID: "user1", Username: "user1", IsActive: true},
		},
	}

	mockService := servicemocks.NewMockTeamService(t)
	mockService.EXPECT().GetTeam(mock.Anything, teamName).Return(expectedTeam, nil)

	handler := teamhandler.NewImplementation(mockService)

	req := httptest.NewRequest(http.MethodGet, "/team/get?team_name="+teamName, nil)
	w := httptest.NewRecorder()

	handler.GetTeam(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	mockService.AssertExpectations(t)
}

func TestGetTeam_InvalidMethod(t *testing.T) {
	mockService := servicemocks.NewMockTeamService(t)
	handler := teamhandler.NewImplementation(mockService)

	req := httptest.NewRequest(http.MethodPost, "/team/get", nil)
	w := httptest.NewRecorder()

	handler.GetTeam(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
}

func TestGetTeam_MissingTeamName(t *testing.T) {
	mockService := servicemocks.NewMockTeamService(t)
	handler := teamhandler.NewImplementation(mockService)

	req := httptest.NewRequest(http.MethodGet, "/team/get", nil)
	w := httptest.NewRecorder()

	handler.GetTeam(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestGetTeam_NotFound(t *testing.T) {
	teamName := "non-existent-team"

	mockService := servicemocks.NewMockTeamService(t)
	mockService.EXPECT().GetTeam(mock.Anything, teamName).Return(nil, service.ErrTeamNotFound)

	handler := teamhandler.NewImplementation(mockService)

	req := httptest.NewRequest(http.MethodGet, "/team/get?team_name="+teamName, nil)
	w := httptest.NewRecorder()

	handler.GetTeam(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	mockService.AssertExpectations(t)
}

func TestGetTeam_InternalError(t *testing.T) {
	teamName := "test-team"

	mockService := servicemocks.NewMockTeamService(t)
	mockService.EXPECT().GetTeam(mock.Anything, teamName).Return(nil, assert.AnError)

	handler := teamhandler.NewImplementation(mockService)

	req := httptest.NewRequest(http.MethodGet, "/team/get?team_name="+teamName, nil)
	w := httptest.NewRecorder()

	handler.GetTeam(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	mockService.AssertExpectations(t)
}
