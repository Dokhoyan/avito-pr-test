package team

import (
	"context"
	"errors"
	"fmt"

	"github.com/Dokhoyan/avito-pr-test/internal/model"
	"github.com/Dokhoyan/avito-pr-test/internal/service"
)

func (s *serv) CreateTeamWithMembers(ctx context.Context, teamName string, members []model.TeamMember) (*model.Team, error) {
	if teamName == "" {
		return nil, errors.New("team name cant be empty")
	}
	if len(members) == 0 {
		return nil, errors.New("team must have at least one member")
	}

	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		exists, err := s.teamRepo.TeamExists(ctx, teamName)
		if err != nil {
			return fmt.Errorf("failed to check team existence: %w", err)
		}
		if exists {
			return service.ErrTeamExists
		}

		if err := s.teamRepo.CreateTeam(ctx, teamName); err != nil {
			return fmt.Errorf("failed to create team: %w", err)
		}

		for i, member := range members {
			if member.UserID == "" || member.Username == "" {
				return fmt.Errorf("member %d has empty user_id or username", i)
			}

			if err := s.userRepo.CreateOrUpdateUser(ctx, member, teamName); err != nil {
				return fmt.Errorf("failed to create user %s: %w", member.UserID, err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.teamRepo.GetTeamWithMembers(ctx, teamName)
}
