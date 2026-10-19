package services

import (
	"context"
	"rest-api/internal/models"
	"rest-api/internal/repository"
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

func (s *teacherService) Register(ctx context.Context, req *models.TeacherRegisterRequest) (*models.TeacherProfile, error)
func (s *teacherService) GetProfile(ctx context.Context, userID string) (*models.TeacherProfile, error)
func (s *teacherService) UpdateProfile(ctx context.Context, userID string, req *models.UpdateTeacherRequest) error
func (s *teacherService) DeleteProfile(ctx context.Context, userID string) error
	// RBAC Read-Only: Teacher can view students belonging to their own class
func (s *teacherService) GetMyClassStudents(ctx context.Context, teacherUserID string) ([]models.StudentProfile, error)