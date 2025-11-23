package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Dokhoyan/avito-pr-test/internal/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreatePR(ctx context.Context, pr *model.PullRequest) error {
	builder := sq.Insert("pull_requests").
		PlaceholderFormat(sq.Dollar).
		Columns("id", "name", "author_id", "status").
		Values(pr.PullRequestID, pr.PullRequestName, pr.AuthorID, pr.Status)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	conn := r.getConn(ctx)

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to create PR: %w", err)
	}

	return nil
}

func (r *Repository) PRExists(ctx context.Context, prID string) (bool, error) {
	builder := sq.Select("COUNT(*)").
		PlaceholderFormat(sq.Dollar).
		From("pull_requests").
		Where(sq.Eq{"id": prID})

	query, args, err := builder.ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to build count query: %w", err)
	}

	conn := r.getConn(ctx)

	var count int

	err = conn.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check PR existence: %w", err)
	}

	return count > 0, nil
}

func (r *Repository) GetPRByID(ctx context.Context, prID string) (*model.PullRequest, error) {
	builder := sq.Select("id", "name", "author_id", "status", "created_at", "merged_at").
		PlaceholderFormat(sq.Dollar).
		From("pull_requests").
		Where(sq.Eq{"id": prID})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	conn := r.getConn(ctx)

	var pr model.PullRequest

	var mergedAt *time.Time

	err = conn.QueryRow(ctx, query, args...).Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&mergedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPRNotFound
		}

		return nil, fmt.Errorf("failed to get PR: %w", err)
	}

	pr.MergedAt = mergedAt

	reviewers, err := r.GetPRReviewers(ctx, prID)
	if err != nil {
		return nil, err
	}

	pr.AssignedReviewers = reviewers

	return &pr, nil
}

func (r *Repository) AssignReviewer(ctx context.Context, prID, userID string) error {
	builder := sq.Insert("pr_reviewers").
		PlaceholderFormat(sq.Dollar).
		Columns("pr_id", "user_id", "assigned_at").
		Values(prID, userID, time.Now())

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	conn := r.getConn(ctx)

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to assign reviewer: %w", err)
	}

	return nil
}

func (r *Repository) GetPRReviewers(ctx context.Context, prID string) ([]string, error) {
	builder := sq.Select("user_id").
		PlaceholderFormat(sq.Dollar).
		From("pr_reviewers").
		Where(sq.Eq{"pr_id": prID}).
		OrderBy("assigned_at")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	conn := r.getConn(ctx)

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR reviewers: %w", err)
	}
	defer rows.Close()

	reviewers := make([]string, 0)
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("failed to scan reviewer: %w", err)
		}

		reviewers = append(reviewers, userID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating reviewers: %w", err)
	}

	return reviewers, nil
}

func (r *Repository) MarkPRAsMerged(ctx context.Context, prID string, status model.PRStatus, now time.Time) error {
	builder := sq.Update("pull_requests").
		PlaceholderFormat(sq.Dollar).
		Set("status", status).
		Set("merged_at", now).
		Where(sq.Eq{"id": prID})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	conn := r.getConn(ctx)

	result, err := conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to mark PR as merged: %w", err)
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return errors.New("PR already merged")
	}

	return nil
}

func (r *Repository) RemoveReviewer(ctx context.Context, prID, userID string) error {
	builder := sq.Delete("pr_reviewers").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"pr_id": prID, "user_id": userID})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	conn := r.getConn(ctx)

	result, err := conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to remove reviewer: %w", err)
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return errors.New("reviewer not found for this PR")
	}

	return nil
}
