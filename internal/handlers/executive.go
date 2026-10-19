package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"rest-api/internal/middleware"
	"rest-api/internal/models"
	"rest-api/internal/repository"
	"rest-api/internal/services"
)

type ExecutiveHandler struct {
	execService services.ExecutiveService
}

func NewExecutiveHandler(execService services.ExecutiveService) *ExecutiveHandler {
	return &ExecutiveHandler{
		execService: execService,
	}
}

// 1. Register handles POST /executives/register (Public route)
func (h *ExecutiveHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.ExecutiveRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	profile, err := h.execService.Register(r.Context(), &req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusConflict, "failed to register executive (username or email may already exist)")
		return
	}

	writeJSON(w, http.StatusCreated, profile)
}

// 2. GetProfile handles GET /executives/me (Protected: RoleExecutive)
func (h *ExecutiveHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	profile, err := h.execService.GetProfile(r.Context(), userID)
	if err != nil {
		if errors.Is(err, repository.ErrExecutiveNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to fetch profile")
		return
	}

	writeJSON(w, http.StatusOK, profile)
}

// 3. UpdateProfile handles PUT /executives/me (Protected: RoleExecutive)
func (h *ExecutiveHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req models.UpdateExecutiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if err := h.execService.UpdateProfile(r.Context(), userID, &req); err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, repository.ErrExecutiveNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update profile")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "executive profile updated successfully"})
}

// 4. DeleteProfile handles DELETE /executives/me (Protected: RoleExecutive)
func (h *ExecutiveHandler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.execService.DeleteProfile(r.Context(), userID); err != nil {
		if errors.Is(err, repository.ErrExecutiveNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete profile")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "executive account deleted successfully"})
}

// 5. GetAllStudents handles GET /executives/students (Protected: RoleExecutive)
func (h *ExecutiveHandler) GetAllStudents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	students, err := h.execService.GetAllStudents(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch all students")
		return
	}

	writeJSON(w, http.StatusOK, students)
}

// 6. GetAllTeachers handles GET /executives/teachers (Protected: RoleExecutive)
func (h *ExecutiveHandler) GetAllTeachers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	teachers, err := h.execService.GetAllTeachers(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch all teachers")
		return
	}

	writeJSON(w, http.StatusOK, teachers)
}

// 7. GetAllExecutives handles GET /executives/all (Protected: RoleExecutive)
func (h *ExecutiveHandler) GetAllExecutives(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	executives, err := h.execService.GetAllExecutives(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch all executives")
		return
	}

	writeJSON(w, http.StatusOK, executives)
}
