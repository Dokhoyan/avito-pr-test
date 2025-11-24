package pr_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Dokhoyan/avito-pr-test/internal/model"
	"github.com/Dokhoyan/avito-pr-test/internal/repository"
	repomocks "github.com/Dokhoyan/avito-pr-test/internal/repository/mocks"
	"github.com/Dokhoyan/avito-pr-test/internal/service"
	servicemocks "github.com/Dokhoyan/avito-pr-test/internal/service/mocks"
	"github.com/Dokhoyan/avito-pr-test/internal/service/pr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMergePR_Success(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"

	existingPR := &model.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   "Test PR",
		AuthorID:          "author1",
		Status:            model.PRStatusOpen,
		AssignedReviewers: []string{"reviewer1"},
	}

	mergedPR := &model.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   "Test PR",
		AuthorID:          "author1",
		Status:            model.PRStatusMerged,
		AssignedReviewers: []string{"reviewer1"},
		MergedAt:          func() *time.Time { t := time.Now(); return &t }(),
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
	// MarkPRAsMerged is inside transaction
	mockPRRepo.EXPECT().MarkPRAsMerged(ctx, prID, model.PRStatusMerged, mock.AnythingOfType("time.Time")).Return(nil).Once()
	// Second GetPRByID call is after transaction
	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(mergedPR, nil).Once()

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, err := svc.MergePR(ctx, prID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, model.PRStatusMerged, result.Status)
	mockPRRepo.AssertExpectations(t)
}

func TestMergePR_AlreadyMerged(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"

	mergedPR := &model.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   "Test PR",
		AuthorID:          "author1",
		Status:            model.PRStatusMerged,
		AssignedReviewers: []string{"reviewer1"},
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
	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(mergedPR, nil)

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, err := svc.MergePR(ctx, prID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, model.PRStatusMerged, result.Status)
	mockPRRepo.AssertExpectations(t)
}

func TestMergePR_NotFound(t *testing.T) {
	ctx := context.Background()
	prID := "non-existent-pr"

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
	result, err := svc.MergePR(ctx, prID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, service.ErrPRNotFound))
	mockPRRepo.AssertExpectations(t)
}

func TestMergePR_EmptyPRID(t *testing.T) {
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
	result, err := svc.MergePR(ctx, "")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "required")
}

func TestMergePR_GetPRByIDError(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"

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
	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(nil, errors.New("pr not found")).Maybe()

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, err := svc.MergePR(ctx, prID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get PR")
	mockPRRepo.AssertExpectations(t)
}

func TestMergePR_MarkPRAsMergedError(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"

	existingPR := &model.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   "Test PR",
		AuthorID:          "author1",
		Status:            model.PRStatusOpen,
		AssignedReviewers: []string{"reviewer1"},
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

	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(existingPR, nil).Once()
	mockPRRepo.EXPECT().MarkPRAsMerged(ctx, prID, model.PRStatusMerged, mock.AnythingOfType("time.Time")).Return(errors.New("database error")).Once()
	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(nil, errors.New("pr not found")).Maybe()

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, err := svc.MergePR(ctx, prID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to mark PR as merged")
	mockPRRepo.AssertExpectations(t)
}

func TestMergePR_GetPRAfterMergeError(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"

	existingPR := &model.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   "Test PR",
		AuthorID:          "author1",
		Status:            model.PRStatusOpen,
		AssignedReviewers: []string{"reviewer1"},
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

	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(existingPR, nil).Once()
	mockPRRepo.EXPECT().MarkPRAsMerged(ctx, prID, model.PRStatusMerged, mock.AnythingOfType("time.Time")).Return(nil).Once()
	mockPRRepo.EXPECT().GetPRByID(ctx, prID).Return(nil, errors.New("database error")).Once()

	svc := pr.NewServiceWithTrManager(mockTeamRepo, mockUserRepo, mockPRRepo, mockTrManager)
	result, err := svc.MergePR(ctx, prID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get PR after merge")
	mockPRRepo.AssertExpectations(t)
}
