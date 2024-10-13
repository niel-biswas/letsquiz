package controllers

import (
	"encoding/json"
	"letsquiz/logger"
	"letsquiz/server/database"
	"letsquiz/server/models"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// GetFeedbacks handles GET requests to fetch all feedbacks
func GetFeedbacks(w http.ResponseWriter, r *http.Request) {
	var feedbacks []models.Feedback
	if err := database.DB.Find(&feedbacks).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(feedbacks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetFeedbackByID handles GET requests to fetch a single feedback by ID
func GetFeedbackByID(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/feedbacks/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid feedback ID", http.StatusBadRequest)
		return
	}

	var feedback models.Feedback
	if err := database.DB.First(&feedback, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Feedback not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(feedback); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CreateFeedback handles POST requests to create a new feedback
func CreateFeedback(w http.ResponseWriter, r *http.Request) {
	logger.Info("CreateFeedback called")

	// Log request details
	logger.Info("CreateFeedback", "Request Method:", r.Method, "Request URL:", r.URL.String())

	var feedback models.Feedback
	if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
		logger.Error("CreateFeedback", "Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Log decoded feedback data
	logger.Info("CreateFeedback", "Decoded feedback data:", feedback)

	if err := database.DB.Create(&feedback).Error; err != nil {
		logger.Error("CreateFeedback", "Error creating feedback in database:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log created feedback
	logger.Info("CreateFeedback", "Feedback created successfully:", feedback)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(feedback); err != nil {
		logger.Error("CreateFeedback", "Error encoding response:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		logger.Info("CreateFeedback", "Response sent successfully")
	}
}

// UpdateFeedback handles PUT requests to update an existing feedback
func UpdateFeedback(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/feedbacks/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid feedback ID", http.StatusBadRequest)
		return
	}

	var feedback models.Feedback
	if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	feedback.ID = id
	if err := database.DB.Save(&feedback).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	database.DB.Commit()

	w.WriteHeader(http.StatusOK)
}
