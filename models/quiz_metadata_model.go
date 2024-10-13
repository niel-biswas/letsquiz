package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"io"
	"letsquiz/common"
	"letsquiz/config"
	"letsquiz/logger"
	"letsquiz/music"
	"letsquiz/server/models"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type QuizMetadata struct {
	ID               int       `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	ContentURL       string    `json:"content_url"`
	CategoryId       int       `json:"category_id"`
	CreatorId        int       `json:"creator_id"`
	TimeLimitInMins  int       `json:"time_limit_in_mins"`
	QuestionCount    int       `json:"question_count"`
	IsActive         bool      `json:"is_active"`
	CreationDate     time.Time `json:"creation_date"`
	LastModifiedDate time.Time `json:"last_modified_date"`
}

var timeLimitInMins, isActive, questionCount string

type QuizMetadataModel struct {
	common.Model
	QuizData QuizMetadata
}

// InitialQuizMetadata initializes the QuizMetadataModel with a form
func InitialQuizMetadata(quizData *QuizMetadata, focused string, buttons []common.Button) QuizMetadataModel {
	fields := QuizMetadata{}
	category := Category{}
	if quizData != nil {
		logger.Info("InitialQuizMetadata called", "quizId", quizData.ID, "focused", focused, "buttons", buttons)
		fields = *quizData
		timeLimitInMins = strconv.Itoa(quizData.TimeLimitInMins)
		questionCount = strconv.Itoa(quizData.QuestionCount)
		category.Name, category.Description, _ = FetchCategoryNameDescByID(quizData.CategoryId)
	}
	logger.Info("InitialQuizMetadata called", "focused", focused, "buttons", buttons)
	categorySuggestions, err := FetchCategories()
	if err != nil {
		logger.Info("Error while fetching categorySuggestions", "error", err)
	} else {
		logger.Info("Fetched category suggestions", "categories", categorySuggestions)
	}

	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Title").Key("title").Description("Stand out giving a kool name to your quiz!").
				Validate(func(val string) error {
					if val == "" {
						return errors.New("This field is required")
					}
					return nil
				}).
				Value(&fields.Title),
			huh.NewInput().Key("description").
				Title("Description").Description("Let others know why your quiz is so kool!").
				Value(&fields.Description),
			huh.NewInput().Key("content_url").
				Title("Content URL").Description("Shine like a star by guiding others to your knowledge repository!").
				Value(&fields.ContentURL),
			huh.NewInput().Key("category").
				Title("Category").Description("Enter a category of your quiz.").
				Validate(func(val string) error {
					if val == "" {
						return errors.New("This field is required")
					}
					return nil
				}).
				Suggestions(categorySuggestions).
				Value(&category.Name),
			huh.NewInput().Key("time_limit_in_mins").
				Title("Time Limit (mins)").Description("Time your quiz wisely so that others get a chance to play your quiz.").
				Validate(func(val string) error {
					if val == "" {
						return errors.New("This field is required")
					}
					return nil
				}).
				Value(&timeLimitInMins),
			huh.NewInput().Key("question_count").
				Title("Question Count").Description("How many questions do you plan?").
				Validate(func(val string) error {
					if val == "" || val == "0" {
						return errors.New("You should have atleast a question in your Quiz.")
					}
					return nil
				}).
				Value(&questionCount),
			huh.NewSelect[string]().
				Key("is_active").
				Options(huh.NewOptions("Yes", "No")...).
				Title("Is Active?").
				Description("Use this in case you need to mark your quiz Inactive or Obsolete.").
				Validate(func(val string) error {
					if val == "" {
						return errors.New("This field is required")
					}
					return nil
				}).
				Value(&isActive),
			huh.NewConfirm().
				Key("done").
				Title("All done?").Description("Seriously...").
				Validate(func(v bool) error {
					if !v {
						return fmt.Errorf("Welp, finish up then")
					}
					return nil
				}).
				Affirmative("Yep").
				Negative("Wait, no"),
		),
	).WithShowHelp(true).WithShowErrors(true)
	model := QuizMetadataModel{
		Model:    common.Model{CurrentScreen: "QuizMetadata", Form: f, Focused: focused, Buttons: buttons},
		QuizData: fields,
	}
	return model
}

func UpdateQuizMetadata(m QuizMetadataModel, msg tea.Msg) (tea.Model, tea.Cmd) {
	logger.Info("UpdateQuizMetadata called", "message", msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.WindowWidth = msg.Width
		m.WindowHeight = msg.Height
		logger.Info("Window size updated", "width", m.WindowWidth, "height", m.WindowHeight)
	case tea.KeyMsg:
		logger.Info("Key pressed", "key", msg.String(), "currentScreen", m.CurrentScreen)
		switch msg.String() {
		case "esc":
			logger.Info("quitting application", "currentScreen", m.CurrentScreen)
			return m, tea.Quit
		case "alt+end":
			logger.Info("Toggling music mute/unmute", "currentScreen", m.CurrentScreen)
			music.ToggleMusicMuteUnmute()
		}
	}

	var cmds []tea.Cmd

	// Process the form
	form, cmd := m.Form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.Form = f
		cmds = append(cmds, cmd)
	}

	if m.Form.State == huh.StateCompleted {
		categoryId, err := FetchCategoryIDByName(m.Form.GetString("category"))
		if err != nil {
			logger.Error("Error fetching category ID by name", "error", err)
			return m, nil
		}

		timeLimitInMins, err := strconv.Atoi(strings.Trim(timeLimitInMins, " "))
		if err != nil {
			logger.Error("Error fetching timeLimitInMins", "error", err)
			return m, nil
		}

		isActive, err := parseBool(isActive)
		if err != nil {
			logger.Error("Error parsing isActive", "error", err)
			return m, nil
		}

		questionCount, err := strconv.Atoi(strings.Trim(questionCount, " "))
		if err != nil {
			logger.Error("Error fetching questionCount", "error", err)
			return m, nil
		} else {
			m.QuizData.QuestionCount = questionCount
		}

		metadata := QuizMetadata{
			ID:               m.QuizData.ID,
			Title:            m.Form.GetString("title"),
			Description:      m.Form.GetString("description"),
			ContentURL:       m.Form.GetString("content_url"),
			CategoryId:       categoryId,
			CreatorId:        m.QuizData.CreatorId,
			TimeLimitInMins:  timeLimitInMins,
			QuestionCount:    questionCount,
			IsActive:         isActive,
			CreationDate:     time.Now(),
			LastModifiedDate: time.Now(),
		}
		logger.Info("Quiz metadata prepared", "metadata", metadata, "m.Focused", m.Focused, "m.Buttons", m.Buttons)
		if m.Focused == "table" {
			if err := saveQuizMetadata(metadata); err != nil {
				logger.Error("Error posting quiz metadata (saveQuizMetadata)", "error", err)
			} else {
				logger.Info("Successfully posted quiz metadata (saveQuizMetadata)", "metadata", metadata)
			}
		}
		if m.Focused == "button" && m.Buttons[0].Label == "Create" {
			if err := postQuizMetadata(metadata); err != nil {
				logger.Error("Error posting quiz metadata (postQuizMetadata)", "error", err)
			} else {
				logger.Info("Successfully posted quiz metadata (postQuizMetadata)", "metadata", metadata)
			}
		}
		// Quit when the form is done.
		// cmds = append(cmds, tea.Quit)
		// Signal the form completion by returning a command.
		return m, func() tea.Msg { return common.QuizMetaDataFormCompletedMsg{} }
	}

	return m, tea.Batch(cmds...)
}

func postQuizMetadata(metadata QuizMetadata) error {
	url := config.AppConfig.BackendURL + "/quizzes" // Read from AppConfig
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to post quiz metadata, status code: %d", resp.StatusCode)
	}
	return nil
}

// saveQuizMetadata saves the quiz metadata using a PUT request
func saveQuizMetadata(metadata QuizMetadata) error {
	url := fmt.Sprintf("%s/quizzes/%d", config.AppConfig.BackendURL, metadata.ID)
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to put quiz metadata, status code: %d", resp.StatusCode)
	}
	return nil
}

// FetchCategories fetches all quiz categories from the backend and returns them as a slice of strings
func FetchCategories() ([]string, error) {
	url := config.AppConfig.BackendURL + "/categories" // Adjust the URL if needed

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching categories: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch categories, status code: %d", resp.StatusCode)
	}

	var categories []models.Category
	if err := json.NewDecoder(resp.Body).Decode(&categories); err != nil {
		return nil, fmt.Errorf("error decoding categories response: %w", err)
	}

	var categoryNames []string
	for _, category := range categories {
		logger.Info("Fetched category", "id", category.ID, "name", category.Name)
		categoryNames = append(categoryNames, category.Name)
	}

	return categoryNames, nil
}

// FetchCategoryNameDescByID fetches the name and description of a category by its ID from the backend
func FetchCategoryNameDescByID(categoryID int) (string, string, error) {
	url := fmt.Sprintf("%s/categories/%d", config.AppConfig.BackendURL, categoryID)

	resp, err := http.Get(url)
	if err != nil {
		return "", "", fmt.Errorf("error fetching category name: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Info("Failed to fetch category name", "status", resp.StatusCode, "body", string(body), "categoryID", categoryID, "error", err)
		return "", "", fmt.Errorf("failed to fetch category name, status code: %d", resp.StatusCode)
	}

	var category Category
	if err := json.NewDecoder(resp.Body).Decode(&category); err != nil {
		return "", "", fmt.Errorf("error decoding category name response: %w", err)
	}

	logger.Info("Fetched category name", "id", category.ID, "name", category.Name)
	return category.Name, category.Description, nil
}

// FetchCategoryIDByName fetches the ID of a category by its name from the backend
func FetchCategoryIDByName(categoryName string) (int, error) {
	// Escape the category name to ensure it's a valid URL component
	escapedCategoryName := url.QueryEscape(categoryName)
	// Replace '+' with '%20' for spaces
	escapedCategoryName = strings.ReplaceAll(escapedCategoryName, "+", "%20")
	url := fmt.Sprintf("%s/categories/byname/%s", config.AppConfig.BackendURL, escapedCategoryName)

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error fetching category ID: %w", err)
	}
	defer resp.Body.Close()

	if err != nil {
		body, _ := io.ReadAll(resp.Body)
		logger.Info("Failed to fetch category ID", "status", resp.StatusCode, "body", string(body), "categoryName", categoryName, "error", err)
		return 0, fmt.Errorf("error fetching category ID: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to fetch category ID, status code: %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("error decoding category ID response: %w", err)
	}

	id, ok := response["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid response format")
	}

	logger.Info("Fetched category ID", "name", categoryName, "id", int(id))
	return int(id), nil
}

// FetchCreatorIDByName fetches the ID of a creator by their name from the backend
func FetchCreatorIDByName(creatorName string) (int, error) {
	// Escape the creator name to ensure it's a valid URL component
	escapedCreatorName := url.QueryEscape(creatorName)
	// Replace '+' with '%20' for spaces
	escapedCreatorName = strings.ReplaceAll(escapedCreatorName, "+", "%20")
	url := fmt.Sprintf("%s/users/byname/%s", config.AppConfig.BackendURL, escapedCreatorName)

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error fetching creator ID: %w", err)
	}
	defer resp.Body.Close()

	if err != nil {
		body, _ := io.ReadAll(resp.Body)
		logger.Info("Failed to fetch creator ID", "status", resp.StatusCode, "body", string(body), "creatorName", creatorName, "error", err)
		return 0, fmt.Errorf("error fetching creator ID: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to fetch creator ID, status code: %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("error decoding creator ID response: %w", err)
	}

	id, ok := response["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid response format")
	}

	logger.Info("Fetched creator ID", "name", creatorName, "id", int(id))
	return int(id), nil
}

// parseBool converts "Yes"/"No" to true/false
func parseBool(input string) (bool, error) {
	switch input {
	case "Yes":
		return true, nil
	case "No":
		return false, nil
	default:
		return false, errors.New("invalid input: expected 'Yes' or 'No'")
	}
}
