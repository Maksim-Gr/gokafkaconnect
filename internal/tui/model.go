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
	title := titleStyle.Render("Kafka Connect Backup")

	var statusLine string
	if m.loading {
		statusLine = fmt.Sprintf("%s %s", m.spinner.View(), m.status)
	} else {
		statusLine = m.status
	}

	if m.err != nil {
		statusLine = errorStyle.Render(m.err.Error())
	}

	panel := panelStyle.Render(
		fmt.Sprintf("Status:\n\n%s", statusLine),
	)

	footer := footerStyle.Render("[B] Backup   [Q] Quit")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		panel,
		"",
		footer,
	)
}
