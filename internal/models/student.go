package models

import "time"

// Student represents a row in the 'students' table
type Student struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Class        string    `json:"class"`
	EnrollmentNo string    `json:"enrollment_no"`
	CreatedAt    time.Time `json:"created_at"`
}

// StudentRegisterRequest is the JSON body for POST /students/register
type StudentRegisterRequest struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	Password     string `json:"password"`
	Class        string `json:"class"`
	EnrollmentNo string `json:"enrollment_no"`
}

// UpdateStudentRequest is the JSON body for PUT/PATCH /students/me
type UpdateStudentRequest struct {
	FullName     string `json:"full_name"`
	Class        string `json:"class"`
	EnrollmentNo string `json:"enrollment_no"`
}

// StudentProfile combines 'users' + 'students' for GET responses
type StudentProfile struct {
	UserID       string    `json:"user_id"`
	StudentID    string    `json:"student_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	Class        string    `json:"class"`
	EnrollmentNo string    `json:"enrollment_no"`
	CreatedAt    time.Time `json:"created_at"`
}