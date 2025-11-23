package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Dokhoyan/avito-pr-test/internal/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateOrUpdateUser(ctx context.Context, user model.TeamMember, teamID string) error {
	builder := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("user_id", "username", "team_name", "is_active", "created_at", "updated_at").
		Values(user.UserID, user.Username, teamID, user.IsActive, sq.Expr("NOW()"), sq.Expr("NOW()")).
		Suffix("ON CONFLICT (user_id) DO UPDATE SET username = EXCLUDED.username, team_name = EXCLUDED.team_name, is_active = EXCLUDED.is_active, updated_at = NOW()")

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	conn := r.getConn(ctx)

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to create or update user: %w", err)
	}

	return nil
}

func (r *Repository) UpdateIsActiveStatus(ctx context.Context, userID string, isActive bool) error {
	builder := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set("is_active", isActive).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"user_id": userID})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	conn := r.getConn(ctx)

	result, err := conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update is_active status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *Repository) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	builder := sq.Select("user_id", "username", "team_name", "is_active").
		PlaceholderFormat(sq.Dollar).
		From("users").
		Where(sq.Eq{"user_id": userID})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	conn := r.getConn(ctx)

	var user model.User

	err = conn.QueryRow(ctx, query, args...).Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

func (r *Repository) GetActiveTeamMembers(ctx context.Context, teamID, authorID string, excludeIDs ...string) ([]*model.User, error) {
	builder := sq.Select("user_id", "username", "team_name", "is_active").
		PlaceholderFormat(sq.Dollar).
		From("users").
		Where(sq.Eq{"team_name": teamID, "is_active": true})

	if authorID != "" {
		excludeIDs = append(excludeIDs, authorID)
	}

	if len(excludeIDs) > 0 {
		interfaceIDs := make([]interface{}, len(excludeIDs))
		for i, id := range excludeIDs {
			interfaceIDs[i] = id
		}
		builder = builder.Where(sq.Expr("user_id NOT IN ("+sq.Placeholders(len(excludeIDs))+")", interfaceIDs...))
	}

	builder = builder.OrderBy("user_id")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	conn := r.getConn(ctx)

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get active team members: %w", err)
	}
	defer rows.Close()

	users := make([]*model.User, 0)
	for rows.Next() {
		var user model.User

		if err := rows.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return users, nil
}

func (r *Repository) GetUserReviews(ctx context.Context, userID string) ([]*model.PullRequestShort, error) {
	builder := sq.Select(
		"p.id as pull_request_id",
		"p.name as pull_request_name",
		"p.author_id",
		"p.status",
	).
		PlaceholderFormat(sq.Dollar).
		From("pull_requests p").
		Join("pr_reviewers pr ON p.id = pr.pr_id").
		Where(sq.Eq{"pr.user_id": userID}).
		OrderBy("p.created_at DESC")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	conn := r.getConn(ctx)

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get user reviews: %w", err)
	}
	defer rows.Close()

	pullRequests := make([]*model.PullRequestShort, 0)

	for rows.Next() {
		var pr model.PullRequestShort

		err := rows.Scan(
			&pr.PullRequestID,
			&pr.PullRequestName,
			&pr.AuthorID,
			&pr.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pull request: %w", err)
		}

		pullRequests = append(pullRequests, &pr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating pull requests: %w", err)
	}

	return pullRequests, nil
}
