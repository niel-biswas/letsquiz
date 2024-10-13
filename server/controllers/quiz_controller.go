package controllers

import (
	"encoding/json"
	"gorm.io/gorm"
	"letsquiz/logger"
	"letsquiz/server/database"
	"letsquiz/server/models"
	"net/http"
	"strconv"
	"strings"
)

// GetQuizzes handles GET requests to fetch all quizzes
func GetQuizzes(w http.ResponseWriter, r *http.Request) {
	logger.Info("GetQuizzes called")

	// Log request details
	logger.Info("GetQuizzes", "Request Method:", r.Method, "Request URL:", r.URL.String())

	var quizzes []models.Quiz
	if err := database.DB.Find(&quizzes).Error; err != nil {
		logger.Error("GetQuizzes", "Error fetching quizzes from database:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log fetched quizzes
	logger.Info("GetQuizzes", "Fetched quizzes from database:", quizzes)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(quizzes); err != nil {
		logger.Error("GetQuizzes", "Error encoding response:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("GetQuizzes", "Response sent successfully")
}

// GetQuizByID handles GET requests to fetch a single quiz by ID
func GetQuizByID(w http.ResponseWriter, r *http.Request) {
	logger.Info("GetQuizByID called")

	// Log request details
	logger.Info("GetQuizByID", "Request Method:", r.Method, "Request URL:", r.URL.String())

	idParam := strings.TrimPrefix(r.URL.Path, "/quizzes/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		logger.Error("GetQuizByID", "Invalid quiz ID:", idParam)
		http.Error(w, "Invalid quiz ID", http.StatusBadRequest)
		return
	}

	logger.Info("GetQuizByID", "Fetching quiz with ID:", id)

	var quiz models.Quiz
	if err := database.DB.First(&quiz, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Error("GetQuizByID", "Quiz not found:", id)
			http.Error(w, "Quiz not found", http.StatusNotFound)
		} else {
			logger.Error("GetQuizByID", "Error fetching quiz from database:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Log fetched quiz
	logger.Info("GetQuizByID", "Fetched quiz from database:", quiz)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(quiz); err != nil {
		logger.Error("GetQuizByID", "Error encoding response:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("GetQuizByID", "Response sent successfully")
}

// CreateQuiz handles POST requests to create a new quiz
func CreateQuiz(w http.ResponseWriter, r *http.Request) {
	logger.Info("CreateQuiz called")

	// Log request details
	logger.Info("CreateQuiz", "Request Method:", r.Method, "Request URL:", r.URL.String())

	var quiz models.Quiz
	if err := json.NewDecoder(r.Body).Decode(&quiz); err != nil {
		logger.Error("CreateQuiz", "Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Log decoded quiz data
	logger.Info("CreateQuiz", "Decoded quiz data:", quiz)

	if err := database.DB.Create(&quiz).Error; err != nil {
		logger.Error("CreateQuiz", "Error creating quiz in database:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Capture the created quiz's ID and return it in the response
	response := map[string]interface{}{
		"id": quiz.ID,
	}

	// Log created quiz
	logger.Info("CreateQuiz", "Quiz created successfully:", quiz)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("CreateQuiz", "Error encoding response:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// UpdateQuiz handles PUT requests to update an existing quiz
func UpdateQuiz(w http.ResponseWriter, r *http.Request) {
	logger.Info("UpdateQuiz called")

	// Log request details
	logger.Info("UpdateQuiz", "Request Method:", r.Method, "Request URL:", r.URL.String())
	idParam := strings.TrimPrefix(r.URL.Path, "/quizzes/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		logger.Error("UpdateQuiz", "Invalid quiz ID:", idParam)
		http.Error(w, "Invalid quiz ID", http.StatusBadRequest)
		return
	}

	logger.Info("UpdateQuiz", "Updating quiz with ID:", id)

	var quiz models.Quiz
	if err := json.NewDecoder(r.Body).Decode(&quiz); err != nil {
		logger.Error("UpdateQuiz", "Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	quiz.ID = id
	logger.Info("UpdateQuiz", "Quiz data to update:", quiz)

	if err := database.DB.Save(&quiz).Error; err != nil {
		logger.Error("UpdateQuiz", "Error updating quiz in database:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log updated quiz
	logger.Info("UpdateQuiz", "Quiz updated successfully:", quiz)
	w.WriteHeader(http.StatusOK)
}
