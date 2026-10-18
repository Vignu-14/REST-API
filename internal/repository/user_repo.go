package repository

import (
	"context"
	"errors"
	"rest-api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	GetByIdentifier(ctx context.Context, identifier string) (*models.User, error)
}

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepo{
		db: db,
	}
}
