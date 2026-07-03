package connector

import (
	"context"
	"fmt"

	"github.com/Maksim-Gr/kkon/internal/connector"

	"github.com/spf13/cobra"
)

var (
	restartIncludeTasks bool
	restartOnlyFailed   bool
)

// RestartCmd restarts one or more connectors.
var RestartCmd = &cobra.Command{
	Use:   "restart [name...]",
	Short: "Restart connectors",
	Long:  "Restarts Kafka Connect connectors (multi-select interactively, pass connector names, or use --all).",
	Args:  cobra.ArbitraryArgs,
	RunE: lifecycleRunE(lifecycleSpec{
		verb:           "restart",
		successFmt:     "Restart requested for %s\n",
		confirmSingle:  true,
		confirmDefault: true,
		dryRunNote: func() string {
			return fmt.Sprintf(" (includeTasks=%t, onlyFailed=%t)", restartIncludeTasks, restartOnlyFailed)
		},
		op: func(ctx context.Context, client *connector.Client, name string) error {
			return client.RestartConnector(ctx, name, restartIncludeTasks, restartOnlyFailed)
		},
	}),
}

func init() {
	RestartCmd.Flags().BoolVar(&restartIncludeTasks, "include-tasks", true, "Also restart the connector's tasks")
	RestartCmd.Flags().BoolVar(&restartOnlyFailed, "only-failed", false, "Restart only FAILED connector and tasks")
	RestartCmd.Flags().Bool("all", false, "Restart all connectors")
	RestartCmd.Flags().BoolP("yes", "y", false, "Skip the confirmation prompt")
}
