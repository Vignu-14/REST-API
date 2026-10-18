package repository

// true

import (
	"context"
	"errors"
	"fmt"
	"rest-api/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrStudentNotFound = errors.New("student not found")

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

func (r *studentRepo) GetByUserID(ctx context.Context, userID string) (*models.StudentProfile, error) {
	query := `
		SELECT u.id, s.id, u.username, u.email, u.full_name, s.class, s.enrollment_no, s.created_at
		FROM users u
		INNER JOIN students s ON u.id = s.user_id
		WHERE u.id = $1 AND u.role = 'student'
	`
	var profile models.StudentProfile
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&profile.UserID,
		&profile.StudentID,
		&profile.Username,
		&profile.Email,
		&profile.FullName,
		&profile.Class,
		&profile.EnrollmentNo,
		&profile.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			return nil, ErrStudentNotFound
		}
		return nil, fmt.Errorf("failed to student by user_id: %w", err)
	}
	return &profile, nil
}
// Update modifies both 'users' (full_name) and 'students' (class, enrollment_no)
func (r *studentRepo) Update(ctx context.Context, userID string, req *models.UpdateStudentRequest) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	userUpdate := `
		UPDATE users
		SET full_name = $1 , UPDATED_AT = CURRENT_TIMESTAMP
		WHERE id = $2 AND role = 'student'
	`
	cmdTag, err := tx.Exec(ctx, userUpdate, req.FullName, userID)
	if err != nil {
		return fmt.Errorf("failed to update user table: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrStudentNotFound
	}

	studentUpdate := `
		UPDATE students
		SET class = $1, enrollment_no = $2
		WHERE user_id = $3
	`
	cmdTag, err = tx.Exec(ctx, studentUpdate, req.Class, req.EnrollmentNo, userID)
	if err != nil {
		return fmt.Errorf("failed to update student table: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrStudentNotFound
	}
	if err = tx.Commit(ctx); err != nil{
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// return tx.Commit(ctx)  // we can also return this because it's returns err 
	return nil
}

func (r *studentRepo) Delete(ctx context.Context, userID string) error {
	query := `DELETE FROM user WHERE id = $1 AND role = 'student'`
	cmdTag, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete student: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrStudentNotFound
	}
	return nil
}

func (r *studentRepo) GetByClass(ctx context.Context, class string) ([]models.StudentProfile, error){
	query := `
		SELECT u.id, s.id, u.username, u.email, u.full_name, s.class, s.enrollment_no, s.created_at
		FROM users u
		INNER JOIN students s ON u.id = s.user_id
		WHERE s.class = $1
		ORDER BY u.full_name ASC
	`

	rows, err := r.db.Query(ctx, query, class)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch student by class: %w", err)
	}
	defer rows.Close()

	var students []models.StudentProfile
	for rows.Next() {
		var u models.StudentProfile
		if err = rows.Scan(&u.UserID, &u.StudentID, &u.Username, &u.Email, &u.FullName, &u.Class, &u.EnrollmentNo, u.CreatedAt); err != nil{
			return nil, fmt.Errorf("failed to scan student row: %w", err)
		}
		students = append(students, u)
	}
	return students, rows.Err()
}

func (r *studentRepo) GetAll(ctx context.Context) ([]models.StudentProfile, error) {
	query := `
		SELECT u.id, s.id, u.username, u.email, u.full_name, s.class, s.enrollment_no, s.created_at
		FROM users u
		INNER JOIN students s ON u.id = s.user_id
		ORDER BY s.class ASC, u.full_name ASC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get student: %w", err)
	}
	defer rows.Close()

	var students []models.StudentProfile
	for rows.Next() {
		var u models.StudentProfile
		if err = rows.Scan(&u.UserID, &u.StudentID, &u.Username, &u.Email, &u.FullName, &u.Class, &u.EnrollmentNo, &u.CreatedAt); err != nil{
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		students = append(students, u)
	}
	return students, rows.Err()
}
