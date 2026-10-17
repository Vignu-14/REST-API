package models

import "github.com/google/uuid"

// Teacher is the raw database model for `teachers` table
type Teacher struct {
	ID      uuid.UUID `json:"id" db:"id"`
	UserID  uuid.UUID `json:"user_id" db:"user_id"`
	Class   string    `json:"class" db:"class"`
	Subject string    `json:"subject" db:"subject"`
}

// TeacherProfileResponse embeds UserResponse
type TeacherProfileResponse struct {
	UserResponse
	Class   string `json:"class"`
	Subject string `json:"subject"`
}

// CreateTeacherRequest is the incoming payload from an executive
type CreateTeacherRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	Class    string `json:"class" binding:"required"`
	Subject  string `json:"subject" binding:"required"`
}