package connector

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "connector",
	Short: "Connector management commands",
	Long:  `Manage Kafka connectors including creation, deletion, listing, and health checks.`,
}

func init() {
	Cmd.AddCommand(CreateCmd)
	Cmd.AddCommand(DeleteCmd)
	Cmd.AddCommand(ListCmd)
	Cmd.AddCommand(HealthCheckCmd)
}
