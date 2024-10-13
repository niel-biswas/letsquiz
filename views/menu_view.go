package views

import (
	"github.com/charmbracelet/lipgloss"
	"letsquiz/common"
	"letsquiz/logger"
)

func ViewMenu(m common.Model) string {
	logger.Info("Rendering Menu View", "choices", m.Choices, "cursor", m.Cursor)

	buttonStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Align(lipgloss.Center).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#04B575"))

	selectedButtonStyle := buttonStyle.
		Foreground(lipgloss.Color("#FFA500")).
		BorderForeground(lipgloss.Color("#FFA500")).
		Bold(true)

	var buttons []string
	var totalHeight int
	maxWidth := 0

	for i, choice := range m.Choices {
		style := buttonStyle
		if m.Cursor == i {
			style = selectedButtonStyle
		}
		button := style.Render(choice)
		buttons = append(buttons, button)
		buttonHeight := lipgloss.Height(button)
		buttonWidth := lipgloss.Width(button)

		if buttonWidth > maxWidth {
			maxWidth = buttonWidth
		}

		totalHeight += buttonHeight
	}

	// Create a vertical stack of buttons
	buttonColumn := lipgloss.JoinVertical(lipgloss.Center, buttons...)

	// Create empty dummy banners for top, left, right, and bottom margins/paddings
	topBanner := lipgloss.NewStyle().
		Width(m.WindowWidth).
		Align(lipgloss.Center).
		Render("")

	leftBanner := lipgloss.NewStyle().
		Height(totalHeight).
		Align(lipgloss.Center).
		Render("")

	rightBanner := lipgloss.NewStyle().
		Height(totalHeight).
		Align(lipgloss.Center).
		Render("")

	bottomBanner := lipgloss.NewStyle().
		Width(m.WindowWidth).
		Align(lipgloss.Center).
		Render("")

	// Center the buttons within the window
	centeredButtonColumn := lipgloss.NewStyle().
		Width(m.WindowWidth-8).   // Increased width for more space within boundary
		Height(m.WindowHeight-8). // Increased height for more space within boundary
		Align(lipgloss.Center).
		Padding(10, 73).
		Render(buttonColumn)

	logger.Info("Button column rendered")

	// Add footer message
	footerMessage := "Press esc to quit, alt+end to mute/unmute."
	footerStyle := lipgloss.NewStyle().
		Align(lipgloss.Right).
		Width(m.WindowWidth - 8). // Adjusted width for boundary
		Render(footerMessage)

	// Create the content view
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		topBanner,
		lipgloss.JoinHorizontal(lipgloss.Center, leftBanner, centeredButtonColumn, rightBanner),
		bottomBanner,
		footerStyle,
	)

	// Create the window boundary
	windowBoundary := lipgloss.NewStyle().
		Width(m.WindowWidth - 4).   // Adjusted width for the outer boundary
		Height(m.WindowHeight - 4). // Adjusted height for the outer boundary
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
