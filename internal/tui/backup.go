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

		connectors, err := connector.ListConnectors(cfg.KafkaConnect.URL)
		if err != nil {
			return commandDoneMsg{err: err}
		}

		file, err := connector.BackupConnectorConfig(cfg.KafkaConnect.URL, connectors, dir)
		return commandDoneMsg{result: file, err: err}
	}
}
