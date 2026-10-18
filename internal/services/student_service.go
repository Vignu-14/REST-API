package services

import (
	"context"
	"errors"
	"rest-api/internal/models"
	"rest-api/internal/repository"
)

var ErrInvalidInput = errors.New("all required fields must be provided")

type StudentService interface {
	Register(ctx context.Context, req *models.StudentRegisterRequest) (*models.StudentProfile, error)
	GetProfile(ctx context.Context, userID string) (*models.StudentProfile, error)
	UpdateProfile(ctx context.Context, userID string, req *models.UpdateStudentRequest) error
	DeleteProfile(ctx context.Context, userID string) error
}

type studentService struct{
	studentRepo repository.StudentRepository
}

func NewStudentService(studentRepo repository.StudentRepository) StudentService {
	return &studentService{
		studentRepo: studentRepo,
	}
}