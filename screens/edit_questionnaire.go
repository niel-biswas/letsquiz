package screens

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"letsquiz/logger"
	"letsquiz/models"
	"letsquiz/views"
	"strconv"
)

type EditQuestionnaire struct {
	model models.EditQuestionnaireModel
}

func InitialEditQuestionnaire() tea.Model {
	logger.Info("InitialEditQuestionnaireModel called")
	model := models.InitialEditQuestionnaireModel()
	return &EditQuestionnaire{model: model}
}

func (m EditQuestionnaire) Init() tea.Cmd {
	logger.Info("EditQuestionnaire Init called")
	return models.FetchQuizzesCmd()
}

func (m EditQuestionnaire) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	logger.Info("EditQuestionnaire Update called", "CurrentScreen", m.model.CurrentScreen, "modelType", fmt.Sprintf("%T", m.model))
	switch m.model.CurrentScreen {
	case "QuizMetadata":
		logger.Info("Transitioning to QuizMetadata screen")

		var quizMetadata *models.QuizMetadata
		if selectedRow := m.model.Table.SelectedRow(); selectedRow != nil && m.model.Focused == "table" {
			logger.Info("Selected row", "row", selectedRow)
			Id, err := strconv.Atoi(selectedRow[0])
			if err != nil {
				logger.Error("Error converting QuizId from selectedRow[0]", "err", err)
			}
			logger.Info("Successful converting QuizId from selectedRow[0]", "Id", Id)
			CategoryId, _ := models.FetchCategoryIDByName(selectedRow[4])
			CreatorId, _ := models.FetchCreatorIDByName(selectedRow[5])
			quizMetadata = &models.QuizMetadata{
				ID:              Id,                           // QuizId
				Title:           selectedRow[1],               // Title
				Description:     selectedRow[2],               // Description
				ContentURL:      selectedRow[3],               // Content URL
				CategoryId:      CategoryId,                   // Category ID
				CreatorId:       CreatorId,                    // Creator ID
				TimeLimitInMins: stringToInt(selectedRow[6]),  // Time Limit in Mins
				IsActive:        stringToBool(selectedRow[7]), // Is Active
				QuestionCount:   stringToInt(selectedRow[8]),  // Question Count
			}
		}

		quizMetadataModel := InitialQuizMetadata(quizMetadata, m.model.Focused, m.model.Buttons)
		return quizMetadataModel, quizMetadataModel.Init()
	default:
		newModel, cmd := models.UpdateEditQuestionnaire(m.model, msg)
		if updatedModel, ok := newModel.(models.EditQuestionnaireModel); ok {
			m.model = updatedModel
			logger.Info("Updated EditQuestionnaireModel", "screen", m.model.CurrentScreen)
		} else {
			logger.Error("Failed to assert model to EditQuestionnaireModel")
			return newModel, cmd // Return the new model and cmd instead of m and nil
		}
		return m, cmd
	}
}

func (m EditQuestionnaire) View() string {
	logger.Info("EditQuestionnaire View called with CurrentScreen", "screen", m.model.CurrentScreen)
	return views.ViewEditQuestionnaire(m.model)
}

func stringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		logger.Error("Error converting string to int", "string", s, "error", err)
		return 0
	}
	return i
}

func stringToBool(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		logger.Error("Error converting string to bool", "string", s, "error", err)
		return false
	}
	return b
}
