package screens

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"letsquiz/common"
	"letsquiz/logger"
	"letsquiz/models"
	"letsquiz/views"
)

type QuizMetadata struct {
	model models.QuizMetadataModel
}

func InitialQuizMetadata(quizData *models.QuizMetadata, focused string, buttons []common.Button) QuizMetadata {
	var model models.QuizMetadataModel
	if focused == "table" {
		logger.Info("InitialQuizMetadata called", "quizData", quizData.ID, "focused", focused, "buttons", buttons)
		model = models.InitialQuizMetadata(quizData, focused, buttons)
	} else {
		logger.Info("InitialQuizMetadata called", "focused", focused, "buttons", buttons)
		model = models.InitialQuizMetadata(nil, focused, buttons)
	}
	return QuizMetadata{model: model}
}

func (m QuizMetadata) Init() tea.Cmd {
	logger.Info("QuizMetadata Init called")
	return m.model.Form.Init()
}

func (m QuizMetadata) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	logger.Info("QuizMetadata Update called", "CurrentScreen", m.model.CurrentScreen, "modelType", fmt.Sprintf("%T", m.model))
	switch msg.(type) {
	case common.QuizMetaDataFormCompletedMsg:
		quizID := m.model.QuizData.ID
		questionCount := m.model.QuizData.QuestionCount
		dynamicQuizScreen := InitialDynamicQuizForms(quizID, questionCount, m.model.Focused, m.model.Buttons)
		return dynamicQuizScreen, dynamicQuizScreen.Init()
	}

	newModel, cmd := models.UpdateQuizMetadata(m.model, msg)
	if updatedModel, ok := newModel.(models.QuizMetadataModel); ok {
		m.model = updatedModel
		logger.Info("Updated QuizMetadataModel", "screen", m.model.CurrentScreen)
	} else {
		logger.Error("Failed to assert model to QuizMetadataModel")
		return newModel, cmd // Return the new model and cmd instead of m and nil
	}
	return m, cmd
}

func (m QuizMetadata) View() string {
	logger.Info("QuizMetadata View called with CurrentScreen", "screen", m.model.CurrentScreen)
	return views.ViewQuizMetadata(m.model)
}
