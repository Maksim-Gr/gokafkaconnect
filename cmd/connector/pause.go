package connector

import (
	"context"

	"github.com/Maksim-Gr/kkon/internal/connector"

	"github.com/spf13/cobra"
)

// PauseCmd pauses a connector.
var PauseCmd = &cobra.Command{
	Use:   "pause [name]",
	Short: "Pause a connector",
	Long:  "Pauses a Kafka Connect connector and its tasks (select interactively or pass the connector name).",
	Args:  cobra.MaximumNArgs(1),
	RunE: lifecycleRunE("pause", func(ctx context.Context, client *connector.Client, name string) error {
		return client.PauseConnector(ctx, name)
	}),
}
