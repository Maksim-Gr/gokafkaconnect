package tui

import tea "github.com/charmbracelet/bubbletea"

type Model struct {
	status string
}

func New() Model {
	return Model{
		status: "Ready",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	return "\n Kafka Connect Status: \n\n" + "Status: " + m.status + "\n\n" + "Press q to quit."
}
