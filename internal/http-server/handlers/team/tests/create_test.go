package team_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	teamhandler "github.com/Dokhoyan/avito-pr-test/internal/http-server/handlers/team"
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

func TestCreateTeam_Success(t *testing.T) {
	teamName := "test-team"
	members := []model.TeamMember{
		{UserID: "user1", Username: "user1", IsActive: true},
	}

	expectedTeam := &model.Team{
		TeamName: teamName,
		Members:  members,
	}

	mockService := servicemocks.NewMockTeamService(t)
	mockService.EXPECT().CreateTeamWithMembers(mock.Anything, teamName, members).Return(expectedTeam, nil)

	handler := teamhandler.NewImplementation(mockService)

	reqBody := map[string]interface{}{
		"team_name": teamName,
		"members":   members,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateTeam(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	mockService.AssertExpectations(t)
}

func TestCreateTeam_InvalidMethod(t *testing.T) {
	mockService := servicemocks.NewMockTeamService(t)
	handler := teamhandler.NewImplementation(mockService)

	req := httptest.NewRequest(http.MethodGet, "/team/add", nil)
	w := httptest.NewRecorder()

	handler.CreateTeam(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
}

func TestCreateTeam_InvalidPayload(t *testing.T) {
	mockService := servicemocks.NewMockTeamService(t)
	handler := teamhandler.NewImplementation(mockService)

	req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	handler.CreateTeam(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestCreateTeam_TeamExists(t *testing.T) {
	teamName := "existing-team"
	members := []model.TeamMember{
		{UserID: "user1", Username: "user1", IsActive: true},
	}

	mockService := servicemocks.NewMockTeamService(t)
	mockService.EXPECT().CreateTeamWithMembers(mock.Anything, teamName, members).Return(nil, service.ErrTeamExists)

	handler := teamhandler.NewImplementation(mockService)

	reqBody := map[string]interface{}{
		"team_name": teamName,
		"members":   members,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateTeam(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	mockService.AssertExpectations(t)
}
