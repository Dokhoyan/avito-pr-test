package repository

import (
	"context"
	"github.com/Dokhoyan/avito-pr-test/internal/model"
)

type UserRepository interface {
	CreateOrUpdateUser(ctx context.Context, user model.TeamMember, teamID string) error
	UpdateIsActiveStatus(ctx context.Context, userID string, isActive bool) error
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	GetActiveTeamMembers(ctx context.Context, teamID, authorID string, excludeIDs ...string) ([]*model.User, error)
	GetUserReviews(ctx context.Context, userID string) ([]*model.PullRequestShort, error)
}

type TeamRepository interface {

}

type PRRepository interface {

}

type StatsRepository interface {

}