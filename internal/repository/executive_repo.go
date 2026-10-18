package repository

import (
	"context"
	"errors"
	"rest-api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrExecutiveNotFound = errors.New("executive not found")

type ExecutiveRepository interface {
	Create(ctx context.Context, user *models.User, exec *models.Executive) error
	GetByUserID(ctx context.Context, userID string) (*models.ExecutiveProfile, error)
	Update(ctx context.Context, userID string, req *models.UpdateExecutiveRequest) error
	Delete(ctx context.Context, userID string) error
	GetAll(ctx context.Context) ([]models.ExecutiveProfile, error)
}

type executiveRepo struct {
	db *pgxpool.Pool
}

func NewExecutiveRepository(db *pgxpool.Pool) ExecutiveRepository {
	return &executiveRepo{
		db: db,
	}
}
