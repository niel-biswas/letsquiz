package screens

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"letsquiz/common"
	"letsquiz/logger"
	"letsquiz/models"
	"letsquiz/views"
)

type DynamicQuizForms struct {
	model models.DynamicQuizModel
}

func InitialDynamicQuizForms(quizID, questionCount int, focused string, buttons []common.Button) DynamicQuizForms {
	logger.Info("InitialDynamicQuizForms called", "quizID", quizID, "questionsCount", questionCount)
	model := models.InitialDynamicQuizModel(quizID, questionCount, focused, buttons)
	logger.Info("Initialized DynamicQuizForms", "quizID", quizID, "questionsCount", questionCount)
	return DynamicQuizForms{model: model}
}

func (m DynamicQuizForms) Init() tea.Cmd {
	logger.Info("DynamicQuizForms Init called", "CurrentFormGroup", m.model.CurrentFormGroup)
	return m.model.QuestionForms[m.model.CurrentFormGroup].Init()
}

func (m DynamicQuizForms) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	logger.Info("DynamicQuizForms Update called", "msgType", fmt.Sprintf("%T", msg))
	newModel, cmd := models.UpdateDynamicQuizModel(m.model, msg)
	if updatedModel, ok := newModel.(models.DynamicQuizModel); ok {
		m.model = updatedModel
		logger.Info("Updated DynamicQuizModel", "CurrentFormGroup", m.model.CurrentFormGroup, "screen", m.model.CurrentScreen)
	} else {
		logger.Error("Failed to assert model to DynamicQuizModel")
		return newModel, cmd // Return the new model and cmd instead of m and nil
	}
	return m, cmd
}

func (m DynamicQuizForms) View() string {
	logger.Info("DynamicQuizForms View called with CurrentForm", "CurrentFormGroup", m.model.CurrentFormGroup)
	return views.ViewDynamicQuizForm(m.model)
}
