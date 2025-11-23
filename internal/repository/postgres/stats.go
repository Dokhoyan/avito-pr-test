package postgres

import (
	"context"
	"fmt"

	"github.com/Dokhoyan/avito-pr-test/internal/model"
	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) GetTeamsCount(ctx context.Context) (int, error) {
	builder := sq.Select("COUNT(*)").
		PlaceholderFormat(sq.Dollar).
		From("teams")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build count query: %w", err)
	}

	conn := r.getConn(ctx)

	var count int

	err = conn.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get teams count: %w", err)
	}

	return count, nil
}

func (r *Repository) GetUsersStats(ctx context.Context) (*model.UsersStats, error) {
	builder := sq.Select(
		"COUNT(*) as total_users",
		"COUNT(CASE WHEN is_active = true THEN 1 END) as active_users",
	).
		PlaceholderFormat(sq.Dollar).
		From("users")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	conn := r.getConn(ctx)

	var totalUsers, activeUsers int

	err = conn.QueryRow(ctx, query, args...).Scan(&totalUsers, &activeUsers)
	if err != nil {
		return nil, fmt.Errorf("failed to get users stats: %w", err)
	}

	return &model.UsersStats{
		TotalUsers:  totalUsers,
		ActiveUsers: activeUsers,
	}, nil
}

func (r *Repository) GetPRStats(ctx context.Context) (*model.PRStats, error) {
	builder := sq.Select(
		"COUNT(*) as total_prs",
		"COUNT(CASE WHEN status = 'OPEN' THEN 1 END) as open_prs",
		"COUNT(CASE WHEN status = 'MERGED' THEN 1 END) as merged_prs",
	).
		PlaceholderFormat(sq.Dollar).
		From("pull_requests")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	conn := r.getConn(ctx)

	var totalPRs, openPRs, mergedPRs int

	err = conn.QueryRow(ctx, query, args...).Scan(&totalPRs, &openPRs, &mergedPRs)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR stats: %w", err)
	}

	return &model.PRStats{
		TotalPRs:  totalPRs,
		OpenPRs:   openPRs,
		MergedPRs: mergedPRs,
	}, nil
}

func (r *Repository) GetAssignmentsCount(ctx context.Context) (int, error) {
	builder := sq.Select("COUNT(*)").
		PlaceholderFormat(sq.Dollar).
		From("pr_reviewers")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build count query: %w", err)
	}

	conn := r.getConn(ctx)

	var count int

	err = conn.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get assignments count: %w", err)
	}

	return count, nil
}
