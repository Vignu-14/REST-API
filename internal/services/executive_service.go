package services

import (
	"context"
	"fmt"
	"rest-api/internal/models"
	"rest-api/internal/repository"
	"rest-api/pkg/password"
	"strings"
)

type ExecutiveService interface {
	Register(ctx context.Context, req *models.ExecutiveRegisterRequest) (*models.ExecutiveProfile, error)
	GetProfile(ctx context.Context, userID string) (*models.ExecutiveProfile, error)
	UpdateProfile(ctx context.Context, userID string, req *models.UpdateExecutiveRequest) error
	DeleteProfile(ctx context.Context, userID string) error
	// RBAC Read-Only: Executive can view all records across the school
	GetAllStudents(ctx context.Context) ([]models.StudentProfile, error)
	GetAllTeachers(ctx context.Context) ([]models.TeacherProfile, error)
	GetAllExecutives(ctx context.Context) ([]models.ExecutiveProfile, error)
}

type executiveService struct {
	exceRepo    repository.ExecutiveRepository
	teacherRepo repository.TeacherRepository
	studentRepo repository.StudentRepository
}

func NewExecutiveService(execRepo repository.ExecutiveRepository, teacherRepo repository.TeacherRepository, studentRepo repository.StudentRepository) ExecutiveService {
	return &executiveService{
		exceRepo:    execRepo,
		teacherRepo: teacherRepo,
		studentRepo: studentRepo,
	}
}

func (s *executiveService) Register(ctx context.Context, req *models.ExecutiveRegisterRequest) (*models.ExecutiveProfile, error) {
	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.FullName) == "" || strings.TrimSpace(req.Password) == "" || strings.TrimSpace(req.Department) == "" {
		return nil, ErrInvalidInput
	}

	hashedPass, err := password.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Username:     strings.TrimSpace(req.Username),
		Email:        strings.TrimSpace(req.Email),
		FullName:     strings.TrimSpace(req.FullName),
		PasswordHash: hashedPass,
		Role:         models.RoleExecutive,
	}

	exec := &models.Executive{
		Department: strings.TrimSpace(req.Department),
	}

	err = s.exceRepo.Create(ctx, user, exec)
	if err != nil {
		return nil, fmt.Errorf("failed to crete executive: %w", err)
	}

	profile := &models.ExecutiveProfile{
		UserID: user.ID,
		ExecutiveID: exec.ID,
		Username: user.Username,
		Email: user.Email,
		FullName: user.FullName,
		Department: exec.Department,
		CreatedAt: exec.CreatedAt,
	}
	return profile, nil
}

func (s *executiveService) GetProfile(ctx context.Context, userID string) (*models.ExecutiveProfile, error) {
	return s.exceRepo.GetByUserID(ctx, userID)
}

func (s *executiveService) UpdateProfile(ctx context.Context, userID string, req *models.UpdateExecutiveRequest) error {
	
}
func (s *executiveService) DeleteProfile(ctx context.Context, userID string) error

// RBAC Read-Only: Executive can view all records across the school
func (s *executiveService) GetAllStudents(ctx context.Context) ([]models.StudentProfile, error)
func (s *executiveService) GetAllTeachers(ctx context.Context) ([]models.TeacherProfile, error)
func (s *executiveService) GetAllExecutives(ctx context.Context) ([]models.ExecutiveProfile, error)
