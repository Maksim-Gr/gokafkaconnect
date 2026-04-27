package task

import (
	"fmt"

	"gokafkaconnect/internal/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks for a connector",
	Long:  "Lists tasks for a selected connector (or --connector).",
	Run: func(cmd *cobra.Command, _ []string) {
		client, ok := util.NewKafkaConnectClient()
		if !ok {
			return
		}

		name, ok := util.ResolveConnectorName(cmd.Context(), client, connectorName)
		if !ok {
			return
		}

		if dryRun != nil && *dryRun {
			color.Yellow("[dry-run] Would list tasks for connector: %s\n", name)
			return
		}

		stop := util.StartSpinner("Fetching tasks...")
		tasks, err := client.ListConnectorTasks(cmd.Context(), name)
		connStatus, _ := client.GetConnectorStatus(cmd.Context(), name)
		stop()

		if err != nil {
			color.Red("Failed to list tasks for %s: %v\n", name, err)
			return
		}
		if len(tasks) == 0 {
			color.Yellow("No tasks found for %s\n", name)
			return
		}

		taskStates := make(map[int]string, len(connStatus.Tasks))
		for _, ts := range connStatus.Tasks {
			taskStates[ts.ID] = ts.State
		}

		color.Cyan("Tasks for %s:", name)
		for _, t := range tasks {
			badge := ""
			if state, ok := taskStates[t.Task]; ok {
				badge = "  " + util.ColorState(state)
			}
			fmt.Printf("  Task %d%s\n", t.Task, badge)
		}
	},
}

func init() {
	Cmd.AddCommand(listCmd)
}
