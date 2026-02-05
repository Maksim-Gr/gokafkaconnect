package config

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management commands",
	Long:  `Manage Kafka Connect configuration including URL setup and backups.`,
}

func init() {
	Cmd.AddCommand(ConfigureCmd)
	Cmd.AddCommand(BackupCmd)
	Cmd.AddCommand(ShowConfigCmd)
}
