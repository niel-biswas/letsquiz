package models

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"io"
	"letsquiz/common"
	"letsquiz/config"
	"letsquiz/logger"
	"letsquiz/music"
	"net/http"
	"strconv"
)

type Quiz struct {
	ID              int
	Title           string
	Description     string
	ContentURL      string
	Category        string
	Creator         string
	TimeLimitInMins int
	IsActive        bool
	QuestionCount   int
}

type EditQuestionnaireModel struct {
	common.Model
	Quizzes []Quiz
}

func InitialEditQuestionnaireModel() EditQuestionnaireModel {
	logger.Info("InitialEditQuestionnaireModel called")
	t := table.New(
		table.WithColumns([]table.Column{
			{Title: "ID", Width: 7},
			{Title: "Title", Width: 20},
			{Title: "Description", Width: 30},
			{Title: "Content URL(Tiny)", Width: 38},
			{Title: "Category", Width: 20},
			{Title: "Creator", Width: 20},
			{Title: "Time Limit(Mins)", Width: 17},
			{Title: "Is Active", Width: 10},
			{Title: "Questions Count", Width: 15},
		}),
		table.WithWidth(195),
	)
	model := EditQuestionnaireModel{
		Model:   common.Model{CurrentScreen: "EditQuestionnaire", Table: t},
		Quizzes: []Quiz{},
	}
	model.Buttons = append(model.Buttons, common.Button{Label: "Create"})
	model.Focused = "button" // Initially focus on the button
	return model
}

func (m *EditQuestionnaireModel) ToggleFocus() {
	if m.Focused == "button" {
		m.Focused = "table"
	} else {
		m.Focused = "button"
	}
}

func FetchQuizzesCmd() tea.Cmd {
	return func() tea.Msg {
		logger.Info("Starting FetchQuizzesCmd")
		// Make the API call to fetch quizzes
		url := fmt.Sprintf("%s/%s", config.AppConfig.BackendURL, "quizzes")
		logger.Info("Fetching quizzes from URL", "url", url)
		resp, err := http.Get(url)
		if err != nil {
			logger.Error("Failed to fetch quizzes", "error", err)
			return nil
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error("Failed to read response body", "error", err)
			return nil
		}

		logger.Info("Fetched response body", "body", string(body))

		// Check if response is an error message
		if string(body) == "invalid transaction\n" {
			logger.Error("Received invalid transaction error from backend")
			return nil
		}

		// Unmarshal the response into a slice of maps
		var quizzes []map[string]interface{}
		if err := json.Unmarshal(body, &quizzes); err != nil {
			logger.Error("Failed to unmarshal quizzes", "error", err)
			logger.Error("Response Body", "body", string(body))
			return nil
		}

		logger.Info("Successfully fetched quizzes", "count", len(quizzes))
		logger.Info("Quizzes Unmarshalled", "Quizzes", quizzes)

		// Resolve quizzes into the Quiz struct
		var resolvedQuizzes []Quiz
		for _, quiz := range quizzes {
			quizId := int(quiz["id"].(float64))
			categoryId := int(quiz["category_id"].(float64))
			creatorId := int(quiz["creator_id"].(float64))

			logger.Info("Processing quiz", "quizId", quizId)

			// Fetch category
			category, err := fetchCategory(categoryId)
			if err != nil {
				logger.Error("Failed to fetch category", "error", err)
				continue
			}
			// Fetch creator
			creator, err := fetchCreator(creatorId)
			if err != nil {
				logger.Error("Failed to fetch creator", "error", err)
				continue
			}
			contentUrl := quiz["content_url"].(string)
			timeLimit := int(quiz["time_limit_in_mins"].(float64))
			isActive := quiz["is_active"].(bool)
			questionCount := int(quiz["question_count"].(float64))

			// Append the resolved quiz to the list
			resolvedQuizzes = append(resolvedQuizzes, Quiz{
				ID:              quizId,
				Title:           quiz["title"].(string),
				Description:     quiz["description"].(string),
				ContentURL:      contentUrl,
				Category:        category,
				Creator:         creator,
				TimeLimitInMins: timeLimit,
				IsActive:        isActive,
				QuestionCount:   questionCount,
			})
		}
		logger.Info("Resolved quizzes", "table values", resolvedQuizzes)
		logger.Info("Resolved quizzes", "resolvedCount", len(resolvedQuizzes))
		return resolvedQuizzes
	}
}

func fetchCategory(categoryId int) (string, error) {
	url := fmt.Sprintf("%s/%s/%d", config.AppConfig.BackendURL, "categories", categoryId)
	logger.Info("Fetching category from URL", "url", url)
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("Failed to fetch category", "error", err)
		return "", err
	}
	defer resp.Body.Close()

	var category map[string]interface{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read response body (category)", "error", err)
		return "", err
	}
	logger.Info("Fetched category response body", "body", string(body))

	if err := json.Unmarshal(body, &category); err != nil {
		logger.Error("Failed to unmarshal category", "error", err)
		logger.Error("Response Body", "body", string(body))
		return "", err
	}

	categoryName, ok := category["name"].(string)
	if !ok {
		logger.Error("Invalid category name format", "category", category)
		return "", fmt.Errorf("invalid category name format")
	}

	return categoryName, nil
}

func fetchCreator(creatorId int) (string, error) {
	url := fmt.Sprintf("%s/%s/%d", config.AppConfig.BackendURL, "users", creatorId)
	logger.Info("Fetching creator from URL", "url", url)
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("Failed to fetch creator", "error", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read response body (creator)", "error", err)
		return "", err
	}
	logger.Info("Fetched creator response body", "body", string(body))

	var user map[string]interface{}
	if err := json.Unmarshal(body, &user); err != nil {
		logger.Error("Failed to unmarshal creator", "error", err)
		logger.Error("Response Body", "body", string(body))
		return "", err
	}

	creatorName, ok := user["user_name"].(string)
	if !ok {
		logger.Error("Invalid creator name format", "user", user)
		return "", fmt.Errorf("invalid creator name format")
	}

	return creatorName, nil
}

func UpdateEditQuestionnaire(m EditQuestionnaireModel, msg tea.Msg) (tea.Model, tea.Cmd) {
	logger.Info("UpdateEditQuestionnaire called", "message", msg, "currentScreen", m.CurrentScreen)
	logger.Info("message type", "msg.(type)", fmt.Sprintf("%T", msg))
	var cmd tea.Cmd

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
		case "tab", "shift+tab":
			m.ToggleFocus()
		case "up", "k":
			if m.Focused == "table" {
				m.Table.MoveUp(1)
			}
		case "down", "j":
			if m.Focused == "table" {
				m.Table.MoveDown(1)
			}
		case "alt+end":
			logger.Info("Toggling music mute/unmute", "currentScreen", m.CurrentScreen)
			music.ToggleMusicMuteUnmute()
		case "enter":
			if m.Focused == "table" {
				m.CurrentScreen = "QuizMetadata" // Set the CurrentScreen here to transition to
				logger.Info("Setting CurrentScreen to QuizMetadata", "screen", m.CurrentScreen)
				selectedRow := m.Table.SelectedRow()
				logger.Info("Selected row", selectedRow)
			} else if m.Focused == "button" && m.Buttons[0].Label == "Create" {
				m.CurrentScreen = "QuizMetadata" // Set the CurrentScreen here to transition to
				logger.Info("Setting CurrentScreen to QuizMetadata", "screen", m.CurrentScreen)
			}
		}
	case []Quiz:
		// Handling the quizzes data received from FetchQuizzesCmd (Call to Database)
		logger.Info("Received quizzes message", "count", len(msg))
		m.Quizzes = msg
		if len(msg) == 0 {
			logger.Info("No quizzes found, setting empty rows")
			m.Table.SetRows([]table.Row{})
		} else {
			logger.Info("Setting table rows with quizzes data")
			rows := convertToTableRows(msg)
			logger.Info("Converted quizzes to table rows", "rows", rows)
			m.Table.SetRows(rows)
			logger.Info("Setting table rows with quizzes data")
		}
	}

	return m, cmd
}

func convertToTableRows(quizzes []Quiz) []table.Row {
	var rows []table.Row
	for _, quiz := range quizzes {
		quizId := strconv.Itoa(quiz.ID)
		quizTimeLimitMins := strconv.Itoa(quiz.TimeLimitInMins)
		questionCount := strconv.Itoa(quiz.QuestionCount)
		quizIsActive := strconv.FormatBool(quiz.IsActive)
		rows = append(rows, table.Row{
			quizId,
			quiz.Title,
			quiz.Description,
			quiz.ContentURL,
			quiz.Category,
			quiz.Creator,
			quizTimeLimitMins,
			quizIsActive,
			questionCount,
		})
	}
	return rows
}
