package pr_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Dokhoyan/avito-pr-test/internal/model"
	"github.com/Dokhoyan/avito-pr-test/internal/repository"
	repomocks "github.com/Dokhoyan/avito-pr-test/internal/repository/mocks"
	"github.com/Dokhoyan/avito-pr-test/internal/service"
	servicemocks "github.com/Dokhoyan/avito-pr-test/internal/service/mocks"
	"github.com/Dokhoyan/avito-pr-test/internal/service/pr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReassignReviewer_Success(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	oldUserID := "reviewer1"
	teamName := "team1"

	existingPR := &model.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   "Test PR",
		AuthorID:          "author1",
		Status:            model.PRStatusOpen,
		AssignedReviewers: []string{oldUserID, "reviewer2"},
	}

	author := &model.User{
		UserID:   "author1",
		Username: "author1",
		TeamName: teamName,
		IsActive: true,
	}

	candidates := []*model.User{
		{UserID: "reviewer3", Username: "reviewer3", IsActive: true},
		{UserID: "reviewer4", Username: "reviewer4", IsActive: true},
	}

	updatedPR := &model.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   "Test PR",
		AuthorID:          "author1",
		Status:            model.PRStatusOpen,
		AssignedReviewers: []string{"reviewer3", "reviewer2"},
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(existingPR, nil)
	mockUserRepo.EXPECT().GetUserByID(ctx, "author1").Return(author, nil)
	mockUserRepo.EXPECT().GetActiveTeamMembers(ctx, teamName, oldUserID, "author1").Return(candidates, nil)
	mockPRRepo.EXPECT().RemoveReviewer(ctx, prID, oldUserID).Return(nil)
	mockPRRepo.EXPECT().AssignReviewer(ctx, prID, mock.AnythingOfType("string")).Return(nil)
	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(updatedPR, nil)

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, newReviewerID, err := svc.ReassignReviewer(ctx, prID, oldUserID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, newReviewerID)
	mockPRRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestReassignReviewer_PRNotFound(t *testing.T) {
	ctx := context.Background()
	prID := "non-existent-pr"
	oldUserID := "reviewer1"

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(nil, repository.ErrPRNotFound)

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, newReviewerID, err := svc.ReassignReviewer(ctx, prID, oldUserID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Empty(t, newReviewerID)
	mockPRRepo.AssertExpectations(t)
}

func TestReassignReviewer_PRMerged(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	oldUserID := "reviewer1"

	mergedPR := &model.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   "Test PR",
		AuthorID:          "author1",
		Status:            model.PRStatusMerged,
		AssignedReviewers: []string{oldUserID},
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(mergedPR, nil)

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, newReviewerID, err := svc.ReassignReviewer(ctx, prID, oldUserID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Empty(t, newReviewerID)
	assert.True(t, errors.Is(err, service.ErrPRMerged))
	mockPRRepo.AssertExpectations(t)
}

func TestReassignReviewer_NotAssigned(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	oldUserID := "not-assigned-reviewer"

	existingPR := &model.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   "Test PR",
		AuthorID:          "author1",
		Status:            model.PRStatusOpen,
		AssignedReviewers: []string{"reviewer1", "reviewer2"},
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(existingPR, nil)

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, newReviewerID, err := svc.ReassignReviewer(ctx, prID, oldUserID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Empty(t, newReviewerID)
	assert.True(t, errors.Is(err, service.ErrNotAssigned))
	mockPRRepo.AssertExpectations(t)
}

func TestReassignReviewer_NoReviewers(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	oldUserID := "reviewer1"
	teamName := "team1"

	existingPR := &model.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   "Test PR",
		AuthorID:          "author1",
		Status:            model.PRStatusOpen,
		AssignedReviewers: []string{oldUserID},
	}

	author := &model.User{
		UserID:   "author1",
		Username: "author1",
		TeamName: teamName,
		IsActive: true,
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(existingPR, nil)
	mockUserRepo.EXPECT().GetUserByID(ctx, "author1").Return(author, nil)
	mockUserRepo.EXPECT().GetActiveTeamMembers(ctx, teamName, oldUserID, "author1").Return([]*model.User{}, nil)

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, newReviewerID, err := svc.ReassignReviewer(ctx, prID, oldUserID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Empty(t, newReviewerID)
	assert.True(t, errors.Is(err, service.ErrNoReviewers))
	mockPRRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestReassignReviewer_EmptyFields(t *testing.T) {
	ctx := context.Background()

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, newReviewerID, err := svc.ReassignReviewer(ctx, "", "reviewer1")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Empty(t, newReviewerID)
	assert.Contains(t, err.Error(), "required")
}

func TestReassignReviewer_GetPRByIDError(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	oldUserID := "reviewer1"

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(nil, errors.New("database error"))

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, newReviewerID, err := svc.ReassignReviewer(ctx, prID, oldUserID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Empty(t, newReviewerID)
	assert.Contains(t, err.Error(), "failed to get PR")
	mockPRRepo.AssertExpectations(t)
}

func TestReassignReviewer_GetAuthorError(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	oldUserID := "reviewer1"

	existingPR := &model.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   "Test PR",
		AuthorID:          "author1",
		Status:            model.PRStatusOpen,
		AssignedReviewers: []string{oldUserID},
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(existingPR, nil)
	mockUserRepo.EXPECT().GetUserByID(ctx, "author1").Return(nil, errors.New("database error"))

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, newReviewerID, err := svc.ReassignReviewer(ctx, prID, oldUserID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Empty(t, newReviewerID)
	assert.Contains(t, err.Error(), "failed to get PR author")
	mockPRRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestReassignReviewer_GetActiveTeamMembersError(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	oldUserID := "reviewer1"
	teamName := "team1"

	existingPR := &model.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   "Test PR",
		AuthorID:          "author1",
		Status:            model.PRStatusOpen,
		AssignedReviewers: []string{oldUserID},
	}

	author := &model.User{
		UserID:   "author1",
		Username: "author1",
		TeamName: teamName,
		IsActive: true,
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(existingPR, nil)
	mockUserRepo.EXPECT().GetUserByID(ctx, "author1").Return(author, nil)
	mockUserRepo.EXPECT().GetActiveTeamMembers(ctx, teamName, oldUserID, "author1").Return(nil, errors.New("database error"))

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, newReviewerID, err := svc.ReassignReviewer(ctx, prID, oldUserID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Empty(t, newReviewerID)
	assert.Contains(t, err.Error(), "failed to get team members")
	mockPRRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestReassignReviewer_NoCandidate(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	oldUserID := "reviewer1"
	teamName := "team1"

	existingPR := &model.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   "Test PR",
		AuthorID:          "author1",
		Status:            model.PRStatusOpen,
		AssignedReviewers: []string{oldUserID, "reviewer2"},
	}

	author := &model.User{
		UserID:   "author1",
		Username: "author1",
		TeamName: teamName,
		IsActive: true,
	}

	// Все кандидаты уже являются ревьюерами
	candidates := []*model.User{
		{UserID: "reviewer2", Username: "reviewer2", IsActive: true},
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(existingPR, nil)
	mockUserRepo.EXPECT().GetUserByID(ctx, "author1").Return(author, nil)
	mockUserRepo.EXPECT().GetActiveTeamMembers(ctx, teamName, oldUserID, "author1").Return(candidates, nil)

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, newReviewerID, err := svc.ReassignReviewer(ctx, prID, oldUserID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Empty(t, newReviewerID)
	assert.True(t, errors.Is(err, service.ErrNoCandidate))
	mockPRRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestReassignReviewer_RemoveReviewerError(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	oldUserID := "reviewer1"
	teamName := "team1"

	existingPR := &model.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   "Test PR",
		AuthorID:          "author1",
		Status:            model.PRStatusOpen,
		AssignedReviewers: []string{oldUserID},
	}

	author := &model.User{
		UserID:   "author1",
		Username: "author1",
		TeamName: teamName,
		IsActive: true,
	}

	candidates := []*model.User{
		{UserID: "reviewer3", Username: "reviewer3", IsActive: true},
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(existingPR, nil)
	mockUserRepo.EXPECT().GetUserByID(ctx, "author1").Return(author, nil)
	mockUserRepo.EXPECT().GetActiveTeamMembers(ctx, teamName, oldUserID, "author1").Return(candidates, nil)
	mockPRRepo.EXPECT().RemoveReviewer(ctx, prID, oldUserID).Return(errors.New("database error"))

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, newReviewerID, err := svc.ReassignReviewer(ctx, prID, oldUserID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Empty(t, newReviewerID)
	assert.Contains(t, err.Error(), "failed to remove old reviewer")
	mockPRRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestReassignReviewer_AssignReviewerError(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	oldUserID := "reviewer1"
	teamName := "team1"

	existingPR := &model.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   "Test PR",
		AuthorID:          "author1",
		Status:            model.PRStatusOpen,
		AssignedReviewers: []string{oldUserID},
	}

	author := &model.User{
		UserID:   "author1",
		Username: "author1",
		TeamName: teamName,
		IsActive: true,
	}

	candidates := []*model.User{
		{UserID: "reviewer3", Username: "reviewer3", IsActive: true},
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(existingPR, nil)
	mockUserRepo.EXPECT().GetUserByID(ctx, "author1").Return(author, nil)
	mockUserRepo.EXPECT().GetActiveTeamMembers(ctx, teamName, oldUserID, "author1").Return(candidates, nil)
	mockPRRepo.EXPECT().RemoveReviewer(ctx, prID, oldUserID).Return(nil)
	mockPRRepo.EXPECT().AssignReviewer(ctx, prID, "reviewer3").Return(errors.New("database error"))

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, newReviewerID, err := svc.ReassignReviewer(ctx, prID, oldUserID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Empty(t, newReviewerID)
	assert.Contains(t, err.Error(), "failed to assign new reviewer")
	mockPRRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestReassignReviewer_GetUpdatedPRError(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	oldUserID := "reviewer1"
	teamName := "team1"

	existingPR := &model.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   "Test PR",
		AuthorID:          "author1",
		Status:            model.PRStatusOpen,
		AssignedReviewers: []string{oldUserID},
	}

	author := &model.User{
		UserID:   "author1",
		Username: "author1",
		TeamName: teamName,
		IsActive: true,
	}

	candidates := []*model.User{
		{UserID: "reviewer3", Username: "reviewer3", IsActive: true},
	}

	mockTeamRepo := repomocks.NewMockTeamRepository(t)
	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockPRRepo := repomocks.NewMockPRRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	// First GetPRByID call is inside transaction
	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(existingPR, nil).Once()
	mockUserRepo.EXPECT().GetUserByID(ctx, "author1").Return(author, nil).Once()
	mockUserRepo.EXPECT().GetActiveTeamMembers(ctx, teamName, oldUserID, "author1").Return(candidates, nil).Once()
	mockPRRepo.EXPECT().RemoveReviewer(ctx, prID, oldUserID).Return(nil).Once()
	mockPRRepo.EXPECT().AssignReviewer(ctx, prID, "reviewer3").Return(nil).Once()
	// Second GetPRByID call is inside transaction and should return error
	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(nil, errors.New("database error")).Once()

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, newReviewerID, err := svc.ReassignReviewer(ctx, prID, oldUserID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Empty(t, newReviewerID)
	assert.Contains(t, err.Error(), "failed to get updated PR")
	mockPRRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}
