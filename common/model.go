package common

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"letsquiz/logger"
)

type Model struct {
	CurrentScreen      string
	Cursor             int
	Choices            []string
	Selected           string
	Tick               tea.Cmd
	Banner             string
	Frame              int
	WindowWidth        int
	WindowHeight       int
	ButtonPos          []Rect
	Table              table.Model
	Buttons            []Button
	Form               *huh.Form
	ConfirmationDialog ConfirmationDialog
	Focused            string
	IsPlaying          bool
}

func (m Model) Init() tea.Cmd {
	logger.Info("Init called")
	return m.Tick
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	logger.Info("Update called")
	return m, nil
}

func (m Model) View() string {
	logger.Info("View called with CurrentScreen", "screen", m.CurrentScreen)
	return "Debug: View called with CurrentScreen = " + m.CurrentScreen
}
