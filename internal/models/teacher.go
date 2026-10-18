package models

// true

import "time"

// Teacher represents a row in the 'teachers' table
type Teacher struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Class     string    `json:"class"`
	Subject   string    `json:"subject"`
	CreatedAt time.Time `json:"created_at"`
}

// TeacherRegisterRequest is the JSON body for POST /teachers/register
type TeacherRegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
	Class    string `json:"class"`
	Subject  string `json:"subject"`
}

// UpdateTeacherRequest is the JSON body for PUT/PATCH /teachers/me
type UpdateTeacherRequest struct {
	FullName string `json:"full_name"`
	Class    string `json:"class"`
	Subject  string `json:"subject"`
}

// TeacherProfile combines 'users' + 'teachers' for GET responses
type TeacherProfile struct {
	UserID    string    `json:"user_id"`
	TeacherID string    `json:"teacher_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Class     string    `json:"class"`
	Subject   string    `json:"subject"`
	CreatedAt time.Time `json:"created_at"`
}