package routes

import (
	"net/http"
	"rest-api/internal/handlers"
	"rest-api/internal/middleware"
	"rest-api/internal/models"
)

func RegisterTeacherRoutes(mux *http.ServeMux, h *handlers.TeacherHandler, jwtSecret string) {
	// Public Route
	mux.HandleFunc("POST /teachers/register", h.Register)

	// Protected Self-CRUD Routes (RoleTeacher only)
	mux.HandleFunc("GET /teachers/me", middleware.Authenticate(
		jwtSecret,
		middleware.RequireRole(h.GetProfile, models.RoleTeacher),
	))

	mux.HandleFunc("PUT /teachers/me", middleware.Authenticate(
		jwtSecret,
		middleware.RequireRole(h.UpdateProfile, models.RoleTeacher),
	))

	mux.HandleFunc("DELETE /teachers/me", middleware.Authenticate(
		jwtSecret,
		middleware.RequireRole(h.DeleteProfile, models.RoleTeacher),
	))

	// Protected RBAC Read-Only Route: View students in this teacher's class
	mux.HandleFunc("GET /teachers/my-students", middleware.Authenticate(
		jwtSecret,
		middleware.RequireRole(h.GetMyClassStudents, models.RoleTeacher),
	))
}
