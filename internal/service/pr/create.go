package pr

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/Dokhoyan/avito-pr-test/internal/model"
	"github.com/Dokhoyan/avito-pr-test/internal/repository"
	"github.com/Dokhoyan/avito-pr-test/internal/service"
)

func (s *serv) CreatePR(ctx context.Context, prID, prName, authorID string) (*model.PullRequest, error) {
	if prID == "" || prName == "" || authorID == "" {
		return nil, errors.New("prID, prName and authorID are required")
	}
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		exists, err := s.prRepo.PRExists(ctx, prID)
		if err != nil {
			return fmt.Errorf("failed to check PR existence: %w", err)
		}
		if exists {
			return service.ErrPRExists
		}
		user, err := s.userRepo.GetUserByID(ctx, authorID)
		if err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				return service.ErrUserNotFound
			}

			return fmt.Errorf("failed to get author: %w", err)
		}
		if user == nil {
			return service.ErrUserNotFound
		}
		if !user.IsActive {
			return service.ErrUserInactive
		}
		teamExists, err := s.teamRepo.TeamExists(ctx, user.TeamName)
		if err != nil {
			return fmt.Errorf("failed to check team existence: %w", err)
		}
		if !teamExists {
			return service.ErrTeamNotFound
		}

		pr := &model.PullRequest{
			PullRequestID:   prID,
			PullRequestName: prName,
			AuthorID:        authorID,
			Status:          model.PRStatusOpen,
		}
		if errCreate := s.prRepo.CreatePR(ctx, pr); errCreate != nil {
			return fmt.Errorf("failed to create PR: %w", errCreate)
		}
		reviewers, err := s.autoAssignReviewers(ctx, user)
		if err != nil && !errors.Is(err, service.ErrNoReviewers) {
			return err
		}
		for _, reviewerID := range reviewers {
			if err := s.prRepo.AssignReviewer(ctx, prID, reviewerID); err != nil {
				return fmt.Errorf("failed to assign reviewer %s: %w", reviewerID, err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	pr, err := s.prRepo.GetPRByID(ctx, prID)
	if err != nil {
		if errors.Is(err, repository.ErrPRNotFound) {
			return nil, service.ErrPRNotFound
		}

		return nil, fmt.Errorf("failed to get PR: %w", err)
	}

	return pr, nil
}

func (s *serv) autoAssignReviewers(ctx context.Context, author *model.User) ([]string, error) {
	candidates, err := s.userRepo.GetActiveTeamMembers(ctx, author.TeamName, author.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}

	if len(candidates) == 0 {
		return nil, service.ErrNoReviewers
	}
	for i := len(candidates) - 1; i > 0; i-- {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return nil, err
		}
		j := int(n.Int64())
		candidates[i], candidates[j] = candidates[j], candidates[i]
	}

	maxReviewers := 2
	if len(candidates) < maxReviewers {
		maxReviewers = len(candidates)
	}
	reviewerIDs := make([]string, 0, maxReviewers)

	for i := 0; i < maxReviewers; i++ {
		reviewerIDs = append(reviewerIDs, candidates[i].UserID)
	}

	return reviewerIDs, nil
}
