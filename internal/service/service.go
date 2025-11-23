package service

import (
	"context"
	"errors"

	"github.com/Dokhoyan/avito-pr-test/internal/model"
)

var (
	ErrPRExists     = errors.New("PR already exists")
	ErrNoReviewers  = errors.New("no available reviewers in team")
	ErrPRNotFound   = errors.New("PR not found")
	ErrPRMerged     = errors.New("cannot reassign reviewer for merged PR")
	ErrNotAssigned  = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate  = errors.New("no active replacement candidate in team")
	ErrUserInactive = errors.New("user is not active")
	ErrTeamExists   = errors.New("team already exists")
	ErrTeamNotFound = errors.New("team not found")
	ErrUserNotFound = errors.New("user not found")
)

type TeamService interface {
	CreateTeamWithMembers(ctx context.Context, teamName string, members []model.TeamMember) (*model.Team, error)
	GetTeam(ctx context.Context, teamName string) (*model.Team, error)
}