package handlers

import (
	"encoding/json"
	"net/http"
	"rest-api/internal/models"
	"rest-api/internal/repository"
)

func StudentSignupHandler(repo *repository.RepoUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.RegisterStudentRequest
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "invalid data form", http.StatusBadRequest)
			return
		}

		if len(user.Password) > 20 {
			http.Error(w, "password should lesser than 21 characters", http.StatusBadRequest)
			return
		}

		if user.Username == "" || user.Email == "" || user.FullName == "" || user.Password == "" || user.Class == "" || user.EnrollmentNo == ""{
			http.Error(w, "data is invalid", http.StatusBadRequest)
			return
		}

		id, err := repo.DbCreateStudent(r.Context(), user.Username, user.FullName, user.Password, user.Class, user.EnrollmentNo)

	}
}


