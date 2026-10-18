package repository

import (
	"context"
	"errors"
	"fmt"
	"rest-api/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrExecutiveNotFound = errors.New("executive not found")

type ExecutiveRepository interface {
	Create(ctx context.Context, user *models.User, exec *models.Executive) error
	GetByUserID(ctx context.Context, userID string) (*models.ExecutiveProfile, error)
	Update(ctx context.Context, userID string, req *models.UpdateExecutiveRequest) error
	Delete(ctx context.Context, userID string) error
	GetAll(ctx context.Context) ([]models.ExecutiveProfile, error)
}

type executiveRepo struct {
	db *pgxpool.Pool
}

func NewExecutiveRepository(db *pgxpool.Pool) ExecutiveRepository {
	return &executiveRepo{
		db: db,
	}
}

// 1. Create inserts into 'users' and 'executives' inside an atomic transaction
func (r *executiveRepo) Create(ctx context.Context, user *models.User, exec *models.Executive) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	userQuery := `
		INSERT INTO users (username, email, full_name, password_hash, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(ctx, userQuery,
		user.Username,
		user.Email,
		user.FullName,
		user.PasswordHash,
		models.RoleExecutive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	exec.UserID = user.ID
	execQuery := `
		INSERT INTO executives (user_id, department)
		VALUES ($1, $2)
		RETURNING id, created_at
	`
	err = tx.QueryRow(ctx, execQuery,
		exec.UserID,
		exec.Department,
	).Scan(&exec.ID, &exec.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert executive: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// 2. GetByUserID fetches an executive's joined profile
func (r *executiveRepo) GetByUserID(ctx context.Context, userID string) (*models.ExecutiveProfile, error) {
	query := `
		SELECT u.id, e.id, u.username, u.email, u.full_name, e.department, e.created_at
		FROM users u
		INNER JOIN executives e ON u.id = e.user_id
		WHERE u.id = $1 AND u.role = 'executive'
	`

	var profile models.ExecutiveProfile
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&profile.UserID,
		&profile.ExecutiveID,
		&profile.Username,
		&profile.Email,
		&profile.FullName,
		&profile.Department,
		&profile.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrExecutiveNotFound
		}
		return nil, fmt.Errorf("failed to get executive by user_id: %w", err)
	}

	return &profile, nil
}
