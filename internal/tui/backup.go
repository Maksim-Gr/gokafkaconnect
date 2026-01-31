package tui

import (
	"gokafkaconnect/internal/connector"
	"gokafkaconnect/internal/util"

	tea "github.com/charmbracelet/bubbletea"
)

type backupDoneMsg struct {
	file string
	err  error
}

func runBackup(dir string) tea.Cmd {
	return func() tea.Msg {
		cfg, err := util.LoadConfig()
		if err != nil {
			return backupDoneMsg{err: err}
		}
		connectors, err := connector.ListConnectors(cfg.KafkaConnect.URL)
		if err != nil {
			return backupDoneMsg{err: err}
		}
		file, err := connector.BackupConnectorConfig(cfg.KafkaConnect.URL, connectors, dir)

		return backupDoneMsg{file: file, err: err}
	}
}
