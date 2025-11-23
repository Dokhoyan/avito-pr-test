package pr_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Dokhoyan/avito-pr-test/internal/model"
	"github.com/Dokhoyan/avito-pr-test/internal/repository"
	repomocks "github.com/Dokhoyan/avito-pr-test/internal/repository/mocks"
	"github.com/Dokhoyan/avito-pr-test/internal/service"
	"github.com/Dokhoyan/avito-pr-test/internal/service/pr"
	"github.com/Dokhoyan/avito-pr-test/internal/service/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreatePR_Success(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	prName := "Test PR"
	authorID := "author1"
	teamName := "team1"

	author := &model.User{
		UserID:   authorID,
		Username: "author1",
		TeamName: teamName,
		IsActive: true,
	}

	candidates := []*model.User{
		{UserID: "reviewer1", Username: "reviewer1", IsActive: true},
		{UserID: "reviewer2", Username: "reviewer2", IsActive: true},
	}

	expectedPR := &model.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   prName,
		AuthorID:          authorID,
		Status:            model.PRStatusOpen,
		AssignedReviewers: []string{"reviewer1", "reviewer2"},
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	mockPRRepo.EXPECT().PRExists(ctx, prID).Return(false, nil)
	mockUserRepo.EXPECT().GetUserByID(ctx, authorID).Return(author, nil)
	mockTeamRepo.EXPECT().TeamExists(ctx, teamName).Return(true, nil)
	mockPRRepo.EXPECT().CreatePR(ctx, mock.AnythingOfType("*model.PullRequest")).Return(nil)
	mockUserRepo.EXPECT().GetActiveTeamMembers(ctx, teamName, authorID).Return(candidates, nil)
	mockPRRepo.EXPECT().AssignReviewer(ctx, prID, "reviewer1").Return(nil)
	mockPRRepo.EXPECT().AssignReviewer(ctx, prID, "reviewer2").Return(nil)
	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(expectedPR, nil)

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, err := svc.CreatePR(ctx, prID, prName, authorID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, prID, result.PullRequestID)
	mockPRRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockTeamRepo.AssertExpectations(t)
}

func TestCreatePR_PRExists(t *testing.T) {
	ctx := context.Background()
	prID := "existing-pr"
	prName := "Test PR"
	authorID := "author1"

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	mockPRRepo.EXPECT().PRExists(ctx, prID).Return(true, nil)

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, err := svc.CreatePR(ctx, prID, prName, authorID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, service.ErrPRExists))
	mockPRRepo.AssertExpectations(t)
}

func TestCreatePR_EmptyFields(t *testing.T) {
	ctx := context.Background()

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, err := svc.CreatePR(ctx, "", "name", "author")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "required")
}

func TestCreatePR_UserNotFound(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	prName := "Test PR"
	authorID := "non-existent-author"

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	mockPRRepo.EXPECT().PRExists(ctx, prID).Return(false, nil)
	mockUserRepo.EXPECT().GetUserByID(ctx, authorID).Return(nil, repository.ErrUserNotFound)

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, err := svc.CreatePR(ctx, prID, prName, authorID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, service.ErrUserNotFound))
	mockPRRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestCreatePR_UserInactive(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	prName := "Test PR"
	authorID := "inactive-author"
	teamName := "team1"

	author := &model.User{
		UserID:   authorID,
		Username: "inactive-author",
		TeamName: teamName,
		IsActive: false,
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	mockPRRepo.EXPECT().PRExists(ctx, prID).Return(false, nil)
	mockUserRepo.EXPECT().GetUserByID(ctx, authorID).Return(author, nil)

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, err := svc.CreatePR(ctx, prID, prName, authorID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, service.ErrUserInactive))
	mockPRRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestCreatePR_NoReviewers(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	prName := "Test PR"
	authorID := "author1"
	teamName := "team1"

	author := &model.User{
		UserID:   authorID,
		Username: "author1",
		TeamName: teamName,
		IsActive: true,
	}

	expectedPR := &model.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   prName,
		AuthorID:          authorID,
		Status:            model.PRStatusOpen,
		AssignedReviewers: []string{},
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	mockPRRepo.EXPECT().PRExists(ctx, prID).Return(false, nil)
	mockUserRepo.EXPECT().GetUserByID(ctx, authorID).Return(author, nil)
	mockTeamRepo.EXPECT().TeamExists(ctx, teamName).Return(true, nil)
	mockPRRepo.EXPECT().CreatePR(ctx, mock.AnythingOfType("*model.PullRequest")).Return(nil)
	mockUserRepo.EXPECT().GetActiveTeamMembers(ctx, teamName, authorID).Return([]*model.User{}, nil)
	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(expectedPR, nil)

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, err := svc.CreatePR(ctx, prID, prName, authorID)

	// PR should be created even without reviewers
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockPRRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockTeamRepo.AssertExpectations(t)
}

func TestCreatePR_PRExistsError(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	prName := "Test PR"
	authorID := "author1"

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	mockPRRepo.EXPECT().PRExists(ctx, prID).Return(false, errors.New("database error"))
	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(nil, errors.New("pr not found")).Maybe()

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, err := svc.CreatePR(ctx, prID, prName, authorID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to check PR existence")
	mockPRRepo.AssertExpectations(t)
}

func TestCreatePR_GetUserError(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	prName := "Test PR"
	authorID := "author1"

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	mockPRRepo.EXPECT().PRExists(ctx, prID).Return(false, nil)
	mockUserRepo.EXPECT().GetUserByID(ctx, authorID).Return(nil, errors.New("database error"))
	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(nil, errors.New("pr not found")).Maybe()

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, err := svc.CreatePR(ctx, prID, prName, authorID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get author")
	mockPRRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestCreatePR_UserIsNil(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	prName := "Test PR"
	authorID := "author1"

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	mockPRRepo.EXPECT().PRExists(ctx, prID).Return(false, nil)
	mockUserRepo.EXPECT().GetUserByID(ctx, authorID).Return(nil, nil)
	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(nil, errors.New("pr not found")).Maybe()

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, err := svc.CreatePR(ctx, prID, prName, authorID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, service.ErrUserNotFound))
	mockPRRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestCreatePR_TeamExistsError(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	prName := "Test PR"
	authorID := "author1"
	teamName := "team1"

	author := &model.User{
		UserID:   authorID,
		Username: "author1",
		TeamName: teamName,
		IsActive: true,
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	mockPRRepo.EXPECT().PRExists(ctx, prID).Return(false, nil)
	mockUserRepo.EXPECT().GetUserByID(ctx, authorID).Return(author, nil)
	mockTeamRepo.EXPECT().TeamExists(ctx, teamName).Return(false, errors.New("database error"))
	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(nil, errors.New("pr not found")).Maybe()

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, err := svc.CreatePR(ctx, prID, prName, authorID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to check team existence")
	mockPRRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockTeamRepo.AssertExpectations(t)
}

func TestCreatePR_TeamNotFound(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	prName := "Test PR"
	authorID := "author1"
	teamName := "team1"

	author := &model.User{
		UserID:   authorID,
		Username: "author1",
		TeamName: teamName,
		IsActive: true,
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	mockPRRepo.EXPECT().PRExists(ctx, prID).Return(false, nil)
	mockUserRepo.EXPECT().GetUserByID(ctx, authorID).Return(author, nil)
	mockTeamRepo.EXPECT().TeamExists(ctx, teamName).Return(false, nil)
	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(nil, errors.New("pr not found")).Maybe()

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, err := svc.CreatePR(ctx, prID, prName, authorID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, service.ErrTeamNotFound))
	mockPRRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockTeamRepo.AssertExpectations(t)
}

func TestCreatePR_CreatePRError(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	prName := "Test PR"
	authorID := "author1"
	teamName := "team1"

	author := &model.User{
		UserID:   authorID,
		Username: "author1",
		TeamName: teamName,
		IsActive: true,
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	mockPRRepo.EXPECT().PRExists(ctx, prID).Return(false, nil)
	mockUserRepo.EXPECT().GetUserByID(ctx, authorID).Return(author, nil)
	mockTeamRepo.EXPECT().TeamExists(ctx, teamName).Return(true, nil)
	mockPRRepo.EXPECT().CreatePR(ctx, mock.AnythingOfType("*model.PullRequest")).Return(errors.New("database error"))
	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(nil, errors.New("pr not found")).Maybe()

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, err := svc.CreatePR(ctx, prID, prName, authorID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create PR")
	mockPRRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockTeamRepo.AssertExpectations(t)
}

func TestCreatePR_AssignReviewerError(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	prName := "Test PR"
	authorID := "author1"
	teamName := "team1"

	author := &model.User{
		UserID:   authorID,
		Username: "author1",
		TeamName: teamName,
		IsActive: true,
	}

	candidates := []*model.User{
		{UserID: "reviewer1", Username: "reviewer1", IsActive: true},
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	mockPRRepo.EXPECT().PRExists(ctx, prID).Return(false, nil)
	mockUserRepo.EXPECT().GetUserByID(ctx, authorID).Return(author, nil)
	mockTeamRepo.EXPECT().TeamExists(ctx, teamName).Return(true, nil)
	mockPRRepo.EXPECT().CreatePR(ctx, mock.AnythingOfType("*model.PullRequest")).Return(nil)
	mockUserRepo.EXPECT().GetActiveTeamMembers(ctx, teamName, authorID).Return(candidates, nil)
	mockPRRepo.EXPECT().AssignReviewer(ctx, prID, "reviewer1").Return(errors.New("database error"))
	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(nil, errors.New("pr not found")).Maybe()

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, err := svc.CreatePR(ctx, prID, prName, authorID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to assign reviewer")
	mockPRRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockTeamRepo.AssertExpectations(t)
}

func TestCreatePR_GetActiveTeamMembersError(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	prName := "Test PR"
	authorID := "author1"
	teamName := "team1"

	author := &model.User{
		UserID:   authorID,
		Username: "author1",
		TeamName: teamName,
		IsActive: true,
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	mockPRRepo.EXPECT().PRExists(ctx, prID).Return(false, nil)
	mockUserRepo.EXPECT().GetUserByID(ctx, authorID).Return(author, nil)
	mockTeamRepo.EXPECT().TeamExists(ctx, teamName).Return(true, nil)
	mockPRRepo.EXPECT().CreatePR(ctx, mock.AnythingOfType("*model.PullRequest")).Return(nil)
	mockUserRepo.EXPECT().GetActiveTeamMembers(ctx, teamName, authorID).Return(nil, errors.New("database error"))
	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(nil, errors.New("pr not found")).Maybe()

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, err := svc.CreatePR(ctx, prID, prName, authorID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get team members")
	mockPRRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockTeamRepo.AssertExpectations(t)
}

func TestCreatePR_GetPRByIDError(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	prName := "Test PR"
	authorID := "author1"
	teamName := "team1"

	author := &model.User{
		UserID:   authorID,
		Username: "author1",
		TeamName: teamName,
		IsActive: true,
	}

	candidates := []*model.User{
		{UserID: "reviewer1", Username: "reviewer1", IsActive: true},
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	mockPRRepo.EXPECT().PRExists(ctx, prID).Return(false, nil)
	mockUserRepo.EXPECT().GetUserByID(ctx, authorID).Return(author, nil)
	mockTeamRepo.EXPECT().TeamExists(ctx, teamName).Return(true, nil)
	mockPRRepo.EXPECT().CreatePR(ctx, mock.AnythingOfType("*model.PullRequest")).Return(nil)
	mockUserRepo.EXPECT().GetActiveTeamMembers(ctx, teamName, authorID).Return(candidates, nil)
	mockPRRepo.EXPECT().AssignReviewer(ctx, prID, "reviewer1").Return(nil)
	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(nil, repository.ErrPRNotFound)

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, err := svc.CreatePR(ctx, prID, prName, authorID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, service.ErrPRNotFound))
	mockPRRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockTeamRepo.AssertExpectations(t)
}

func TestCreatePR_GetPRByIDGenericError(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	prName := "Test PR"
	authorID := "author1"
	teamName := "team1"

	author := &model.User{
		UserID:   authorID,
		Username: "author1",
		TeamName: teamName,
		IsActive: true,
	}

	candidates := []*model.User{
		{UserID: "reviewer1", Username: "reviewer1", IsActive: true},
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	mockPRRepo.EXPECT().PRExists(ctx, prID).Return(false, nil)
	mockUserRepo.EXPECT().GetUserByID(ctx, authorID).Return(author, nil)
	mockTeamRepo.EXPECT().TeamExists(ctx, teamName).Return(true, nil)
	mockPRRepo.EXPECT().CreatePR(ctx, mock.AnythingOfType("*model.PullRequest")).Return(nil)
	mockUserRepo.EXPECT().GetActiveTeamMembers(ctx, teamName, authorID).Return(candidates, nil)
	mockPRRepo.EXPECT().AssignReviewer(ctx, prID, "reviewer1").Return(nil)
	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(nil, errors.New("database error"))

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, err := svc.CreatePR(ctx, prID, prName, authorID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get PR")
	mockPRRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockTeamRepo.AssertExpectations(t)
}
