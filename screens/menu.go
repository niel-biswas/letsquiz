package screens

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"letsquiz/common"
	"letsquiz/logger"
	"letsquiz/models"
	"letsquiz/views"
)

type Menu struct {
	model models.MenuModel
}

func InitialMenu() tea.Model {
	logger.Info("InitialMenuModel called")
	model := models.InitialMenuModel()
	return &Menu{
		model: model,
	}
}

func (m *Menu) Init() tea.Cmd {
	logger.Info("Menu Init called")
	return m.model.Init()
}

func (m *Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	logger.Info("Menu Update called", "CurrentScreen", m.model.CurrentScreen, "modelType", fmt.Sprintf("%T", m.model))
	switch m.model.CurrentScreen {
	case "menu":
		newModel, cmd := models.UpdateMenu(m.model.Model, msg)
		m.model.Model = newModel.(common.Model)
		return m, cmd

	case "EditQuestionnaire":
		editQuestionnaireModel := InitialEditQuestionnaire()
		return editQuestionnaireModel, editQuestionnaireModel.Init()

	default:
		return m, nil
	}
}

func (m *Menu) View() string {
	logger.Info("Menu View called with CurrentScreen", "screen", m.model.CurrentScreen)
	return views.ViewMenu(m.model.Model)
}
