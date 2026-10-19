package services

import (
	"context"
	"errors"
	"fmt"
	"rest-api/internal/models"
	"rest-api/internal/repository"
	"rest-api/pkg/password"
	"strings"
)

var ErrInvalidInput = errors.New("all required fields must be provided")

type StudentService interface {
	Register(ctx context.Context, req *models.StudentRegisterRequest) (*models.StudentProfile, error)
	GetProfile(ctx context.Context, userID string) (*models.StudentProfile, error)
	UpdateProfile(ctx context.Context, userID string, req *models.UpdateStudentRequest) error
	DeleteProfile(ctx context.Context, userID string) error
}

type studentService struct {
	studentRepo repository.StudentRepository
}

func NewStudentService(studentRepo repository.StudentRepository) StudentService {
	return &studentService{
		studentRepo: studentRepo,
	}
}

func (s *studentService) Register(ctx context.Context, req *models.StudentRegisterRequest) (*models.StudentProfile, error) {
	if strings.TrimSpace(req.Username) == "" ||
		strings.TrimSpace(req.Email) == "" ||
		strings.TrimSpace(req.FullName) == "" ||
		strings.TrimSpace(req.Password) == "" ||
		strings.TrimSpace(req.Class) == "" {
		return nil, ErrInvalidInput
	}

	hashedPwd, err := password.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w",  err)
	}

	user := &models.User{
		Username: strings.TrimSpace(req.Username),
		Email: strings.TrimSpace(req.Email),
		FullName: strings.TrimSpace(req.FullName),
		PasswordHash: hashedPwd,
		Role: models.RoleStudent,
	}

	student := &models.Student{
		Class: strings.TrimSpace(req.Class),
		EnrollmentNo: strings.TrimSpace(req.EnrollmentNo),
	}
	err = s.studentRepo.Create(ctx, user, student)
	if err != nil {
		return nil, fmt.Errorf("failed to create student: %w", err)
	}

	profile := &models.StudentProfile{
		UserID: user.ID,
		StudentID: student.ID,
		Username: user.Username,
		Email: user.Email,
		FullName: user.FullName,
		Class: student.Class,
		EnrollmentNo: student.EnrollmentNo,
		CreatedAt: user.CreatedAt,
	}
	return profile, nil
}
func (s *studentService) GetProfile(ctx context.Context, userID string) (*models.StudentProfile, error)
func (s *studentService) UpdateProfile(ctx context.Context, userID string, req *models.UpdateStudentRequest) error
func (s *studentService) DeleteProfile(ctx context.Context, userID string) error
