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

func NewTeacherRepository(db *pgxpool.Pool) TeacherRepository {
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
		RETURNING id, created_at
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
			return nil, ErrTeacherNotFound
		}
		return nil, fmt.Errorf("failed to get teacher by user_id: %w", err)
	}
	return &profile, nil
}

func (t *teacherRepo) Update(ctx context.Context, userID string, req *models.UpdateTeacherRequest) error {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	userUpdate := `
		UPDATE users
		SET full_name = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND role = 'teacher'
	`
	cmdTag, err := tx.Exec(ctx, userUpdate, req.FullName, userID)
	if err != nil {
		return fmt.Errorf("failed to update user table: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrTeacherNotFound
	}

	teacherUpdate := `
		UPDATE teachers
		SET class = $1, subject = $2
		WHERE user_id = $3
	`
	_, err = tx.Exec(ctx, teacherUpdate, req.Class, req.Subject, userID)
	if err != nil {
		return fmt.Errorf("failed to update teacher table: %w", err)
	}
	return tx.Commit(ctx)
}

func (t *teacherRepo) Delete(ctx context.Context, userID string) error {
	query := `DELETE FROM users WHERE id = $1 AND role = 'teacher'`
	cmdTag, err := t.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete teacher: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrTeacherNotFound
	}
	return nil
}

// Used by Executives
func (t *teacherRepo) GetAll(ctx context.Context) ([]models.TeacherProfile, error) {
	query := `
		SELECT u.id, t.id, u.username, u.email, u.full_name, t.class, t.subject, t.created_at
		FROM users u
		INNER JOIN teachers t ON u.id = t.user_id
		ORDER BY t.class ASC, u.full_name ASC
	`

	rows, err := t.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all teachers: %w", err)
	}
	defer rows.Close()

	var teachers []models.TeacherProfile
	for rows.Next() {
		var p models.TeacherProfile
		if err := rows.Scan(
			&p.UserID, &p.TeacherID, &p.Username, &p.Email,
			&p.FullName, &p.Class, &p.Subject, &p.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan teacher row: %w", err)
		}
		teachers = append(teachers, p)
	}
	return teachers, rows.Err()
}
