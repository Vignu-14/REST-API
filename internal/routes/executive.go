package routes

import (
	"net/http"
	"rest-api/internal/handlers"
	"rest-api/internal/middleware"
	"rest-api/internal/models"
)

func RegisterExecutiveRoutes(mux *http.ServeMux, h *handlers.ExecutiveHandler, jwtSecret string) {
	// Public Route
	mux.HandleFunc("POST /executives/register", h.Register)

	// Protected Self-CRUD Routes (RoleExecutive only)
	mux.HandleFunc("GET /executives/me", middleware.Authenticate(
		jwtSecret,
		middleware.RequireRole(h.GetProfile, models.RoleExecutive),
	))

	mux.HandleFunc("PUT /executives/me", middleware.Authenticate(
		jwtSecret,
		middleware.RequireRole(h.UpdateProfile, models.RoleExecutive),
	))

	mux.HandleFunc("DELETE /executives/me", middleware.Authenticate(
		jwtSecret,
		middleware.RequireRole(h.DeleteProfile, models.RoleExecutive),
	))

	// Protected RBAC Read-All Routes (RoleExecutive only)
	mux.HandleFunc("GET /executives/students", middleware.Authenticate(
		jwtSecret,
		middleware.RequireRole(h.GetAllStudents, models.RoleExecutive),
	))

	mux.HandleFunc("GET /executives/teachers", middleware.Authenticate(
		jwtSecret,
		middleware.RequireRole(h.GetAllTeachers, models.RoleExecutive),
	))

	mux.HandleFunc("GET /executives/all", middleware.Authenticate(
		jwtSecret,
		middleware.RequireRole(h.GetAllExecutives, models.RoleExecutive),
	))
}
