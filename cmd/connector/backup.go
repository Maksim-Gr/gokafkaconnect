// Package connector provides CLI commands for managing Kafka Connect connectors.
package connector

import (
	"encoding/json"
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
	Short: "Back up all connector configs to a file",
	Long: `Fetch every connector's configuration from Kafka Connect and save it to a
JSON backup file for later restore with 'kkon connector restore'.`,
	Example: `  # Back up all connectors to ./backup
  kkon connector backup

  # Back up to a custom directory
  kkon connector backup --dir ./my-backups

  # Preview what would be backed up
  kkon connector backup --dry-run`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		jsonMode := util.IsJSONOutput(cmd)

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
			if jsonMode {
				b, _ := json.Marshal(map[string]any{"connectors": connectors, "dir": backupDir, "result": "dry-run"})
				fmt.Println(string(b))
				return nil
			}
			color.Yellow("[dry-run] Would back up %d connector(s) to %s\n", len(connectors), backupDir)
			return nil
		}

		backupFile, err := connector.BackupConnectorConfig(cmd.Context(), client, connectors, backupDir)
		stop()
		if err != nil {
			return fmt.Errorf("failed to back up connectors config: %w", err)
		}
		if jsonMode {
			b, _ := json.Marshal(map[string]any{"backedUp": len(connectors), "file": backupFile, "result": "ok"})
			fmt.Println(string(b))
			return nil
		}
		color.Green("Successfully backed up %d connector(s) → %s\n", len(connectors), backupFile)
		return nil
	},
}

func init() {
	BackupCmd.Flags().StringVar(&backupDir, "dir", "./backup", "Directory to save backup files")
}
