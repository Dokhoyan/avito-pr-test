package team

import (
	"context"
	"errors"
	"fmt"

	"github.com/Dokhoyan/avito-pr-test/internal/model"
	"github.com/Dokhoyan/avito-pr-test/internal/repository"
	"github.com/Dokhoyan/avito-pr-test/internal/service"
)

func (s *serv) GetTeam(ctx context.Context, teamName string) (*model.Team, error) {
	if teamName == "" {
		return nil, errors.New("team name cant be empty")
	}

	team, err := s.teamRepo.GetTeamWithMembers(ctx, teamName)
	if err != nil {
		if errors.Is(err, repository.ErrTeamNotFound) {
			return nil, service.ErrTeamNotFound
		}

		return nil, fmt.Errorf("failed to get team with members: %w", err)
	}

	return team, nil
}
