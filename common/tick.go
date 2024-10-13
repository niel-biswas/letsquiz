package common

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func Tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

type TickMsg time.Time
