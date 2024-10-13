package views

import (
	"github.com/charmbracelet/lipgloss"
	"letsquiz/logger"
	"letsquiz/models"
)

// ViewDynamicQuizForm renders the form for the current quiz question
func ViewDynamicQuizForm(m models.DynamicQuizModel) string {
	logger.Info("Rendering Dynamic Quiz Form View", "CurrentFormGroup", m.CurrentFormGroup)

	// Get current question form view from the model layer
	formView := m.QuestionForms[m.CurrentFormGroup].View()
	logger.Info("Form view rendered", "CurrentFormGroup", m.CurrentFormGroup)

	// Add footer message
	footerMessage := "Press Ctrl+Right Arrow to proceed, Ctrl+Left Arrow to go back, Esc to quit."
	footerStyle := lipgloss.NewStyle().
		Width(m.WindowWidth - 20). // Adjust width for boundary
		Render(footerMessage)

	// Combine form view and footer message
	content := lipgloss.JoinVertical(lipgloss.Top, formView, footerStyle)
	logger.Info("Combined form view and footer", "CurrentFormGroup", m.CurrentFormGroup)

	// Create the window boundary
	windowBoundary := lipgloss.NewStyle().
		Width(m.WindowWidth-20).   // Adjust width for the outer boundary
		Height(m.WindowHeight-10). // Adjust height for the outer boundary
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#04B575")).
		Padding(1, 1, 1, 1). // Adding padding to ensure the content is within the boundary
		Margin(1, 1, 1, 1).  // Adding margin to ensure the boundary doesn't exceed the window size
		Align(lipgloss.Center).
		Render(content)

	logger.Info("Window Boundary Applied", "width", m.WindowWidth-20, "height", m.WindowHeight-10)

	// Render the final view with the window boundary
	finalView := lipgloss.NewStyle().
		Width(m.WindowWidth).
		Height(m.WindowHeight).
		Align(lipgloss.Center).
		Render(windowBoundary)

	logger.Info("Final View Rendered", "width", m.WindowWidth, "height", m.WindowHeight)

	return finalView
}
