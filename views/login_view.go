package views

import (
	"github.com/charmbracelet/lipgloss"
	"letsquiz/common"
	"letsquiz/logger"
)

func ViewLogin(m common.Model) string {
	logger.Info("Rendering Login View")
	logger.Info("Current Banner", "banner", m.Banner)

	bannerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575")).
		Align(lipgloss.Center).
		Width(m.WindowWidth - 4). // Adjusted width
		Render(m.Banner)

	logger.Info("Banner rendered")

	buttonStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Margin(3, 3).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#04B575"))

	selectedButtonStyle := buttonStyle.Copy().
		Foreground(lipgloss.Color("#FFA500")).
		BorderForeground(lipgloss.Color("#FFA500")).
		Bold(true)

	var buttons []string
	var buttonWidths []int
	var buttonHeights []int
	totalWidth := 0

	for i, choice := range m.Choices {
		style := buttonStyle
		if m.Cursor == i {
			style = selectedButtonStyle
		}
		button := style.Render(choice)
		buttons = append(buttons, button)
		buttonWidth := lipgloss.Width(button)
		buttonHeight := lipgloss.Height(button)
		buttonWidths = append(buttonWidths, buttonWidth)
		buttonHeights = append(buttonHeights, buttonHeight)
		totalWidth += buttonWidth

		if i < len(m.Choices)-1 {
			totalWidth += 3
		}
	}

	logger.Info("Buttons rendered")

	xOffset := (m.WindowWidth - totalWidth) / 2
	yOffset := m.WindowHeight/2 - buttonHeights[0]/2
	for i, button := range buttons {
		m.ButtonPos[i] = common.Rect{
			X:      xOffset,
			Y:      yOffset,
			Width:  buttonWidths[i],
			Height: lipgloss.Height(button),
		}
		xOffset += buttonWidths[i] + 3
	}

	logger.Info("Button positions calculated")

	buttonRow := lipgloss.JoinHorizontal(lipgloss.Top, buttons...)
	buttonRowWidth := lipgloss.Width(buttonRow)

	buttonRowStyle := lipgloss.NewStyle().
		Width(buttonRowWidth).
		Align(lipgloss.Center).
		Render(buttonRow)

	logger.Info("Button row rendered")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		bannerStyle,
		lipgloss.NewStyle().Height(1).Render(""),
		buttonRowStyle,
		lipgloss.NewStyle().Width(m.WindowWidth-4).Align(lipgloss.Center).Render("Press esc to quit, alt+end to mute/unmute."),
	)

	windowBoundary := lipgloss.NewStyle().
		Width(m.WindowWidth - 2).   // Adjusted width
		Height(m.WindowHeight - 2). // Adjusted height
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#04B575")).
		Align(lipgloss.Center).
		Render(content)

	return lipgloss.NewStyle().
		Width(m.WindowWidth).
		Height(m.WindowHeight).
		Align(lipgloss.Center).
		Render(windowBoundary)
}
