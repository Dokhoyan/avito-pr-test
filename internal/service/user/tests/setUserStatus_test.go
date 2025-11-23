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

func TestUpdateUserStatus_Success(t *testing.T) {
	ctx := context.Background()
	userID := "user1"
	isActive := true

	expectedUser := &model.User{
		UserID:   userID,
		Username: "user1",
		IsActive: isActive,
	}

	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockUserRepo.EXPECT().UpdateIsActiveStatus(ctx, userID, isActive).Return(nil)
	mockUserRepo.EXPECT().GetUserByID(ctx, userID).Return(expectedUser, nil)

	svc := user.NewServiceWithTrManager(mockUserRepo, mockTrManager)
	result, err := svc.UpdateUserStatus(ctx, userID, isActive)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, isActive, result.IsActive)
	mockUserRepo.AssertExpectations(t)
}

func TestUpdateUserStatus_UserNotFound(t *testing.T) {
	ctx := context.Background()
	userID := "non-existent-user"
	isActive := true

	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockUserRepo.EXPECT().UpdateIsActiveStatus(ctx, userID, isActive).Return(repository.ErrUserNotFound)
	// GetUserByID is called after transaction, but since transaction returns error, it shouldn't be called
	// However, if the mock doesn't return error from Do, GetUserByID will be called
	// So we need to mock it to avoid unexpected call
	mockUserRepo.EXPECT().GetUserByID(ctx, userID).Return(nil, repository.ErrUserNotFound).Maybe()

	svc := user.NewServiceWithTrManager(mockUserRepo, mockTrManager)
	result, err := svc.UpdateUserStatus(ctx, userID, isActive)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, service.ErrUserNotFound))
	mockUserRepo.AssertExpectations(t)
}

func TestUpdateUserStatus_EmptyUserID(t *testing.T) {
	ctx := context.Background()
	isActive := true

	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	svc := user.NewServiceWithTrManager(mockUserRepo, mockTrManager)
	result, err := svc.UpdateUserStatus(ctx, "", isActive)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user ID cannot be empty")
}

func TestUpdateUserStatus_RepositoryError(t *testing.T) {
	ctx := context.Background()
	userID := "user1"
	isActive := true

	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockUserRepo.EXPECT().UpdateIsActiveStatus(ctx, userID, isActive).Return(errors.New("database error"))
	// GetUserByID is called after transaction, but since transaction returns error, it shouldn't be called
	// However, if the mock doesn't return error from Do, GetUserByID will be called
	// So we need to mock it to avoid unexpected call
	mockUserRepo.EXPECT().GetUserByID(ctx, userID).Return(nil, errors.New("database error")).Maybe()

	svc := user.NewServiceWithTrManager(mockUserRepo, mockTrManager)
	result, err := svc.UpdateUserStatus(ctx, userID, isActive)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to update user status")
	mockUserRepo.AssertExpectations(t)
}

func TestUpdateUserStatus_GetUserByIDError(t *testing.T) {
	ctx := context.Background()
	userID := "user1"
	isActive := true

	mockUserRepo := repomocks.NewMockUserRepository(t)
	mockTrManager := servicemocks.NewMockTrManager(t)
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	mockUserRepo.EXPECT().UpdateIsActiveStatus(ctx, userID, isActive).Return(nil)
	mockUserRepo.EXPECT().GetUserByID(ctx, userID).Return(nil, errors.New("database error"))

	svc := user.NewServiceWithTrManager(mockUserRepo, mockTrManager)
	result, err := svc.UpdateUserStatus(ctx, userID, isActive)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get user by ID")
	mockUserRepo.AssertExpectations(t)
}
