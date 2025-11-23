package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Dokhoyan/avito-pr-test/internal/model"
	"github.com/Dokhoyan/avito-pr-test/internal/repository"
	"github.com/Dokhoyan/avito-pr-test/internal/service"
)

func (s *serv) GetUserReviews(ctx context.Context, userID string) ([]*model.PullRequestShort, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, service.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, service.ErrUserNotFound
	}

	pullRequests, err := s.userRepo.GetUserReviews(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user reviews: %w", err)
	}

	if pullRequests == nil {
		pullRequests = []*model.PullRequestShort{}
	}

	return pullRequests, nil
}
