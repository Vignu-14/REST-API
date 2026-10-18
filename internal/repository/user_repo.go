package repository

import (
	"context"
	"errors"
	"rest-api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	GetByIdentifier(ctx context.Context, identifier string) (*models.User, error)
}

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepo{
		db: db,
	}
}

// GetByIdentifier looks up a user by EITHER username OR email for Login
func (r *userRepo) GetByIdentifier(ctx context.Context, identifier string) (*models.User, error) {
	query := `
		SELECT id, username, email, full_name, password_hash, role, created_at, updated_at
		FROM users
		WHERE username = $1 OR email = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, identifier).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to query user by identifier: %w", err)
	}

	return &user, nil
}