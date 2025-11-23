package stats

import (
	"context"

	"github.com/Dokhoyan/avito-pr-test/internal/model"
)

func (s *serv) GetStats(ctx context.Context) (*model.Stats, error) {
	var stats *model.Stats

	err := s.tr.Do(ctx, func(ctx context.Context) error {
		teamsCount, err := s.repo.GetTeamsCount(ctx)
		if err != nil {
			return err
		}

		usersStats, err := s.repo.GetUsersStats(ctx)
		if err != nil {
			return err
		}

		prStats, err := s.repo.GetPRStats(ctx)
		if err != nil {
			return err
		}

		assignmentsCount, err := s.repo.GetAssignmentsCount(ctx)
		if err != nil {
			return err
		}

		stats = &model.Stats{
			TotalTeams:       teamsCount,
			TotalUsers:       usersStats.TotalUsers,
			ActiveUsers:      usersStats.ActiveUsers,
			TotalPRs:         prStats.TotalPRs,
			OpenPRs:          prStats.OpenPRs,
			MergedPRs:        prStats.MergedPRs,
			TotalAssignments: assignmentsCount,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return stats, nil
}
