package tui

import (
	"strings"

	"gokafkaconnect/internal/connector"
	"gokafkaconnect/internal/util"

	tea "github.com/charmbracelet/bubbletea"
)

func runListConnectors() tea.Cmd {
	return func() tea.Msg {
		cfg, err := util.LoadConfig()
		if err != nil {
			return commandDoneMsg{err: err}
		}

		connectors, err := connector.ListConnectors(cfg.KafkaConnect.URL)
		if err != nil {
			return commandDoneMsg{err: err}
		}

		if len(connectors) == 0 {
			return commandDoneMsg{result: "No connectors found"}
		}

		return commandDoneMsg{
			result: strings.Join(connectors, "\n"),
		}
	}
}
