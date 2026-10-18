package repository

import (
	"context"
	"rest-api/internal/models"
)

type StudentRepository interface {
	Create(ctx context.Context, user *models.User, student *models.Student) error
	GetByUserID(ctx context.Context, userID string) (*models.StudentProfile, error)
	Update(ctx context.Context, userID string, req *models.UpdateStudentRequest) error
	Delete(ctx context.Context, userID string) error
	// Used by Teachers (their class) and Executives (all classes)
	GetByClass(ctx context.Context, class string) ([]models.StudentProfile, error)
	GetAll(ctx context.Context) ([]models.StudentProfile, error)
}
