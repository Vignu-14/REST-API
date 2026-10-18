package repository

import (
	"context"
	"rest-api/internal/models"
)

type ExecutiveRepository interface {
	Create(ctx context.Context, user *models.User, exec *models.Executive) error
	GetByUserID(ctx context.Context, userID string) (*models.ExecutiveProfile, error)
	Update(ctx context.Context, userID string, req *models.UpdateExecutiveRequest) error
	Delete(ctx context.Context, userID string) error
	GetAll(ctx context.Context) ([]models.ExecutiveProfile, error)
}
