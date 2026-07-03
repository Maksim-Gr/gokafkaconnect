package connector

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Maksim-Gr/kkon/internal/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var healthState string

// HealthCheckCmd shows connector and task statuses.
var HealthCheckCmd = &cobra.Command{
	Use:   "health-check",
	Short: "Show connector statuses",
	Long:  `Show each connector status`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := util.NewKafkaConnectClient()
		if err != nil {
			return err
		}

		stop := util.StartSpinner("Fetching connector statuses...")
		connectorStatuses, err := client.ListConnectorStatuses(cmd.Context())
		stop()
		if err != nil {
			return fmt.Errorf("failed to list connector statuses: %w", err)
		}

		if healthState != "" {
			for name, status := range connectorStatuses {
				if !strings.EqualFold(status.Connector.State, healthState) {
					delete(connectorStatuses, name)
				}
			}
		}

		if util.IsJSONOutput(cmd) {
			b, _ := json.MarshalIndent(connectorStatuses, "", "  ")
			fmt.Println(string(b))
			return nil
		}

		color.Cyan("Connector Statuses:")
		for name, status := range connectorStatuses {
			printConnectorStatus(cmd.Context(), client, name, status)
		}
		return nil
	},
}

func init() {
	HealthCheckCmd.Flags().StringVar(&healthState, "state", "", "Only show connectors in this state (e.g. RUNNING, FAILED, PAUSED)")
}
