package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	status    string
	backupDir string
	loading   bool
	spinner   spinner.Model
	err       error

	width  int
	height int
}

func New(backupDir string) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return Model{
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

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "b":
			m.loading = true
			m.status = "Backing up..."
			m.err = nil
			return m, tea.Batch(
				m.spinner.Tick,
				runBackup(m.backupDir),
			)
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case backupDoneMsg:
		m.loading = false
		if msg.err != nil {
			m.status = "Backup failed"
			m.err = msg.err
		} else {
			m.status = "Backup created"
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	header := titleStyle.
		Width(m.width).
		Render(" Kafka Connect Backup ")

	divider := lipgloss.NewStyle().
		Foreground(lipgloss.Color("238")).
		Render(strings.Repeat("─", m.width))

	var status string
	switch {
	case m.err != nil:
		status = statusErrorStyle.Render(m.err.Error())
	case m.loading:
		status = statusRunningStyle.Render(
			fmt.Sprintf("%s %s", m.spinner.View(), m.status),
		)
	default:
		status = statusIdleStyle.Render(m.status)
	}

	body := panelStyle.
		Width(m.width).
		Height(m.height - 5).
		Render(
			lipgloss.Place(
				m.width,
				m.height-5,
				lipgloss.Left,
				lipgloss.Top,
				"Status:\n\n"+status,
			),
		)

	footer := footerStyle.
		Width(m.width).
		Render(" b:backup   q:quit ")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		divider,
		body,
		divider,
		footer,
	)
}
