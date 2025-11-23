package postgres

import (
	"context"
	"errors"

	txmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db     *pgxpool.Pool
	getter *txmpgx.CtxGetter
}

var (
	ErrUserNotFound = errors.New("user not found")
	ErrTeamNotFound = errors.New("team not found")
	ErrPRNotFound   = errors.New("PR not found")
)

func New(db *pgxpool.Pool, getter *txmpgx.CtxGetter) *Repository {
	return &Repository{
		db:     db,
		getter: getter,
	}
}

func (r *Repository) Close() {
	r.db.Close()
}

func (r *Repository) getConn(ctx context.Context) txmpgx.Tr {
	return r.getter.DefaultTrOrDB(ctx, r.db)
}
