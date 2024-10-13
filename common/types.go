package common

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

type Button struct {
	Label string
}

type ConfirmationDialog struct {
	Message    string
	Active     bool
	Confirmed  bool
	ConfirmKey key.Binding
	CancelKey  key.Binding
}

// QuizMetaDataFormCompletedMsg is a message used to signal that a metadata form has been completed.
type QuizMetaDataFormCompletedMsg struct{}

// QuizDynamicFormsCompletedMsg is a message used to signal that all dynamic forms has been completed.
type QuizDynamicFormsCompletedMsg struct{}

func NewConfirmationDialog() ConfirmationDialog {
	return ConfirmationDialog{
		ConfirmKey: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "yes"),
		),
		CancelKey: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "no"),
		),
	}
}

func (d *ConfirmationDialog) Update(msg tea.Msg) (bool, bool) {
	if !d.Active {
		return false, false
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case d.ConfirmKey.Help().Key:
			d.Confirmed = true
			d.Active = false
			return true, true
		case d.CancelKey.Help().Key:
			d.Confirmed = false
			d.Active = false
			return true, false
		}
	}
	return false, false
}

func (d *ConfirmationDialog) View() string {
	if !d.Active {
		return ""
	}
	return "Do you really want to exit? (y/n)"
}
