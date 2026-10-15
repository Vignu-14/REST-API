package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPostgres(ctx context.Context, DbUrl string) (*pgxpool.Pool, error) {
	Pool, err := pgxpool.New(ctx, DbUrl)
	if err != nil {
		return nil, errors.New("failed to connect to database")
	}
	return Pool, nil
}
