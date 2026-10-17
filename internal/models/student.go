package models

import "github.com/google/uuid"

// Student is the raw database model for `students` table
type Student struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	Class        string    `json:"class" db:"class"`
	EnrollmentNo string    `json:"enrollment_no" db:"enrollment_no"`
}

// StudentProfileResponse embeds UserResponse automatically
// When converted to JSON, all UserResponse fields appear at top level
type StudentProfileResponse struct {
	UserResponse
	Class        string `json:"class"`
	EnrollmentNo string `json:"enrollment_no"`
}

// RegisterStudentRequest is the incoming payload
type RegisterStudentRequest struct {
	Username     string `json:"username" binding:"required,min=3"`
	Email        string `json:"email" binding:"required,email"`
	FullName     string `json:"full_name" binding:"required"`
	Password     string `json:"password" binding:"required,min=8"`
	Class        string `json:"class" binding:"required"`
	EnrollmentNo string `json:"enrollment_no" binding:"required"`
}
