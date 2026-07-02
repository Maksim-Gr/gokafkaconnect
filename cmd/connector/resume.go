package connector

import (
	"context"

	"github.com/Maksim-Gr/kkon/internal/connector"

	"github.com/spf13/cobra"
)

// ResumeCmd resumes a connector.
var ResumeCmd = &cobra.Command{
	Use:   "resume [name]",
	Short: "Resume a connector",
	Long:  "Resumes a paused Kafka Connect connector and its tasks (select interactively or pass the connector name).",
	Args:  cobra.MaximumNArgs(1),
	RunE: lifecycleRunE("resume", func(ctx context.Context, client *connector.Client, name string) error {
		return client.ResumeConnector(ctx, name)
	}),
}
