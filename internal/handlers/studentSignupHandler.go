package handlers

import (
	"encoding/json"
	"net/http"
	"rest-api/internal/models"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var user models.RegisterStudentRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "invalid data form", http.StatusBadRequest)
		return
	}

	// Checking if password is greater than characters.
	if len(user.Password) > 20 {
		http.Error(w, "password must be at most 20 characters", http.StatusBadRequest)
		return
	}

	


}
