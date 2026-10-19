package services

import (
	"context"
	"fmt"
	"rest-api/internal/models"
	"rest-api/internal/repository"
	"rest-api/pkg/password"
	"strings"
)

type TeacherService interface {
	Register(ctx context.Context, req *models.TeacherRegisterRequest) (*models.TeacherProfile, error)
	GetProfile(ctx context.Context, userID string) (*models.TeacherProfile, error)
	UpdateProfile(ctx context.Context, userID string, req *models.UpdateTeacherRequest) error
	DeleteProfile(ctx context.Context, userID string) error
	// RBAC Read-Only: Teacher can view students belonging to their own class
	GetMyClassStudents(ctx context.Context, teacherUserID string) ([]models.StudentProfile, error)
}

type teacherService struct {
	teacherRepo repository.TeacherRepository
	studentRepo repository.StudentRepository
}

func NewTeacherService(teacherRepo repository.TeacherRepository, studentRepo repository.StudentRepository,) TeacherService {
	return &teacherService{
		teacherRepo: teacherRepo,
		studentRepo: studentRepo,
	}
}

func (s *teacherService) Register(ctx context.Context, req *models.TeacherRegisterRequest) (*models.TeacherProfile, error) {
	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.FullName) == "" || strings.TrimSpace(req.Password) == "" || strings.TrimSpace(req.Class) == "" || strings.TrimSpace(req.Subject) == ""{
		return nil, ErrInvalidInput
	}
	hashedPass, err := password.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Username: strings.TrimSpace(req.Username),
		Email: strings.TrimSpace(req.Email),
		FullName: strings.TrimSpace(req.FullName),
		PasswordHash: hashedPass,
		Role: models.RoleTeacher,
	}

	teacher := &models.Teacher{
		Class: req.Class,
		Subject: req.Subject,
	}

	err = s.teacherRepo.Create(ctx, user, teacher)
	if err != nil {
		return nil, fmt.Errorf("failed to create teacher: %w", err)
	}

	profile := &models.TeacherProfile{
		UserID: user.ID,
		TeacherID: teacher.ID,
		Username: user.Username,
		Email: user.Email,
		FullName: user.FullName,
		Class: teacher.Class,
		Subject: teacher.Subject,
		CreatedAt: teacher.CreatedAt,
	}
	return profile, nil
}
func (s *teacherService) GetProfile(ctx context.Context, userID string) (*models.TeacherProfile, error)
func (s *teacherService) UpdateProfile(ctx context.Context, userID string, req *models.UpdateTeacherRequest) error
func (s *teacherService) DeleteProfile(ctx context.Context, userID string) error
	// RBAC Read-Only: Teacher can view students belonging to their own class
func (s *teacherService) GetMyClassStudents(ctx context.Context, teacherUserID string) ([]models.StudentProfile, error)