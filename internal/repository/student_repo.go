package repository

import "github.com/jackc/pgx/v5/pgxpool"

type RepoUser struct {
	pool *pgxpool.Pool
}

func NewRepoUser(pool *pgxpool.Pool) *RepoUser{
	return &RepoUser{pool: pool}
}

func DbCreateStudent(username, email, fullname, password, class, enroll_no string) error
