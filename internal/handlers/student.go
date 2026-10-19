package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"rest-api/internal/models"
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
func (h *StudentHandler) Register (w http.ResponseWriter,r  *http.Request) {
	if r.Method != http.MethodPost{
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.StudentRegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalied JSON body")
		return
	}
	profile, err := h.studentService.Register(r.Context(), &req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput){
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusConflict, "failed to register student (username, email, or enrollment_no may already exist)")
		return
	}
	writeJSON(w, http.StatusCreated, profile)
}

func (h * 	)