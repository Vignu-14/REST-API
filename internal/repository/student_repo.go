package repository

import (
	"context"
	"fmt"
	"rest-api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
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

type studentRepo struct {
	db *pgxpool.Pool
}

func NewStudentRepository(db *pgxpool.Pool) StudentRepository {
	return &studentRepo{
		db: db,
	}
}

func (r *studentRepo) Create(ctx context.Context, user *models.User, student *models.Student) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	userQuery := `
			INSERT INTO users(username, email, full_name, password_hash, role)
			VALUES($1, $2, $3, $4, $5)
			RETURNING(id, created_at, updated_at)
	`
	err = tx.QueryRow(ctx, userQuery,
		user.Username,
		user.Email,
		user.FullName,
		user.PasswordHash,
		user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert to user: %w", err)
	}

	student.UserID = user.ID
	studentQuery := `
			INSERT INTO students(user_id, class, enrollment_no)
			VAlUES($1, $2, $3)
			RETURNING(id, created_at)
	`
	err = r.db.QueryRow(ctx, studentQuery,
			student.UserID,
			student.Class,
			student.EnrollmentNo,
	).Scan(&student.ID, &student.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert student")
	}
	
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction")
	}

	return ctx.Err()
}

func (r *studentRepo) GetByUserID(ctx context.Context, userID string) (*models.StudentProfile, error)

func (r *studentRepo) Update(ctx context.Context, userID string, req *models.UpdateStudentRequest) error

func (r *studentRepo) Delete(ctx context.Context, userID string) error

func (r *studentRepo) GetByClass(ctx context.Context, class string) ([]models.StudentProfile, error)

func (r *studentRepo) GetAll(ctx context.Context) ([]models.StudentProfile, error)
