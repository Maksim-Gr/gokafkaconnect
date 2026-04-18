package task

import "github.com/spf13/cobra"

var (
	connectorName string
	taskID        int

	dryRun *bool
)

// Cmd is the root command for task management.
var Cmd = &cobra.Command{
	Use:   "task",
	Short: "Manage Kafka Connect tasks",
	Long:  "Task operations for Kafka Connect (list, get status, restart).",
}

// BindGlobals wires root-level flags into the task subcommands.
func BindGlobals(rootDryRun *bool) {
	dryRun = rootDryRun
}

func init() {
	Cmd.PersistentFlags().StringVarP(&connectorName, "connector", "c", "", "Connector name")
	Cmd.PersistentFlags().IntVarP(&taskID, "id", "i", -1, "Task id (integer)")
}
