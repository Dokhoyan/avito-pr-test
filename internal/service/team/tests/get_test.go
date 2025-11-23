package team_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Dokhoyan/avito-pr-test/internal/model"
	"github.com/Dokhoyan/avito-pr-test/internal/repository"
	repomocks "github.com/Dokhoyan/avito-pr-test/internal/repository/mocks"
	"github.com/Dokhoyan/avito-pr-test/internal/service"
	"github.com/Dokhoyan/avito-pr-test/internal/service/team"
	"github.com/Dokhoyan/avito-pr-test/internal/service/testutil"
	"github.com/stretchr/testify/assert"
)

func TestGetTeam_Success(t *testing.T) {
	ctx := context.Background()
	teamName := "test-team"
	expectedTeam := &model.Team{
		TeamName: teamName,
		Members: []model.TeamMember{
			{UserID: "user1", Username: "user1", IsActive: true},
		},
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	mockTeamRepo.EXPECT().GetTeamWithMembers(ctx, teamName).Return(expectedTeam, nil)

	svc := team.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockTrManager)
	result, err := svc.GetTeam(ctx, teamName)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, teamName, result.TeamName)
	assert.Len(t, result.Members, 1)
	mockTeamRepo.AssertExpectations(t)
}

func TestGetTeam_NotFound(t *testing.T) {
	ctx := context.Background()
	teamName := "non-existent-team"

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	mockTeamRepo.EXPECT().GetTeamWithMembers(ctx, teamName).Return(nil, repository.ErrTeamNotFound)

	svc := team.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockTrManager)
	result, err := svc.GetTeam(ctx, teamName)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, service.ErrTeamNotFound))
	mockTeamRepo.AssertExpectations(t)
}

func TestGetTeam_EmptyTeamName(t *testing.T) {
	ctx := context.Background()

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	svc := team.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockTrManager)
	result, err := svc.GetTeam(ctx, "")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "team name cant be empty")
}

func TestGetTeam_RepositoryError(t *testing.T) {
	ctx := context.Background()
	teamName := "test-team"

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	mockTeamRepo.EXPECT().GetTeamWithMembers(ctx, teamName).Return(nil, errors.New("database error"))

	svc := team.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockTrManager)
	result, err := svc.GetTeam(ctx, teamName)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get team with members")
	mockTeamRepo.AssertExpectations(t)
}
