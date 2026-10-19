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

type TeacherHandler struct {
	teacherService services.TeacherService
}

func NewTeacherHandler(teacherService services.TeacherService) *TeacherHandler {
	return &TeacherHandler{
		teacherService: teacherService,
	}
}

// 1. Register handles POST /teachers/register (Public route)
func (h *TeacherHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.TeacherRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	profile, err := h.teacherService.Register(r.Context(), &req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusConflict, "failed to register teacher (username or email may already exist)")
		return
	}

	writeJSON(w, http.StatusCreated, profile)
}

// 2. GetProfile handles GET /teachers/me (Protected: RoleTeacher)
func (h *TeacherHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	profile, err := h.teacherService.GetProfile(r.Context(), userID)
	if err != nil {
		if errors.Is(err, repository.ErrTeacherNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to fetch profile")
		return
	}

	writeJSON(w, http.StatusOK, profile)
}

// 3. UpdateProfile handles PUT /teachers/me (Protected: RoleTeacher)
func (h *TeacherHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req models.UpdateTeacherRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if err := h.teacherService.UpdateProfile(r.Context(), userID, &req); err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, repository.ErrTeacherNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update profile")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "teacher profile updated successfully"})
}

// 4. DeleteProfile handles DELETE /teachers/me (Protected: RoleTeacher)
func (h *TeacherHandler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.teacherService.DeleteProfile(r.Context(), userID); err != nil {
		if errors.Is(err, repository.ErrTeacherNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete profile")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "teacher account deleted successfully"})
}

// 5. GetMyClassStudents handles GET /teachers/my-students (Protected: RoleTeacher)
func (h *TeacherHandler) GetMyClassStudents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	students, err := h.teacherService.GetMyClassStudents(r.Context(), userID)
	if err != nil {
		if errors.Is(err, repository.ErrTeacherNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to fetch class students")
		return
	}

	writeJSON(w, http.StatusOK, students)
}
