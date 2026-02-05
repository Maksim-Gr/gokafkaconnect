package config

import (
	"gokafkaconnect/internal/connector"
	"gokafkaconnect/internal/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var backupDir string

var BackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup connectors config from Kafka Connect API",
	Long:  `Backup connectors config from Kafka Connect API and save to file for future usage `,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := util.LoadConfig()
		if err != nil {
			color.Red("Failed to load config: %v\n", err)
			return
		}
		client := connector.NewClient(cfg.KafkaConnect.URL)

		connectors, err := client.ListConnectors()
		if err != nil {
			color.Red("Failed to dump  connector config: %v\n", err)
			return
		}
		backupFile, err := connector.BackupConnectorConfig(client, connectors, backupDir)
		if err != nil {
			color.Red("Failed to back up  connectors config: %v\n", err)
			return
		}
		color.Green("Successfully created backup: %s", backupFile)
	},
}

func init() {
	BackupCmd.Flags().StringVarP(&backupDir, "dir", "o", "./backup", "Directory to save backup files")
}
