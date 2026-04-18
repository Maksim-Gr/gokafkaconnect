// Package connector provides CLI commands for managing Kafka Connect connectors.
package connector

import (
	"gokafkaconnect/internal/connector"
	"gokafkaconnect/internal/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var backupDir string

// BackupCmd backs up connector configs from the Kafka Connect API.
var BackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup connectors config from Kafka Connect API",
	Long:  `Backup connectors config from Kafka Connect API and save to file for future usage `,
	Run: func(cmd *cobra.Command, _ []string) {
		cfg, err := util.LoadConfig()
		if err != nil {
			color.Red("Failed to load config: %v\n", err)
			return
		}
		client := connector.NewClient(cfg.KafkaConnect.URL)
		if cfg.KafkaConnect.Username != "" {
			client.SetBasicAuth(cfg.KafkaConnect.Username, cfg.KafkaConnect.Password)
		}

		connectors, err := client.ListConnectors(cmd.Context())
		if err != nil {
			color.Red("Failed to dump connector config: %v\n", err)
			return
		}
		backupFile, err := connector.BackupConnectorConfig(cmd.Context(), client, connectors, backupDir)
		if err != nil {
			color.Red("Failed to back up connectors config: %v\n", err)
			return
		}
		color.Green("Successfully created backup: %s\n", backupFile)
	},
}

func init() {
	BackupCmd.Flags().StringVarP(&backupDir, "dir", "o", "./backup", "Directory to save backup files")
}
