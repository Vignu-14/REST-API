package routes

import (
	"net/http"
	"rest-api/internal/handlers"
)

func NewRouter(
	authHandler *handlers.AuthHandler,
	studentHandler *handlers.StudentHandler,
	teacherHandler *handlers.TeacherHandler,
	execHandler *handlers.ExecutiveHandler,
	jwtSecret string,
) *http.ServeMux {
	mux := http.NewServeMux()

	// Unified Login Route (Works for Student, Teacher, and Executive)
	mux.HandleFunc("POST /login", authHandler.Login)

	// Role-specific route groups
	RegisterStudentRoutes(mux, studentHandler, jwtSecret)
	RegisterTeacherRoutes(mux, teacherHandler, jwtSecret)
	RegisterExecutiveRoutes(mux, execHandler, jwtSecret)

	return mux
}
