package connector

import (
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
		if cfg.KafkaConnect.Username != "" {
			client.SetBasicAuth(cfg.KafkaConnect.Username, cfg.KafkaConnect.Password)
		}
		connectorStatuses, err := client.ListConnectorStatuses()
		if err != nil {
			color.Red("Failed to list connector statuses: %v", err)
			return
		}

		color.Cyan("Connector Statuses:")
		for name, status := range connectorStatuses {
			fmt.Printf("\t%s - Connector: %s\n", name, status.Connector.State)
			for _, t := range status.Tasks {
				fmt.Printf("\t\tTask %d: %s\n", t.ID, t.State)
			}
		}
	},
}
