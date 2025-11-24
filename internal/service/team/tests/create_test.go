package team_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Dokhoyan/avito-pr-test/internal/model"
	repomocks "github.com/Dokhoyan/avito-pr-test/internal/repository/mocks"
	"github.com/Dokhoyan/avito-pr-test/internal/service"
	servicemocks "github.com/Dokhoyan/avito-pr-test/internal/service/mocks"
	"github.com/Dokhoyan/avito-pr-test/internal/service/team"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTeamWithMembers_Success(t *testing.T) {
	ctx := context.Background()
	teamName := "test-team"
	members := []model.TeamMember{
		{UserID: "user1", Username: "user1", IsActive: true},
		{UserID: "user2", Username: "user2", IsActive: true},
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockTeamRepo.EXPECT().TeamExists(ctx, teamName).Return(false, nil)
	mockTeamRepo.EXPECT().CreateTeam(ctx, teamName).Return(nil)
	mockUserRepo.EXPECT().CreateOrUpdateUser(ctx, members[0], teamName).Return(nil)
	mockUserRepo.EXPECT().CreateOrUpdateUser(ctx, members[1], teamName).Return(nil)

	expectedTeam := &model.Team{
		TeamName: teamName,
		Members:  members,
	}
	mockTeamRepo.EXPECT().GetTeamWithMembers(ctx, teamName).Return(expectedTeam, nil)

	svc := team.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockTrManager)
	result, err := svc.CreateTeamWithMembers(ctx, teamName, members)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, teamName, result.TeamName)
	assert.Len(t, result.Members, 2)
	mockTeamRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestCreateTeamWithMembers_TeamExists(t *testing.T) {
	ctx := context.Background()
	teamName := "existing-team"
	members := []model.TeamMember{
		{UserID: "user1", Username: "user1", IsActive: true},
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockTeamRepo.EXPECT().TeamExists(ctx, teamName).Return(true, nil)
	// GetTeamWithMembers is called after transaction, but since transaction returns error, it shouldn't be called
	// However, if the mock doesn't return error from Do, GetTeamWithMembers will be called
	// So we need to mock it to avoid unexpected call
	mockTeamRepo.EXPECT().GetTeamWithMembers(ctx, teamName).Return(nil, errors.New("team not found")).Maybe()

	svc := team.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockTrManager)
	result, err := svc.CreateTeamWithMembers(ctx, teamName, members)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, service.ErrTeamExists))
	mockTeamRepo.AssertExpectations(t)
}

func TestCreateTeamWithMembers_EmptyTeamName(t *testing.T) {
	ctx := context.Background()
	members := []model.TeamMember{
		{UserID: "user1", Username: "user1", IsActive: true},
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	svc := team.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockTrManager)
	result, err := svc.CreateTeamWithMembers(ctx, "", members)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "team name cant be empty")
}

func TestCreateTeamWithMembers_EmptyMembers(t *testing.T) {
	ctx := context.Background()
	teamName := "test-team"

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	svc := team.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockTrManager)
	result, err := svc.CreateTeamWithMembers(ctx, teamName, []model.TeamMember{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "team must have at least one member")
}

func TestCreateTeamWithMembers_InvalidMember(t *testing.T) {
	ctx := context.Background()
	teamName := "test-team"
	members := []model.TeamMember{
		{UserID: "", Username: "user1", IsActive: true},
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockTeamRepo.EXPECT().TeamExists(ctx, teamName).Return(false, nil)
	mockTeamRepo.EXPECT().CreateTeam(ctx, teamName).Return(nil)
	// GetTeamWithMembers is called after transaction, but since transaction returns error, it shouldn't be called
	// However, if the mock doesn't return error from Do, GetTeamWithMembers will be called
	// So we need to mock it to avoid unexpected call
	mockTeamRepo.EXPECT().GetTeamWithMembers(ctx, teamName).Return(nil, errors.New("team not found")).Maybe()

	svc := team.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockTrManager)
	result, err := svc.CreateTeamWithMembers(ctx, teamName, members)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "empty user_id or username")
	mockTeamRepo.AssertExpectations(t)
}

func TestCreateTeamWithMembers_TeamExistsError(t *testing.T) {
	ctx := context.Background()
	teamName := "test-team"
	members := []model.TeamMember{
		{UserID: "user1", Username: "user1", IsActive: true},
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockTeamRepo.EXPECT().TeamExists(ctx, teamName).Return(false, errors.New("database error"))
	// GetTeamWithMembers is called after transaction, but since transaction returns error, it shouldn't be called
	// However, if the mock doesn't return error from Do, GetTeamWithMembers will be called
	// So we need to mock it to avoid unexpected call
	mockTeamRepo.EXPECT().GetTeamWithMembers(ctx, teamName).Return(nil, errors.New("team not found")).Maybe()

	svc := team.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockTrManager)
	result, err := svc.CreateTeamWithMembers(ctx, teamName, members)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockTeamRepo.AssertExpectations(t)
}

func TestCreateTeamWithMembers_CreateTeamError(t *testing.T) {
	ctx := context.Background()
	teamName := "test-team"
	members := []model.TeamMember{
		{UserID: "user1", Username: "user1", IsActive: true},
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockTeamRepo.EXPECT().TeamExists(ctx, teamName).Return(false, nil)
	mockTeamRepo.EXPECT().CreateTeam(ctx, teamName).Return(errors.New("database error"))
	mockTeamRepo.EXPECT().GetTeamWithMembers(ctx, teamName).Return(nil, errors.New("team not found")).Maybe()

	svc := team.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockTrManager)
	result, err := svc.CreateTeamWithMembers(ctx, teamName, members)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create team")
	mockTeamRepo.AssertExpectations(t)
}

func TestCreateTeamWithMembers_CreateUserError(t *testing.T) {
	ctx := context.Background()
	teamName := "test-team"
	members := []model.TeamMember{
		{UserID: "user1", Username: "user1", IsActive: true},
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockTeamRepo.EXPECT().TeamExists(ctx, teamName).Return(false, nil)
	mockTeamRepo.EXPECT().CreateTeam(ctx, teamName).Return(nil)
	mockUserRepo.EXPECT().CreateOrUpdateUser(ctx, members[0], teamName).Return(errors.New("database error"))
	mockTeamRepo.EXPECT().GetTeamWithMembers(ctx, teamName).Return(nil, errors.New("team not found")).Maybe()

	svc := team.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockTrManager)
	result, err := svc.CreateTeamWithMembers(ctx, teamName, members)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create user")
	mockTeamRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestCreateTeamWithMembers_GetTeamError(t *testing.T) {
	ctx := context.Background()
	teamName := "test-team"
	members := []model.TeamMember{
		{UserID: "user1", Username: "user1", IsActive: true},
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockTeamRepo.EXPECT().TeamExists(ctx, teamName).Return(false, nil)
	mockTeamRepo.EXPECT().CreateTeam(ctx, teamName).Return(nil)
	mockUserRepo.EXPECT().CreateOrUpdateUser(ctx, members[0], teamName).Return(nil)
	mockTeamRepo.EXPECT().GetTeamWithMembers(ctx, teamName).Return(nil, errors.New("database error"))

	svc := team.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockTrManager)
	result, err := svc.CreateTeamWithMembers(ctx, teamName, members)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockTeamRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}
