// Package task provides CLI commands for managing Kafka Connect tasks.
package task

import (
	"encoding/json"
	"fmt"

	"github.com/Maksim-Gr/kkon/internal/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get task status",
	Long:  "Fetches status for a single task (select interactively or use --connector and --id).",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := util.NewKafkaConnectClient()
		if err != nil {
			return err
		}

		name, err := util.ResolveConnectorName(cmd.Context(), client, connectorName)
		if err != nil {
			return err
		}

		isDryRun := util.IsDryRun(cmd)
		id, err := util.ResolveTaskID(cmd.Context(), client, name, taskID, isDryRun)
		if err != nil {
			return err
		}

		if isDryRun {
			color.Yellow("[dry-run] Would get status for %s\n", util.FormatTaskRef(name, id))
			return nil
		}

		status, err := client.GetConnectorTaskStatus(cmd.Context(), name, id)
		if err != nil {
			return fmt.Errorf("failed to get status for %s: %w", util.FormatTaskRef(name, id), err)
		}

		if util.IsJSONOutput(cmd) {
			b, _ := json.MarshalIndent(status, "", "  ")
			fmt.Println(string(b))
			return nil
		}

		color.Cyan("Task status:")
		fmt.Printf("\tConnector: %s\n", name)
		fmt.Printf("\tTask ID:   %d\n", status.ID)
		fmt.Printf("\tState:     %s\n", util.ColorState(status.State))
		fmt.Printf("\tWorker:    %s\n", status.WorkerID)
		if status.Trace != "" {
			color.Yellow("\tTrace:\n%s\n", status.Trace)
		}
		return nil
	},
}

func init() {
	Cmd.AddCommand(getCmd)
}
