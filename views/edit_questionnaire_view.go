package views

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"letsquiz/logger"
	"letsquiz/models"
	"os"
)

const (
	purple    = lipgloss.Color("99")
	gray      = lipgloss.Color("245")
	lightGray = lipgloss.Color("241")
)

var (
	// HeaderStyle is the lipgloss style used for the table headers.
	HeaderStyle = lipgloss.NewRenderer(os.Stdout).NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
	// CellStyle is the base lipgloss style used for the table rows.
	CellStyle = lipgloss.NewRenderer(os.Stdout).NewStyle().Padding(0, 1).Width(14)
	// OddRowStyle is the lipgloss style used for odd-numbered table rows.
	OddRowStyle = CellStyle.Foreground(gray)
	// EvenRowStyle is the lipgloss style used for even-numbered table rows.
	EvenRowStyle = CellStyle.Foreground(lightGray)
	// BorderStyle is the lipgloss style used for the table border.
	BorderStyle = lipgloss.NewStyle().Foreground(purple)
)

func ViewEditQuestionnaire(m models.EditQuestionnaireModel) string {
	logger.Info("Rendering EditQuestionnaire View")

	// Get the default table styles
	style := table.DefaultStyles()

	// Conditionally style the table header text based on focus
	if m.Focused == "table" {
		style.Header = style.Header.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57"))
	} else {
		// Reset header style to default when the table is not focused
		style.Header = style.Header.
			Foreground(lipgloss.Color("240")).
			Background(lipgloss.Color("0"))
	}

	// Apply the styles to the table
	m.Table.SetStyles(style)

	// Render the button with the appropriate focus style
	var createButtonView string
	if m.Focused == "button" {
		focusedStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(true).
			Render
		createButtonView = focusedStyle(fmt.Sprintf("[ %s ]", m.Buttons[0].Label))
	} else {
		normalStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render
		createButtonView = normalStyle(fmt.Sprintf("[ %s ]", m.Buttons[0].Label))
	}

	// Render the button above the table
	view := "\n" + createButtonView + "\n\n" + m.Table.View()

	// Add footer message
	footerMessage := "Press esc to quit, alt+end to mute/unmute."
	footerStyle := lipgloss.NewStyle().
		Align(lipgloss.Right).
		Width(m.WindowWidth - 10). // Adjusted width for boundary
		Render(footerMessage)

	// Create the content view
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		view,
		footerStyle,
	)

	// Create the window boundary
	windowBoundary := lipgloss.NewStyle().
		Width(m.WindowWidth - 10).   // Adjusted width for the outer boundary
		Height(m.WindowHeight - 10). // Adjusted height for the outer boundary
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#04B575")).
		Align(lipgloss.Center).
		Render(content)

	// Render the final view with the window boundary
	finalView := lipgloss.NewStyle().
		Width(m.WindowWidth).
		Height(m.WindowHeight).
		Align(lipgloss.Center).
		Render(windowBoundary)

	return finalView
}
