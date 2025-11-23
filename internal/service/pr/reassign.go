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

func (s *serv) ReassignReviewer(ctx context.Context, pullRequestID, oldUserID string) (*model.PullRequest, string, error) {
	if pullRequestID == "" || oldUserID == "" {
		return nil, "", errors.New("pullRequestID and oldUserID are required")
	}
	var newReviewerID string
	var resultPR *model.PullRequest
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		pr, err := s.prRepo.GetPRByID(ctx, pullRequestID)
		if err != nil {
			return fmt.Errorf("failed to get PR: %w", err)
		}
		if pr == nil {
			return service.ErrPRNotFound
		}
		if pr.Status == model.PRStatusMerged {
			return service.ErrPRMerged
		}

		isAssigned := false
		for _, reviewer := range pr.AssignedReviewers {
			if reviewer == oldUserID {
				isAssigned = true

				break
			}
		}
		if !isAssigned {
			return service.ErrNotAssigned
		}

		author, err := s.userRepo.GetUserByID(ctx, pr.AuthorID)
		if err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				return service.ErrUserNotFound
			}

			return fmt.Errorf("failed to get PR author: %w", err)
		}
		candidates, err := s.userRepo.GetActiveTeamMembers(ctx, author.TeamName, oldUserID, author.UserID)
		if err != nil {
			return fmt.Errorf("failed to get team members: %w", err)
		}

		if len(candidates) == 0 {
			return service.ErrNoReviewers
		}
		var availableCandidates []*model.User
		for _, candidate := range candidates {
			isAlreadyReviewer := false
			for _, reviewer := range pr.AssignedReviewers {
				if candidate.UserID == reviewer {
					isAlreadyReviewer = true

					break
				}
			}

			if !isAlreadyReviewer {
				availableCandidates = append(availableCandidates, candidate)
			}
		}
		if len(availableCandidates) == 0 {
			return service.ErrNoCandidate
		}
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(availableCandidates))))
		if err != nil {
			return fmt.Errorf("failed to generate random index: %w", err)
		}

		selectedIndex := int(n.Int64())
		newReviewerID = availableCandidates[selectedIndex].UserID

		if err = s.prRepo.RemoveReviewer(ctx, pullRequestID, oldUserID); err != nil {
			return fmt.Errorf("failed to remove old reviewer: %w", err)
		}

		if err = s.prRepo.AssignReviewer(ctx, pullRequestID, newReviewerID); err != nil {
			return fmt.Errorf("failed to assign new reviewer: %w", err)
		}

		updatedPR, err := s.prRepo.GetPRByID(ctx, pullRequestID)
		if err != nil {
			return fmt.Errorf("failed to get updated PR: %w", err)
		}
		resultPR = updatedPR

		return nil
	})
	if err != nil {
		return nil, "", err
	}

	return resultPR, newReviewerID, nil
}
