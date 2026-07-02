package connector

import (
	"context"

	"github.com/Maksim-Gr/kkon/internal/connector"

	"github.com/spf13/cobra"
)

// ResumeCmd resumes one or more connectors.
var ResumeCmd = &cobra.Command{
	Use:   "resume [name...]",
	Short: "Resume connectors",
	Long:  "Resumes paused Kafka Connect connectors and their tasks (multi-select interactively, pass connector names, or use --all).",
	Args:  cobra.ArbitraryArgs,
	RunE: lifecycleRunE(lifecycleSpec{
		verb:           "resume",
		successFmt:     "Resume requested for %s\n",
		confirmDefault: true,
		op: func(ctx context.Context, client *connector.Client, name string) error {
			return client.ResumeConnector(ctx, name)
		},
	}),
}

func init() {
	ResumeCmd.Flags().Bool("all", false, "Resume all connectors")
	ResumeCmd.Flags().BoolP("yes", "y", false, "Skip the confirmation prompt for multiple connectors")
}
