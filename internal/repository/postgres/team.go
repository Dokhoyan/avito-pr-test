package postgres

import (
	"context"
	"fmt"

	"github.com/Dokhoyan/avito-pr-test/internal/model"
	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) CreateTeam(ctx context.Context, teamName string) error {
	builder := sq.Insert("teams").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "created_at").
		Values(teamName, sq.Expr("NOW()"))

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	conn := r.getConn(ctx)

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	return nil
}

func (r *Repository) TeamExists(ctx context.Context, name string) (bool, error) {
	builder := sq.Select("COUNT(*)").
		PlaceholderFormat(sq.Dollar).
		From("teams").
		Where(sq.Eq{"name": name})

	query, args, err := builder.ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to build count query: %w", err)
	}

	conn := r.getConn(ctx)

	var count int

	err = conn.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check team existence: %w", err)
	}

	return count > 0, nil
}

func (r *Repository) GetTeamWithMembers(ctx context.Context, name string) (*model.Team, error) {
	conn := r.db

	exists, err := r.TeamExists(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to check team existence: %w", err)
	}

	if !exists {
		return nil, ErrTeamNotFound
	}

	builder := sq.Select("user_id", "username", "is_active").
		PlaceholderFormat(sq.Dollar).
		From("users").
		Where(sq.Eq{"team_name": name}).
		OrderBy("username")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}
	defer rows.Close()

	members := make([]model.TeamMember, 0)
	for rows.Next() {
		var member model.TeamMember
		if err := rows.Scan(&member.UserID, &member.Username, &member.IsActive); err != nil {
			return nil, fmt.Errorf("failed to scan team member: %w", err)
		}

		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating team members: %w", err)
	}

	return &model.Team{
		TeamName: name,
		Members:  members,
	}, nil
}
