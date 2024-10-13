package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"letsquiz/common"
	"letsquiz/logger"
	"letsquiz/music"
)

type MenuModel struct {
	common.Model
}

func InitialMenuModel() MenuModel {
	logger.Info("InitialMenuModel called")
	model := common.InitializeChoices("menu")
	model.Tick = common.Tick()
	return MenuModel{Model: model}
}

func UpdateMenu(m common.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	logger.Info("UpdateMenu called", "message", msg, "currentScreen", m.CurrentScreen)
	switch msg := msg.(type) {
	case common.TickMsg:
		m.Frame = (m.Frame + 1) % len(common.Frames)
		m.Banner = common.Frames[m.Frame]
		logger.Info("Banner updated", "banner", m.Banner)
		return m, m.Tick
	case tea.WindowSizeMsg:
		m.WindowWidth = msg.Width
		m.WindowHeight = msg.Height
		logger.Info("Window size updated", "width", m.WindowWidth, "height", m.WindowHeight)
		return m, nil
	case tea.KeyMsg:
		logger.Info("Key pressed", "key", msg.String(), "currentScreen", m.CurrentScreen)
		switch msg.String() {
		case "esc":
			logger.Info("quitting application", "currentScreen", m.CurrentScreen)
			return m, tea.Quit
		case "alt+end":
			logger.Info("Toggling music mute/unmute", "currentScreen", m.CurrentScreen)
			music.ToggleMusicMuteUnmute()
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			}
		case "enter":
			m.Selected = m.Choices[m.Cursor]
			logger.Info("Option selected", "selected", m.Selected)
			switch m.Selected {
			case "Setup (for Admins Only)":
				return m, func() tea.Msg { return "setup" }
			case "Select category & Start Quiz":
				return m, func() tea.Msg { return "start_quiz" }
			case "Create/Edit Questionnaire & Answers":
				logger.Info("Transitioning to Edit Questionnaire")
				editQuestionnaireModel := InitialEditQuestionnaireModel()
				return editQuestionnaireModel.Model, editQuestionnaireModel.Init()
			case "View Leaderboard score":
				return m, func() tea.Msg { return "leaderboard" }
			case "Submit Enhancement Request":
				return m, func() tea.Msg { return "feedback" }
			case "Exit":
				return m, tea.Quit
			}
		}
	}
	return m, nil
}
