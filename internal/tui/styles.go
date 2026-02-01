package tui

import "github.com/charmbracelet/lipgloss"

// Gruvbox (dark) palette (approx)
var (
	gbBg0   = lipgloss.Color("#282828")
	gbBg1   = lipgloss.Color("#3c3836")
	gbBg2   = lipgloss.Color("#504945")
	gbFg0   = lipgloss.Color("#fbf1c7")
	gbFg1   = lipgloss.Color("#ebdbb2")
	gbFgDim = lipgloss.Color("#a89984")

	gbRed    = lipgloss.Color("#fb4934")
	gbGreen  = lipgloss.Color("#b8bb26")
	gbYellow = lipgloss.Color("#fabd2f")
	gbBlue   = lipgloss.Color("#83a598")
	gbAqua   = lipgloss.Color("#8ec07c")
	gbPurple = lipgloss.Color("#d3869b")
)

var (
	// App chrome
	appStyle = lipgloss.NewStyle().
			Background(gbBg0).
			Foreground(gbFg1)

	headerBarStyle = lipgloss.NewStyle().
			Background(gbBg1).
			Foreground(gbFg1).
			Padding(0, 1)

	headerTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(gbFg0)

	// Panels
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(gbBg2).
			Padding(1, 2).
			Foreground(gbFg1)

	panelTitleStyle = lipgloss.NewStyle().
			Foreground(gbAqua).
			Bold(true)

	// Statusline (k9s/vim-ish)
	statusBarStyle = lipgloss.NewStyle().
			Background(gbBg1).
			Foreground(gbFg1).
			Padding(0, 1)

	statusKeyStyle = lipgloss.NewStyle().
			Background(gbBg2).
			Foreground(gbFg0).
			Bold(true).
			Padding(0, 1)

	statusSepStyle = lipgloss.NewStyle().
			Foreground(gbBg2)

	statusIdleStyle = lipgloss.NewStyle().
			Foreground(gbFgDim)

	statusRunningStyle = lipgloss.NewStyle().
				Foreground(gbYellow)

	statusSuccessStyle = lipgloss.NewStyle().
				Foreground(gbGreen)

	statusErrorStyle = lipgloss.NewStyle().
				Foreground(gbRed)
)
