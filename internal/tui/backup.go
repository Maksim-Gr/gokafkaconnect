package tui

import (
	"gokafkaconnect/internal/connector"
	"gokafkaconnect/internal/util"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type backupDoneMsg struct {
	file string
	err  error
}

func runBackup(dir string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(10 * time.Second)
		cfg, err := util.LoadConfig()
		if err != nil {
			return commandDoneMsg{err: err}
		}
		client := connector.NewClient(cfg.KafkaConnect.URL)
		connectors, err := client.ListConnectors()
		if err != nil {
			return commandDoneMsg{err: err}
		}

		file, err := connector.BackupConnectorConfig(client, connectors, dir)
		return commandDoneMsg{result: file, err: err}
	}
}
