package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Dokhoyan/avito-pr-test/internal/model"
	"github.com/Dokhoyan/avito-pr-test/internal/repository"
	repomocks "github.com/Dokhoyan/avito-pr-test/internal/repository/mocks"
	"github.com/Dokhoyan/avito-pr-test/internal/service"
	servicemocks "github.com/Dokhoyan/avito-pr-test/internal/service/mocks"
	"github.com/Dokhoyan/avito-pr-test/internal/service/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUserReviews_Success(t *testing.T) {
	ctx := context.Background()
	userID := "user1"

	expectedUser := &model.User{
		UserID:   userID,
		Username: "user1",
		IsActive: true,
	}

	expectedReviews := []*model.PullRequestShort{
		{
			PullRequestID:   "pr1",
			PullRequestName: "PR 1",
			AuthorID:        "author1",
			Status:          model.PRStatusOpen,
		},
		{
			PullRequestID:   "pr2",
			PullRequestName: "PR 2",
			AuthorID:        "author2",
			Status:          model.PRStatusMerged,
		},
	}

	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockUserRepo.EXPECT().GetUserByID(ctx, userID).Return(expectedUser, nil)
	mockUserRepo.EXPECT().GetUserReviews(ctx, userID).Return(expectedReviews, nil)

	svc := user.NewServiceWithTrManager(mockUserRepo, mockTrManager)
	result, err := svc.GetUserReviews(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedReviews[0].PullRequestID, result[0].PullRequestID)
	mockUserRepo.AssertExpectations(t)
}

func TestGetUserReviews_UserNotFound(t *testing.T) {
	ctx := context.Background()
	userID := "non-existent-user"

	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockUserRepo.EXPECT().GetUserByID(ctx, userID).Return(nil, repository.ErrUserNotFound)

	svc := user.NewServiceWithTrManager(mockUserRepo, mockTrManager)
	result, err := svc.GetUserReviews(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, service.ErrUserNotFound))
	mockUserRepo.AssertExpectations(t)
}

func TestGetUserReviews_EmptyUserID(t *testing.T) {
	ctx := context.Background()

	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	svc := user.NewServiceWithTrManager(mockUserRepo, mockTrManager)
	result, err := svc.GetUserReviews(ctx, "")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user ID cannot be empty")
}

func TestGetUserReviews_NoReviews(t *testing.T) {
	ctx := context.Background()
	userID := "user1"

	expectedUser := &model.User{
		UserID:   userID,
		Username: "user1",
		IsActive: true,
	}

	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockUserRepo.EXPECT().GetUserByID(ctx, userID).Return(expectedUser, nil)
	mockUserRepo.EXPECT().GetUserReviews(ctx, userID).Return(nil, nil)

	svc := user.NewServiceWithTrManager(mockUserRepo, mockTrManager)
	result, err := svc.GetUserReviews(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
	mockUserRepo.AssertExpectations(t)
}

func TestGetUserReviews_UserIsNil(t *testing.T) {
	ctx := context.Background()
	userID := "user1"

	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockUserRepo.EXPECT().GetUserByID(ctx, userID).Return(nil, nil)

	svc := user.NewServiceWithTrManager(mockUserRepo, mockTrManager)
	result, err := svc.GetUserReviews(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, service.ErrUserNotFound))
	mockUserRepo.AssertExpectations(t)
}

func TestGetUserReviews_GetUserError(t *testing.T) {
	ctx := context.Background()
	userID := "user1"

	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockUserRepo.EXPECT().GetUserByID(ctx, userID).Return(nil, errors.New("database error"))

	svc := user.NewServiceWithTrManager(mockUserRepo, mockTrManager)
	result, err := svc.GetUserReviews(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get user")
	mockUserRepo.AssertExpectations(t)
}

func TestGetUserReviews_GetReviewsError(t *testing.T) {
	ctx := context.Background()
	userID := "user1"

	expectedUser := &model.User{
		UserID:   userID,
		Username: "user1",
		IsActive: true,
	}

	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockUserRepo.EXPECT().GetUserByID(ctx, userID).Return(expectedUser, nil)
	mockUserRepo.EXPECT().GetUserReviews(ctx, userID).Return(nil, errors.New("database error"))

	svc := user.NewServiceWithTrManager(mockUserRepo, mockTrManager)
	result, err := svc.GetUserReviews(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get user reviews")
	mockUserRepo.AssertExpectations(t)
}
