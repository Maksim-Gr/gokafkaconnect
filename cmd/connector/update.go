package connector

import (
	"errors"
	"fmt"

	"github.com/Maksim-Gr/kkon/internal/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// UpdateCmd interactively updates an existing connector's configuration.
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing connector's configuration",
	Long: `Select a connector, edit its live configuration interactively, review the
changes, and apply them. This command is interactive and does not support
--output json.`,
	Example: `  # Pick a connector and edit its config
  kkon connector update

  # Preview the edit without applying it
  kkon connector update --dry-run`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		if util.IsJSONOutput(cmd) {
			return errors.New("--output json is not supported for update (interactive command)")
		}

		client, err := util.NewKafkaConnectClient()
		if err != nil {
			return err
		}

		connectors, err := client.ListConnectors(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to list connectors: %w", err)
		}
		if len(connectors) == 0 {
			color.Yellow("No connectors found\n")
			return nil
		}

		var selected string
		if err := survey.AskOne(&survey.Select{
			Message: "Select connector to update:",
			Options: connectors,
		}, &selected); err != nil {
			return util.ErrCanceled
		}

		return editConnectorConfig(cmd.Context(), client, selected, util.IsDryRun(cmd))
	},
}
