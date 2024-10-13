package views

import (
	"github.com/charmbracelet/lipgloss"
	"letsquiz/logger"
	"letsquiz/models"
)

// ViewQuizMetadata renders the form for QuizMetadata
func ViewQuizMetadata(m models.QuizMetadataModel) string {
	logger.Info("Rendering Quiz Metadata View")

	// Get huh form view from the model layer
	formView := m.Form.View()

	// Add footer message
	footerMessage := "Press esc to quit, alt+end to mute/unmute."

	footerStyle := lipgloss.NewStyle().
		Width(m.WindowWidth - 20). // Adjusted width for boundary
		Render(footerMessage)

	logger.Info("Footer Style Applied", "width", m.WindowWidth-20)

	// Combine form view and footer message
	content := lipgloss.JoinVertical(lipgloss.Top, formView, footerStyle)

	// Create the window boundary
	windowBoundary := lipgloss.NewStyle().
		Width(m.WindowWidth-20).   // Adjusted width for the outer boundary
		Height(m.WindowHeight-10). // Adjusted height for the outer boundary
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
