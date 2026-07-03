package connector

import (
	"context"

	"github.com/Maksim-Gr/kkon/internal/connector"

	"github.com/spf13/cobra"
)

// DeleteCmd represents the delete command.
var DeleteCmd = &cobra.Command{
	Use:   "delete [name...]",
	Short: "Delete connectors",
	Long:  "Delete connectors from Kafka Connect (multi-select interactively, pass connector names, or use --all).",
	Args:  cobra.ArbitraryArgs,
	RunE: lifecycleRunE(lifecycleSpec{
		verb:          "delete",
		successFmt:    "Connector %s deleted\n",
		confirmSingle: true,
		op: func(ctx context.Context, client *connector.Client, name string) error {
			return client.DeleteConnector(ctx, name)
		},
	}),
}

func init() {
	DeleteCmd.Flags().Bool("all", false, "Delete all connectors")
	DeleteCmd.Flags().BoolP("yes", "y", false, "Skip the confirmation prompt")
}
