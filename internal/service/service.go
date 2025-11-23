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

// TeamService определяет методы бизнес-логики для управления командами
type TeamService interface {
	CreateTeamWithMembers(ctx context.Context, teamName string, members []model.TeamMember) (*model.Team, error)
	GetTeam(ctx context.Context, teamName string) (*model.Team, error)
}

// UserService определяет методы бизнес-логики для управления пользователями
type UserService interface {
	UpdateUserStatus(ctx context.Context, userID string, isActive bool) (*model.User, error)
	GetUserReviews(ctx context.Context, userID string) ([]*model.PullRequestShort, error)
}

// PRService определяет методы бизнес-логики для управления pull request
type PRService interface {
	CreatePR(ctx context.Context, prID string, prName string, authorID string) (*model.PullRequest, error)
	MergePR(ctx context.Context, prID string) (*model.PullRequest, error)
	ReassignReviewer(ctx context.Context, pullRequestID, oldUserID string) (*model.PullRequest, string, error)
}

// StatsService определяет методы бизнес-логики для получения статистики
type StatsService interface {
	GetStats(ctx context.Context) (*model.Stats, error)
}
