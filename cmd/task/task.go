package task

import "github.com/spf13/cobra"

var (
	connectorName string
	taskID        int
)

// Cmd is the root command for task management.
var Cmd = &cobra.Command{
	Use:   "task",
	Short: "Manage Kafka Connect tasks",
	Long: `List, inspect, and restart the tasks of a connector. All subcommands accept
--connector and --id; omit them to select interactively.`,
	Example: `  # List tasks for a connector
  kkon task list --connector my-connector

  # Get the status of a single task
  kkon task get --connector my-connector --id 0`,
}

func init() {
	Cmd.PersistentFlags().StringVarP(&connectorName, "connector", "c", "", "Connector name")
	Cmd.PersistentFlags().IntVarP(&taskID, "id", "i", -1, "Task id (integer)")
}
