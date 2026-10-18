package repository

import (
	"context"
	"rest-api/internal/models"
)

type TeacherRepository interface {
	Create(ctx context.Context, user *models.User, teacher *models.Teacher) error
	GetByUserID(ctx context.Context, userID string) (*models.TeacherProfile, error)
	Update(ctx context.Context, userID string, req *models.UpdateTeacherRequest) error
	Delete(ctx context.Context, userID string) error
	// Used by Executives
	GetAll(ctx context.Context) ([]models.TeacherProfile, error)
}
