package tui

import (
	"encoding/json"
	"fmt"
	"strings"

	"gokafkaconnect/internal/connector"
	"gokafkaconnect/internal/util"

	tea "github.com/charmbracelet/bubbletea"
)

// runHealthcheck fetches connector statuses asynchronously
func runHealthcheck() tea.Cmd {
	return func() tea.Msg {
		cfg, err := util.LoadConfig()
		if err != nil {
			return commandDoneMsg{err: err}
		}

		// Get raw statuses
		rawStatuses, err := connector.ListConnectorStatuses(cfg.KafkaConnect.URL)
		if err != nil {
			return commandDoneMsg{err: err}
		}

		rawJSON, err := json.Marshal(rawStatuses)
		if err != nil {
			return commandDoneMsg{err: err}
		}

		var connectorStatuses connector.ConnectorsStatusResponse
		if err := json.Unmarshal(rawJSON, &connectorStatuses); err != nil {
			return commandDoneMsg{err: err}
		}

		// Format for TUI output
		lines := []string{"🔗 Connector Statuses:"}
		for name, status := range connectorStatuses {
			lines = append(lines, fmt.Sprintf("  %s - Status: %s", name, status.Connector.State))
		}

		return commandDoneMsg{
			result: strings.Join(lines, "\n"),
		}
	}
}
