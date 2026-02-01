package tui

import (
	"fmt"

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

	case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c":
			return m, tea.Quit

		case "b":
			m.loading = true
			m.status = "Backing up..."
			m.err = nil
			return m, runBackup(m.backupDir)
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
	header := titleStyle.Render(" Kafka Connect Backup ")

	divider := lipgloss.NewStyle().
		Foreground(lipgloss.Color("238")).
		Render("────────────────────────────────────────────")

	var status string
	switch {
	case m.err != nil:
		status = statusErrorStyle.Render(m.err.Error())
	case m.loading:
		status = statusRunningStyle.Render(
			fmt.Sprintf("%s %s", m.spinner.View(), m.status),
		)
	case m.status == "Backup created":
		status = statusSuccessStyle.Render(m.status)
	default:
		status = statusIdleStyle.Render(m.status)
	}

	body := panelStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			"Status:",
			"",
			status,
		),
	)

	footer := footerStyle.Render(" b:Backup   q:Quit ")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		divider,
		body,
		divider,
		footer,
	)
}
