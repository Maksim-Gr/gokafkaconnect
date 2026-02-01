package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	mode viewMode

	status    string
	result    string
	err       error
	loading   bool
	backupDir string

	spinner spinner.Model

	width  int
	height int
}

func New(backupDir string) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return Model{
		mode:      viewHome,
		status:    "Ready",
		backupDir: backupDir,
		spinner:   s,
	}
}

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Update terminal size
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	// Handle key presses
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC: // Ctrl+C quits
			return m, tea.Quit
		}

		switch msg.String() {
		case "q":
			return m, tea.Quit

		case "esc", "\x1b":
			m.mode = viewHome
			m.loading = false
			m.err = nil
			m.status = "Ready"
			return m, nil

		case "b":
			m.mode = viewBackup
			m.loading = true
			m.status = "Running backup..."
			m.err = nil
			return m, tea.Batch(
				m.spinner.Tick,
				runBackup(m.backupDir),
			)

		case "l":
			m.mode = viewList
			m.loading = true
			m.status = "Fetching connectors..."
			m.err = nil
			return m, tea.Batch(
				m.spinner.Tick,
				runListConnectors(),
			)

		case "h":
			m.mode = viewHealthcheck
			m.loading = true
			m.status = "Checking health..."
			m.err = nil
			return m, tea.Batch(
				m.spinner.Tick,
				runHealthcheck(),
			)
		}
	}

	if msg, ok := msg.(spinner.TickMsg); ok {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	if msg, ok := msg.(commandDoneMsg); ok {
		m.loading = false
		m.mode = viewResult
		m.err = msg.err
		m.result = msg.result

		// Optional: colorize status line based on result
		if m.err != nil {
			m.status = statusErrorStyle.Render("Error")
		} else {
			m.status = statusSuccessStyle.Render("Done")
		}
	}

	return m, nil
}

func (m Model) View() string {
	switch m.mode {
	case viewHome:
		return m.viewHome()
	case viewBackup:
		return m.viewBackup()
	case viewList:
		return m.viewList()
	case viewHealthcheck:
		return m.viewHealthcheck()
	case viewResult:
		return m.viewResult()
	default:
		return "unknown view"
	}
}

func (m Model) viewHome() string {
	body := `
[B] Backup
[L] List connectors
[H] Healthcheck
[C] Configure
[D] Delete
[N] Create

[Q] Quit
`
	return m.layout("Kafka Connect", body)
}

func (m Model) viewHealthcheck() string {
	if m.loading {
		return m.layout("Healthcheck", m.spinner.View()+" "+statusRunningStyle.Render(m.status))
	}
	return m.layout("Healthcheck", m.result)
}

func (m Model) viewBackup() string {
	if m.loading {
		return m.layout("Backup", m.spinner.View()+" "+statusRunningStyle.Render(m.status))
	}
	return m.layout("Backup", m.status)
}

func (m Model) viewList() string {
	if m.loading {
		return m.layout("List", m.spinner.View()+" "+statusRunningStyle.Render(m.status))
	}
	return m.layout("List", m.result)
}

func (m Model) viewResult() string {
	if m.err != nil {
		return m.layout("Error", m.err.Error())
	}
	return m.layout("Result", m.result)
}

func (m Model) layout(title, body string) string {
	w := m.width
	h := m.height
	if w <= 0 {
		w = 80
	}
	if h <= 0 {
		h = 24
	}

	header := headerBarStyle.Width(w).Render(
		headerTitleStyle.Render(" " + title + " "),
	)

	left := lipgloss.JoinHorizontal(
		lipgloss.Left,
		statusKeyStyle.Render("STATUS"),
		statusSepStyle.Render(" | "),
		" "+m.status,
	)
	right := lipgloss.JoinHorizontal(
		lipgloss.Left,
		statusKeyStyle.Render("ESC"),
		" Back  ",
		statusKeyStyle.Render("Q"),
		" Quit",
	)

	spacerW := w - lipgloss.Width(left) - lipgloss.Width(right)
	if spacerW < 1 {
		right = statusKeyStyle.Render("Q") + " Quit"
		spacerW = w - lipgloss.Width(left) - lipgloss.Width(right)
		if spacerW < 1 {
			spacerW = 1
		}
	}
	footer := statusBarStyle.Width(w).Render(left + strings.Repeat(" ", spacerW) + right)

	panelW := w - 2
	if panelW < 30 {
		panelW = 30
	}
	panel := panelStyle.Width(panelW).Render(
		panelTitleStyle.Render(title) + "\n" + strings.TrimSpace(body),
	)
	panel = lipgloss.PlaceHorizontal(w, lipgloss.Left, panel)

	usedH := lipgloss.Height(header) + lipgloss.Height(panel) + lipgloss.Height(footer)
	padH := h - usedH
	if padH < 0 {
		padH = 0
	}
	padding := strings.Repeat("\n", padH)

	screen := header + "\n" + panel + padding + "\n" + footer
	return appStyle.Width(w).Height(h).Render(screen)
}
