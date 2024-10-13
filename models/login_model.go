package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"letsquiz/common"
	"letsquiz/logger"
	"letsquiz/music"
)

type LoginModel struct {
	common.Model
}

func InitialLoginModel() LoginModel {
	logger.Info("InitialLoginModel called")
	model := common.InitializeChoices("login")
	model.Tick = common.Tick()
	model.Banner = common.Frames[0]
	model.IsPlaying = true
	model.ConfirmationDialog = common.NewConfirmationDialog()
	return LoginModel{Model: model}
}

func UpdateLogin(m common.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	logger.Info("UpdateLogin called", "message", msg, "currentScreen", m.CurrentScreen)
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
			return m, nil
		case "left", "h":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "right", "l":
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			}
		case "enter":
			m.Selected = m.Choices[m.Cursor]
			if m.Selected == "Login via SSO" || m.Selected == "Signup via SSO" {
				logger.Info("Transitioning to menu")
				return InitialMenuModel().Model, nil
			}
			return m, tea.Quit
		}
	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			xPos, yPos := msg.X, msg.Y
			for i, rect := range m.ButtonPos {
				if xPos >= rect.X && xPos < rect.X+rect.Width && yPos >= rect.Y && yPos < rect.Y+rect.Height {
					m.Cursor = i
					m.Selected = m.Choices[m.Cursor]
					if m.Selected == "Login via SSO" || m.Selected == "Signup via SSO" {
						logger.Info("Transitioning to menu")
						return common.InitializeChoices("menu"), nil
					}
					return m, tea.Quit
				}
			}
		}
	}
	return m, nil
}
