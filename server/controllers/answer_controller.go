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

// GetAnswers handles GET requests to fetch all answers
func GetAnswers(w http.ResponseWriter, r *http.Request) {
	var answers []models.Answer
	if err := database.DB.Find(&answers).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(answers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetAnswerByID handles GET requests to fetch a single answer by ID
func GetAnswerByID(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/answers/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid answer ID", http.StatusBadRequest)
		return
	}

	var answer models.Answer
	if err := database.DB.First(&answer, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Answer not found", http.StatusNotFound)
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

// GetAnswersByQuestionID handles GET requests to fetch all answers for a specific question ID
func GetAnswersByQuestionID(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/questions/")
	idParam = strings.TrimSuffix(idParam, "/answers")
	questionID, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	var answers []models.Answer
	if err := database.DB.Where("question_id = ?", questionID).Find(&answers).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(answers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CreateAnswer handles POST requests to create a new answer
func CreateAnswer(w http.ResponseWriter, r *http.Request) {
	logger.Info("CreateAnswer called")

	// Log request details
	logger.Info("CreateAnswer", "Request Method:", r.Method, "Request URL:", r.URL.String())

	var answer models.Answer
	if err := json.NewDecoder(r.Body).Decode(&answer); err != nil {
		logger.Error("CreateAnswer", "Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Log decoded answer data
	logger.Info("CreateAnswer", "Decoded answer data:", answer)

	if err := database.DB.Create(&answer).Error; err != nil {
		logger.Error("CreateAnswer", "Error creating answer in database:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log created answer
	logger.Info("CreateAnswer", "Answer created successfully:", answer)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(answer); err != nil {
		logger.Error("CreateAnswer", "Error encoding response:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		logger.Info("CreateAnswer", "Response sent successfully")
	}
}

// UpdateAnswer handles PUT requests to update an existing answer
func UpdateAnswer(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/answers/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid answer ID", http.StatusBadRequest)
		return
	}

	var answer models.Answer
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
