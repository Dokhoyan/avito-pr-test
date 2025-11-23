package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Dokhoyan/avito-pr-test/internal/model"
	"github.com/Dokhoyan/avito-pr-test/internal/repository"
	"github.com/Dokhoyan/avito-pr-test/internal/service"
)

func (s *serv) UpdateUserStatus(ctx context.Context, userID string, isActive bool) (*model.User, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		err := s.userRepo.UpdateIsActiveStatus(ctx, userID, isActive)
		if err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				return service.ErrUserNotFound
			}

			return fmt.Errorf("failed to update user status: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}
