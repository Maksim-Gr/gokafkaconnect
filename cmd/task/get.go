package task

import (
	"fmt"
	"gokafkaconnect/internal/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get task status",
	Long:  "Fetches status for a single task (select interactively or use --connector and --id).",
	Run: func(cmd *cobra.Command, args []string) {
		client, ok := util.NewKafkaConnectClient()
		if !ok {
			return
		}

		name, ok := util.ResolveConnectorName(client, connectorName)
		if !ok {
			return
		}

		isDryRun := dryRun != nil && *dryRun
		id, ok := util.ResolveTaskID(client, name, taskID, isDryRun)
		if !ok {
			return
		}

		if isDryRun {
			color.Yellow("[dry-run] Would get status for %s\n", util.FormatTaskRef(name, id))
			return
		}

		status, err := client.GetConnectorTaskStatus(name, id)
		if err != nil {
			color.Red("Failed to get status for %s: %v\n", util.FormatTaskRef(name, id), err)
			return
		}

		color.Cyan("Task status:")
		fmt.Printf("\tConnector: %s\n", name)
		fmt.Printf("\tTask ID:   %d\n", status.ID)
		fmt.Printf("\tState:     %s\n", status.State)
		fmt.Printf("\tWorker:    %s\n", status.WorkerID)
		if status.Trace != "" {
			color.Yellow("\tTrace:\n%s\n", status.Trace)
		}
	},
}

func init() {
	Cmd.AddCommand(getCmd)
}
