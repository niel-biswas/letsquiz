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

// GetQuestions handles GET requests to fetch all questions
func GetQuestions(w http.ResponseWriter, r *http.Request) {
	var questions []models.Question
	if err := database.DB.Find(&questions).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(questions); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetQuestionByID handles GET requests to fetch a single question by ID
func GetQuestionByID(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/questions/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	var question models.Question
	if err := database.DB.First(&question, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Question not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(question); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetQuestionsByQuizID handles GET requests to fetch all questions for a specific quiz ID
func GetQuestionsByQuizID(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/quizzes/")
	idParam = strings.TrimSuffix(idParam, "/questions")
	quizID, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid quiz ID", http.StatusBadRequest)
		return
	}

	var questions []models.Question
	if err := database.DB.Where("quiz_id = ?", quizID).Find(&questions).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(questions); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CreateQuestion handles POST requests to create a new question
func CreateQuestion(w http.ResponseWriter, r *http.Request) {
	var question models.Question
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.DB.Create(&question).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Capture the created question's ID and return it in the response
	response := map[string]interface{}{
		"id": question.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// UpdateQuestion handles PUT requests to update an existing question
func UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/questions/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	var question models.Question
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	question.ID = id
	if err := database.DB.Save(&question).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	database.DB.Commit()

	w.WriteHeader(http.StatusOK)
}
