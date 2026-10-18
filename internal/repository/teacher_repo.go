package repository

import (
	"context"
	"errors"
	"fmt"
	"rest-api/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrTeacherNotFound = errors.New("teacher not found")

type TeacherRepository interface {
	Create(ctx context.Context, user *models.User, teacher *models.Teacher) error
	GetByUserID(ctx context.Context, userID string) (*models.TeacherProfile, error)
	Update(ctx context.Context, userID string, req *models.UpdateTeacherRequest) error
	Delete(ctx context.Context, userID string) error
	// Used by Executives
	GetAll(ctx context.Context) ([]models.TeacherProfile, error)
}

type teacherRepo struct {
	db *pgxpool.Pool
}

func NewTeacherRepo(db *pgxpool.Pool) TeacherRepository {
	return &teacherRepo{
		db: db,
	}
}

func (t *teacherRepo) Create(ctx context.Context, user *models.User, teacher *models.Teacher) error{
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to create teacher")
	}
	defer tx.Rollback(ctx)

	userQuery := `
		INSERT INTO users(username, email, full_name, password_hash, role)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(ctx, userQuery, user.Username, user.Email, user.FullName, user.PasswordHash, user.Role).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create teacher: %w", err)
	}

	teacher.UserID = user.ID
	teacherQuery := `
		INSERT INTO teachers(user_id, class, subject)
		VALUES($1, $2, $3)
		RETURNING(id, created_at)
	`
	err = tx.QueryRow(ctx, teacherQuery, teacher.UserID, teacher.Class, teacher.Subject).Scan(&teacher.ID, &teacher.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create teacher: %w", err)
	}
	
	if err = tx.Commit(ctx); err != nil{
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}


func (t *teacherRepo) GetByUserID(ctx context.Context, userID string) (*models.TeacherProfile, error){
	query := `
		SELECT u.id, t.id, u.username, u.email, u.full_name, t.class, t.subject, t.created_at
		FROM users u
		INNER JOIN teachers t ON u.id = t.user_id
		WHERE u.id = $1 AND u.role = 'teacher'
	`
	var profile models.TeacherProfile
	err := t.db.QueryRow(ctx, query, userID).Scan(&profile.UserID, &profile.TeacherID, &profile.Username, &profile.Email, &profile.FullName, &profile.Class, &profile.Subject, &profile.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			return nil, ErrStudentNotFound
		}
		return nil, fmt.Errorf("failed to get teacher by user_id: %w", err)
	}
	return &profile, nil
}

func (t *teacherRepo) Update(ctx context.Context, userID string, req *models.UpdateTeacherRequest) error
func (t *teacherRepo) Delete(ctx context.Context, userID string) error

// Used by Executives
func (t *teacherRepo) GetAll(ctx context.Context) ([]models.TeacherProfile, error)
