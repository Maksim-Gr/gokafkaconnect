package connector

import (
	"context"

	"github.com/Maksim-Gr/kkon/internal/connector"

	"github.com/spf13/cobra"
)

// PauseCmd pauses one or more connectors.
var PauseCmd = &cobra.Command{
	Use:   "pause [name...]",
	Short: "Pause connectors",
	Long:  "Pauses Kafka Connect connectors and their tasks (multi-select interactively, pass connector names, or use --all).",
	Args:  cobra.ArbitraryArgs,
	RunE: lifecycleRunE(lifecycleSpec{
		verb:           "pause",
		successFmt:     "Pause requested for %s\n",
		confirmDefault: true,
		op: func(ctx context.Context, client *connector.Client, name string) error {
			return client.PauseConnector(ctx, name)
		},
	}),
}

func init() {
	PauseCmd.Flags().Bool("all", false, "Pause all connectors")
	PauseCmd.Flags().BoolP("yes", "y", false, "Skip the confirmation prompt for multiple connectors")
}
