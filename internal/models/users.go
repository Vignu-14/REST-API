package models

// true

import "time"

type Role string

const (
	RoleStudent   Role = "student"
	RoleTeacher   Role = "teacher"
	RoleExecutive Role = "executive"
)

// User represents a row in the 'users' table
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	PasswordHash string    `json:"-"` // "-" ensures password hash is NEVER sent in JSON responses
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// LoginRequest allows login via either username or email
type LoginRequest struct {
	Identifier string `json:"identifier"` // User can type their username OR email here
	Password   string `json:"password"`
}

// LoginResponse returns the JWT token and basic user info
type LoginResponse struct {
	Token string `json:"token"`
	Role  Role   `json:"role"`
	ID    string `json:"user_id"`
}
