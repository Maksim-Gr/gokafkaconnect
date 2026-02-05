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
	Run: func(cmd *cobra.Command, args []string) {
		client, ok := util.NewKafkaConnectClient()
		if !ok {
			return
		}

		name, ok := util.ResolveConnectorName(client, connectorName)
		if !ok {
			return
		}

		if dryRun != nil && *dryRun {
			color.Yellow("[dry-run] Would list tasks for connector: %s\n", name)
			return
		}

		tasks, err := client.ListConnectorTasks(name)
		if err != nil {
			color.Red("Failed to list tasks for %s: %v\n", name, err)
			return
		}
		if len(tasks) == 0 {
			color.Yellow("No tasks found for %s\n", name)
			return
		}

		color.Cyan("Tasks for %s:", name)
		for _, t := range tasks {
			fmt.Printf("\t%d\n", t.Task)
		}
	},
}

func init() {
	Cmd.AddCommand(listCmd)
}
