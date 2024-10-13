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

// GetLeaderboards handles GET requests to fetch all leaderboards
func GetLeaderboards(w http.ResponseWriter, r *http.Request) {
	var leaderboards []models.Leaderboard
	if err := database.DB.Find(&leaderboards).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(leaderboards); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetLeaderboardByID handles GET requests to fetch a single leaderboard by ID
func GetLeaderboardByID(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/leaderboards/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid leaderboard ID", http.StatusBadRequest)
		return
	}

	var leaderboard models.Leaderboard
	if err := database.DB.First(&leaderboard, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Leaderboard not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(leaderboard); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CreateLeaderboard handles POST requests to create a new leaderboard
func CreateLeaderboard(w http.ResponseWriter, r *http.Request) {
	var leaderboard models.Leaderboard
	if err := json.NewDecoder(r.Body).Decode(&leaderboard); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.DB.Create(&leaderboard).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	database.DB.Commit()

	w.WriteHeader(http.StatusCreated)
}

// UpdateLeaderboard handles PUT requests to update an existing leaderboard
func UpdateLeaderboard(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/leaderboards/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid leaderboard ID", http.StatusBadRequest)
		return
	}

	var leaderboard models.Leaderboard
	if err := json.NewDecoder(r.Body).Decode(&leaderboard); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	leaderboard.ID = id
	if err := database.DB.Save(&leaderboard).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	database.DB.Commit()

	w.WriteHeader(http.StatusOK)
}
