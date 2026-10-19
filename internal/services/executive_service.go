package services

import (
	"context"
	"rest-api/internal/models"
	"rest-api/internal/repository"
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
	exceRepo repository.ExecutiveRepository
	teacherRepo repository.StudentRepository
	studentRepo repository.StudentRepository
}

func NewExecutiveService(execRepo repository.ExecutiveRepository, teacherRepo repository.TeacherRepository, studentRepo repository.StudentRepository) ExecutiveService {
	return &executiveService{
		exceRepo: execRepo,
		teacherRepo: teacherRepo,
		studentRepo: studentRepo,
	}
}

func (s *executiveService) Register(ctx context.Context, req *models.ExecutiveRegisterRequest) (*models.ExecutiveProfile, error) {
	
} 

func (s *executiveService) GetProfile(ctx context.Context, userID string) (*models.ExecutiveProfile, error)
func (s *executiveService) UpdateProfile(ctx context.Context, userID string, req *models.UpdateExecutiveRequest) error
func (s *executiveService) DeleteProfile(ctx context.Context, userID string) error

	// RBAC Read-Only: Executive can view all records across the school
func (s *executiveService) GetAllStudents(ctx context.Context) ([]models.StudentProfile, error)
func (s *executiveService) GetAllTeachers(ctx context.Context) ([]models.TeacherProfile, error)
func (s *executiveService) GetAllExecutives(ctx context.Context) ([]models.ExecutiveProfile, error)