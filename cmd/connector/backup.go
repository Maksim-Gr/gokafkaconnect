// Package connector provides CLI commands for managing Kafka Connect connectors.
package connector

import (
	"fmt"

	"github.com/Maksim-Gr/kkon/internal/connector"
	"github.com/Maksim-Gr/kkon/internal/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var backupDir string

// BackupCmd backs up connector configs from the Kafka Connect API.
var BackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup connectors config from Kafka Connect API",
	Long:  `Backup connectors config from Kafka Connect API and save to file for future usage `,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := util.NewKafkaConnectClient()
		if err != nil {
			return err
		}

		stop := util.StartSpinner("Backing up connectors...")
		connectors, err := client.ListConnectors(cmd.Context())
		if err != nil {
			stop()
			return fmt.Errorf("failed to list connectors: %w", err)
		}

		if util.IsDryRun(cmd) {
			stop()
			color.Yellow("[dry-run] Would back up %d connector(s) to %s\n", len(connectors), backupDir)
			return nil
		}

		backupFile, err := connector.BackupConnectorConfig(cmd.Context(), client, connectors, backupDir)
		stop()
		if err != nil {
			return fmt.Errorf("failed to back up connectors config: %w", err)
		}
		color.Green("Successfully backed up %d connector(s) → %s\n", len(connectors), backupFile)
		return nil
	},
}

func init() {
	BackupCmd.Flags().StringVar(&backupDir, "dir", "./backup", "Directory to save backup files")
}
