package connector

import (
	"fmt"
	"gokafkaconnect/internal/connector"
	"gokafkaconnect/internal/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func stateColor(state string) string {
	switch state {
	case "RUNNING":
		return color.GreenString(state)
	case "FAILED":
		return color.RedString(state)
	default:
		return color.YellowString(state)
	}
}

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
		connectorStatuses, err := client.ListConnectorStatuses(cmd.Context())
		if err != nil {
			color.Red("Failed to list connector statuses: %v\n", err)
			return
		}

		maxLen := 0
		for name := range connectorStatuses {
			if len(name) > maxLen {
				maxLen = len(name)
			}
		}

		color.Cyan("Connector Statuses:")
		for name, status := range connectorStatuses {
			fmt.Printf("  %-*s  %s\n", maxLen, name, stateColor(status.Connector.State))
			for _, t := range status.Tasks {
				fmt.Printf("  %-*s    Task %d: %s\n", maxLen, "", t.ID, stateColor(t.State))
			}
		}
	},
}
