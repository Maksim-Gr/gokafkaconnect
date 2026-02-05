package connector

import (
	"encoding/json"
	"fmt"
	"gokafkaconnect/internal/connector"
	"gokafkaconnect/internal/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// health-check show statuses for connectors
var HealthCheckCmd = &cobra.Command{
	Use:   "health-check",
	Short: "Show connector statuses",
	Long:  `Show each connector status`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := util.LoadConfig()
		if err != nil {
			color.Red("Failed to load config: %v\n", err)
			return
		}
		client := connector.NewClient(cfg.KafkaConnect.URL)
		rawStatuses, err := client.ListConnectorStatuses()
		if err != nil {
			color.Red("Failed to list connector statuses: %v", err)
			return
		}

		rawJSON, err := json.Marshal(rawStatuses)
		if err != nil {
			color.Red("Failed to marshal raw connector statuses: %v", err)
			return
		}

		var connectorStatuses connector.ConnectorsStatusResponse
		if err := json.Unmarshal(rawJSON, &connectorStatuses); err != nil {
			color.Red("Failed to unmarshal into typed connector statuses: %v", err)
			return
		}

		color.Cyan("🔗 Connector Statuses:")
		for name, status := range connectorStatuses {
			fmt.Printf("\t%s - Status: %s\n", name, status.Connector.State)
		}
	},
}
