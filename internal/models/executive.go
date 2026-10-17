package models

import "github.com/google/uuid"

// Executive is the raw database model for the `executives` table
type Executive struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	Department string    `json:"department" db:"department"`
}

// ExecutiveProfileResponse embeds UserResponse automatically
// Keeps the JSON response flat with all base identity fields included
type ExecutiveProfileResponse struct {
	UserResponse
	Department string `json:"department"`
}

// CreateExecutiveRequest is used when bootstrapping or creating admin accounts
type CreateExecutiveRequest struct {
	Username   string `json:"username" binding:"required,min=3"`
	Email      string `json:"email" binding:"required,email"`
	FullName   string `json:"full_name" binding:"required"`
	Password   string `json:"password" binding:"required,min=8"`
	Department string `json:"department" binding:"required"`
}