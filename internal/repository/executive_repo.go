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

// 3. Update modifies both 'users' (full_name) and 'executives' (department)
func (r *executiveRepo) Update(ctx context.Context, userID string, req *models.UpdateExecutiveRequest) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	userUpdate := `
		UPDATE users
		SET full_name = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND role = 'executive'
	`
	cmdTag, err := tx.Exec(ctx, userUpdate, req.FullName, userID)
	if err != nil {
		return fmt.Errorf("failed to update user table: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrExecutiveNotFound
	}

	execUpdate := `
		UPDATE executives
		SET department = $1
		WHERE user_id = $2
	`
	_, err = tx.Exec(ctx, execUpdate, req.Department, userID)
	if err != nil {
		return fmt.Errorf("failed to update executive table: %w", err)
	}

	return tx.Commit(ctx)
}

// 4. Delete removes the user row; ON DELETE CASCADE deletes the executive row
func (r *executiveRepo) Delete(ctx context.Context, userID string) error {
	query := `DELETE FROM users WHERE id = $1 AND role = 'executive'`
	cmdTag, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete executive: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrExecutiveNotFound
	}
	return nil
}

// 5. GetAll returns all executives
func (r *executiveRepo) GetAll(ctx context.Context) ([]models.ExecutiveProfile, error) {
	query := `
		SELECT u.id, e.id, u.username, u.email, u.full_name, e.department, e.created_at
		FROM users u
		INNER JOIN executives e ON u.id = e.user_id
		ORDER BY e.department ASC, u.full_name ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all executives: %w", err)
	}
	defer rows.Close()

	var executives []models.ExecutiveProfile
	for rows.Next() {
		var p models.ExecutiveProfile
		if err := rows.Scan(
			&p.UserID, &p.ExecutiveID, &p.Username, &p.Email,
			&p.FullName, &p.Department, &p.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan executive row: %w", err)
		}
		executives = append(executives, p)
	}

	return executives, rows.Err()
}