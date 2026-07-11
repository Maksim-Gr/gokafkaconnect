package connector

import (
	"github.com/spf13/cobra"
)

// Cmd is the root command for connector management.
var Cmd = &cobra.Command{
	Use:   "connector",
	Short: "Manage Kafka Connect connectors",
	Long: `Create, update, delete, list, pause, resume, restart, back up, and restore
Kafka Connect connectors, and inspect the plugins installed on the cluster.`,
	Example: `  # Create a connector interactively
  kkon connector create

  # Show the status of every connector
  kkon connector health-check`,
}

func init() {
	Cmd.AddCommand(CreateCmd)
	Cmd.AddCommand(DeleteCmd)
	Cmd.AddCommand(ListCmd)
	Cmd.AddCommand(UpdateCmd)
	Cmd.AddCommand(HealthCheckCmd)
	Cmd.AddCommand(BackupCmd)
	Cmd.AddCommand(RestoreCmd)
	Cmd.AddCommand(PluginsCmd)
	Cmd.AddCommand(PauseCmd)
	Cmd.AddCommand(ResumeCmd)
	Cmd.AddCommand(RestartCmd)
}
