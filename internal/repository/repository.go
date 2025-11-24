package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Dokhoyan/avito-pr-test/internal/model"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrTeamNotFound = errors.New("team not found")
	ErrPRNotFound   = errors.New("PR not found")
)

// TeamRepository определяет методы для доступа к данным команд
type TeamRepository interface {
	CreateTeam(ctx context.Context, teamName string) error
	TeamExists(ctx context.Context, name string) (bool, error)
	GetTeamWithMembers(ctx context.Context, name string) (*model.Team, error)
}

// UserRepository определяет методы для доступа к данным пользователей
type UserRepository interface {
	CreateOrUpdateUser(ctx context.Context, user model.TeamMember, teamID string) error
	UpdateIsActiveStatus(ctx context.Context, userID string, isActive bool) error
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	GetActiveTeamMembers(ctx context.Context, teamID, authorID string, excludeIDs ...string) ([]*model.User, error)
	GetUserReviews(ctx context.Context, userID string) ([]*model.PullRequestShort, error)
}

// PRRepository определяет методы для доступа к данным pull request
type PRRepository interface {
	CreatePR(ctx context.Context, pr *model.PullRequest) error
	PRExists(ctx context.Context, prID string) (bool, error)
	GetPRByID(ctx context.Context, prID string) (*model.PullRequest, error)
	AssignReviewer(ctx context.Context, prID, userID string) error
	MarkPRAsMerged(ctx context.Context, prID string, status model.PRStatus, now time.Time) error
	RemoveReviewer(ctx context.Context, prID, userID string) error
}

// StatsRepository определяет методы для доступа к статистическим данным
type StatsRepository interface {
	GetTeamsCount(ctx context.Context) (int, error)
	GetUsersStats(ctx context.Context) (*model.UsersStats, error)
	GetPRStats(ctx context.Context) (*model.PRStats, error)
	GetAssignmentsCount(ctx context.Context) (int, error)
}
