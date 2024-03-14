package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB = pgxpool.Pool

func New(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, dsn)
}
