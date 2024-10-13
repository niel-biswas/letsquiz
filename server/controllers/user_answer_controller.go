package controllers

import (
	"encoding/json"
	"letsquiz/server/database"
	"letsquiz/server/models"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// GetUserAnswers handles GET requests to fetch all user answers
func GetUserAnswers(w http.ResponseWriter, r *http.Request) {
	var answers []models.UserAnswer
	if err := database.DB.Find(&answers).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(answers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetUserAnswerByID handles GET requests to fetch a single user answer by ID
func GetUserAnswerByID(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/user-answers/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid user answer ID", http.StatusBadRequest)
		return
	}

	var answer models.UserAnswer
	if err := database.DB.First(&answer, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "User answer not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(answer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CreateUserAnswer handles POST requests to create a new user answer
func CreateUserAnswer(w http.ResponseWriter, r *http.Request) {
	var answer models.UserAnswer
	if err := json.NewDecoder(r.Body).Decode(&answer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.DB.Create(&answer).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	database.DB.Commit()

	w.WriteHeader(http.StatusCreated)
}

// UpdateUserAnswer handles PUT requests to update an existing user answer
func UpdateUserAnswer(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/user-answers/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid user answer ID", http.StatusBadRequest)
		return
	}

	var answer models.UserAnswer
	if err := json.NewDecoder(r.Body).Decode(&answer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	answer.ID = id
	if err := database.DB.Save(&answer).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	database.DB.Commit()

	w.WriteHeader(http.StatusOK)
}
