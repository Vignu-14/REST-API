package models

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleStudent   Role = "student"
	RoleTeacher   Role = "teacher"
	RoleExecutive Role = "executive"
)

// User is the raw database model for the `users` table
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	FullName     string    `json:"full_name" db:"full_name"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         Role      `json:"role" db:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// UserResponse is the safe public identity (no password)
// Other structs will EMBED this to avoid repeating fields
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginRequest is shared across everyone
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}