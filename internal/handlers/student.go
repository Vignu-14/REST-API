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

type StudentHandler struct {
	studentService services.StudentService
}

func NewStudentHandler(studentService services.StudentService) *StudentHandler {
	return &StudentHandler{
		studentService: studentService,
	}
}

// POST students/register	(public route)
func (h *StudentHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.StudentRegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	profile, err := h.studentService.Register(r.Context(), &req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusConflict, "failed to register student (username, email, or enrollment_no may already exist)")
		return
	}
	writeJSON(w, http.StatusCreated, profile)
}

// GET students/me (protected: RoleStudent)
func (h *StudentHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	profile, err := h.studentService.GetProfile(r.Context(), userID)
	if err != nil {
		if errors.Is(err, repository.ErrStudentNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to fetch profile")
		return
	}
	writeJSON(w, http.StatusOK, profile)
}

// PUT students/me	 (protected: RoleStudent)
func (h *StudentHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req models.UpdateStudentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if err := h.studentService.UpdateProfile(r.Context(), userID, &req); err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, repository.ErrStudentNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update profile")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "student profile updated successfully"})
}

func (h *StudentHandler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.studentService.DeleteProfile(r.Context(), userID); err != nil {
		if errors.Is(err, repository.ErrStudentNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete profile")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "student account deleted successfully"})
}
