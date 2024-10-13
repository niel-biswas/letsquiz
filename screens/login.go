package screens

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"letsquiz/common"
	"letsquiz/logger"
	"letsquiz/models"
	"letsquiz/views"
)

type Login struct {
	model models.LoginModel
}

func InitialModel() tea.Model {
	logger.Info("InitialModel called")
	model := models.InitialLoginModel()
	return &Login{
		model: model,
	}
}

func (m *Login) Init() tea.Cmd {
	logger.Info("Login Init called")
	return m.model.Init()
}

func (m *Login) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	logger.Info("Login Update called", "CurrentScreen", m.model.CurrentScreen, "modelType", fmt.Sprintf("%T", m.model))
	if m.model.CurrentScreen == "login" {
		newModel, cmd := models.UpdateLogin(m.model.Model, msg)
		m.model.Model = newModel.(common.Model)
		return m, cmd
	}
	if m.model.CurrentScreen == "menu" {
		menuModel := InitialMenu()
		return menuModel, menuModel.Init()
	}
	return m, nil
}

func (m *Login) View() string {
	logger.Info("Login View called with CurrentScreen", "screen", m.model.CurrentScreen)
	if m.model.CurrentScreen == "login" {
		return views.ViewLogin(m.model.Model)
	} else if m.model.CurrentScreen == "menu" {
		return views.ViewMenu(m.model.Model)
	}
	return ""
}
