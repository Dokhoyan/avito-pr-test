package pr

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Dokhoyan/avito-pr-test/internal/model"
	"github.com/Dokhoyan/avito-pr-test/internal/repository"
	"github.com/Dokhoyan/avito-pr-test/internal/service"
)

func (s *serv) MergePR(ctx context.Context, prID string) (*model.PullRequest, error) {
	if prID == "" {
		return nil, errors.New("prID is required")
	}
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		pr, err := s.prRepo.GetPRByID(ctx, prID)
		if err != nil {
			if errors.Is(err, repository.ErrPRNotFound) {
				return service.ErrPRNotFound
			}

			return fmt.Errorf("failed to get PR: %w", err)
		}
		if pr == nil {
			return service.ErrPRNotFound
		}
		if pr.Status == model.PRStatusMerged {
			return nil
		}
		err = s.prRepo.MarkPRAsMerged(ctx, prID, model.PRStatusMerged, time.Now())
		if err != nil {
			return fmt.Errorf("failed to mark PR as merged: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	pr, err := s.prRepo.GetPRByID(ctx, prID)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR after merge: %w", err)
	}

	return pr, nil
}
