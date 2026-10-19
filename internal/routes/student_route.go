package routes

import (
	"net/http"
	"rest-api/internal/handlers"
	"rest-api/internal/middleware"
	"rest-api/internal/models"
)

func RegisterStudentRoutes(mux *http.ServeMux, h *handlers.StudentHandler, jwtSecret string) {

	// public route: anyone can register as a student
	mux.HandleFunc("POST /students/register", h.Register)

	// Potected route: only authenticated user with role = 'student'
	mux.HandleFunc("GET /students/me", middleware.Authenticate(jwtSecret, middleware.RequireRole(h.GetProfile, models.RoleStudent)))

	mux.HandleFunc("PUT /students/me", middleware.Authenticate(jwtSecret, middleware.RequireRole(h.UpdateProfile, models.RoleStudent)))

	mux.HandleFunc("DELETE /students/me", middleware.Authenticate(jwtSecret,middleware.RequireRole(h.DeleteProfile, models.RoleStudent),
	))

}