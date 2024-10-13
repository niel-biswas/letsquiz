package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"letsquiz/logger"
	"letsquiz/server/database"
	"letsquiz/server/models"

	"gorm.io/gorm"
)

// GetCategories handles GET requests to fetch all categories
func GetCategories(w http.ResponseWriter, r *http.Request) {
	var categories []models.Category
	if err := database.DB.Find(&categories).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(categories); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetCategoryByID handles GET requests to fetch a single category by ID
func GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var category models.Category
	if err := database.DB.First(&category, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Category not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(category); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetCategoryIDByName handles GET requests to fetch a category ID by its name
func GetCategoryIDByName(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/categories/byname/")
	if name == "" {
		http.Error(w, "Category name is required", http.StatusBadRequest)
		return
	}

	var category models.Category
	if err := database.DB.Where("name = ?", name).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Category not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"id": category.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CreateCategory handles POST requests to create a new category
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	logger.Info("CreateCategory called")

	// Log request details
	logger.Info("CreateCategory", "Request Method:", r.Method, "Request URL:", r.URL.String())

	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		logger.Error("CreateCategory", "Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Log decoded category data
	logger.Info("CreateCategory", "Decoded category data:", category)

	if err := database.DB.Create(&category).Error; err != nil {
		logger.Error("CreateCategory", "Error creating category in database:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log created category
	logger.Info("CreateCategory", "Category created successfully:", category)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(category); err != nil {
		logger.Error("CreateCategory", "Error encoding response:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		logger.Info("CreateCategory", "Response sent successfully")
	}
}

// UpdateCategory handles PUT requests to update an existing category
func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	category.ID = id
	if err := database.DB.Save(&category).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	database.DB.Commit()

	w.WriteHeader(http.StatusOK)
}
