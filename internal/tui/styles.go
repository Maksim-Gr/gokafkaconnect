package tui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("237")).
			Padding(0, 1).
			Width(60)

	panelStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Width(60)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("250")).
			Background(lipgloss.Color("237")).
			Padding(0, 1).
			Width(60)

	statusIdleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("250"))

	statusRunningStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("220"))

	statusSuccessStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("82"))

	statusErrorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("196"))
)
