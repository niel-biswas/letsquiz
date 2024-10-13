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

// GetUserQuizAttempts handles GET requests to fetch all user quiz attempts
func GetUserQuizAttempts(w http.ResponseWriter, r *http.Request) {
	var attempts []models.UserQuizAttempt
	if err := database.DB.Find(&attempts).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(attempts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetUserQuizAttemptByID handles GET requests to fetch a single user quiz attempt by ID
func GetUserQuizAttemptByID(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/attempts/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid attempt ID", http.StatusBadRequest)
		return
	}

	var attempt models.UserQuizAttempt
	if err := database.DB.First(&attempt, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Attempt not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(attempt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CreateUserQuizAttempt handles POST requests to create a new user quiz attempt
func CreateUserQuizAttempt(w http.ResponseWriter, r *http.Request) {
	var attempt models.UserQuizAttempt
	if err := json.NewDecoder(r.Body).Decode(&attempt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.DB.Create(&attempt).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	database.DB.Commit()

	w.WriteHeader(http.StatusCreated)
}

// UpdateUserQuizAttempt handles PUT requests to update an existing user quiz attempt
func UpdateUserQuizAttempt(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/attempts/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid attempt ID", http.StatusBadRequest)
		return
	}

	var attempt models.UserQuizAttempt
	if err := json.NewDecoder(r.Body).Decode(&attempt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	attempt.ID = id
	if err := database.DB.Save(&attempt).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	database.DB.Commit()

	w.WriteHeader(http.StatusOK)
}
