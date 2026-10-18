package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RepoUser struct {
	pool *pgxpool.Pool
}

func NewRepoUser(pool *pgxpool.Pool) *RepoUser {
	return &RepoUser{pool: pool}
}

func (repo RepoUser) DbCreateStudent(ctx context.Context, username, email, fullname, password, class, enroll_no string) (int, error) {
	return 0, errors.New("")
}
