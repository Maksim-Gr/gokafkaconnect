package tui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7DD3FC")).
			Padding(0, 1)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#38BDF8")).
			Padding(1, 2).
			Width(50)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#94A3B8")).
			PaddingTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F87171")).
			Bold(true)
)
