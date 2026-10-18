package models

import "time"

// Executive represents a row in the 'executives' table
type Executive struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Department string    `json:"department"`
	CreatedAt  time.Time `json:"created_at"`
}

// ExecutiveRegisterRequest is the JSON body for POST /executives/register
type ExecutiveRegisterRequest struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	FullName   string `json:"full_name"`
	Password   string `json:"password"`
	Department string `json:"department"`
}

// UpdateExecutiveRequest is the JSON body for PUT/PATCH /executives/me
type UpdateExecutiveRequest struct {
	FullName   string `json:"full_name"`
	Department string `json:"department"`
}

// ExecutiveProfile combines 'users' + 'executives' for GET responses
type ExecutiveProfile struct {
	UserID      string    `json:"user_id"`
	ExecutiveID string    `json:"executive_id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	FullName    string    `json:"full_name"`
	Department  string    `json:"department"`
	CreatedAt   time.Time `json:"created_at"`
}