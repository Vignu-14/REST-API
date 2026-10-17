package handlers

import (
	"encoding/json"
	"net/http"
	"rest-api/internal/models"
	"rest-api/internal/repository"
)

func StudentSignupHandler(r *repository.RepoUser) http.HandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.RegisterStudentRequest
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "invalid data form", http.StatusBadRequest)
			return
		}


	}
}


