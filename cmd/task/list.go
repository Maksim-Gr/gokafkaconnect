package task

import (
	"encoding/json"
	"fmt"

	"github.com/Maksim-Gr/kkon/internal/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks for a connector",
	Long: `List a connector's tasks and their states. Prompts for the connector when
--connector is not given.`,
	Example: `  # Pick a connector and list its tasks
  kkon task list

  # List tasks for a specific connector
  kkon task list --connector my-connector

  # Machine-readable output
  kkon task list --connector my-connector --output json`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := util.NewKafkaConnectClient()
		if err != nil {
			return err
		}

		name, err := util.ResolveConnectorName(cmd.Context(), client, connectorName)
		if err != nil {
			return err
		}

		stop := util.StartSpinner("Fetching tasks...")
		tasks, err := client.ListConnectorTasks(cmd.Context(), name)
		connStatus, _ := client.GetConnectorStatus(cmd.Context(), name)
		stop()

		if err != nil {
			return fmt.Errorf("failed to list tasks for %s: %w", name, err)
		}

		if util.IsJSONOutput(cmd) {
			type taskJSON struct {
				Connector string `json:"connector"`
				Task      int    `json:"task"`
				State     string `json:"state,omitempty"`
			}
			taskStates := make(map[int]string, len(connStatus.Tasks))
			for _, ts := range connStatus.Tasks {
				taskStates[ts.ID] = ts.State
			}
			out := make([]taskJSON, 0, len(tasks))
			for _, t := range tasks {
				e := taskJSON{Connector: t.Connector, Task: t.Task}
				if state, ok := taskStates[t.Task]; ok {
					e.State = state
				}
				out = append(out, e)
			}
			b, _ := json.MarshalIndent(out, "", "  ")
			fmt.Println(string(b))
			return nil
		}

		if len(tasks) == 0 {
			color.Yellow("No tasks found for %s\n", name)
			return nil
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
		return nil
	},
}

func init() {
	Cmd.AddCommand(listCmd)
}
