package connector

import (
	"fmt"

	"github.com/Maksim-Gr/kkon/internal/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var deleteYes bool

// DeleteCmd represents the delete command.
var DeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a connector",
	Long:  "Delete a connector from Kafka Connect (select interactively or pass the connector name).",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if util.IsJSONOutput(cmd) && (argOrEmpty(args) == "" || !deleteYes) {
			return errNonInteractiveJSON
		}

		client, err := util.NewKafkaConnectClient()
		if err != nil {
			return err
		}

		name, err := util.ResolveConnectorName(cmd.Context(), client, argOrEmpty(args))
		if err != nil {
			return err
		}

		if util.IsDryRun(cmd) {
			color.Yellow("[dry-run] Would delete connector %s\n", name)
			return nil
		}

		if !deleteYes {
			var confirmed bool
			if err := survey.AskOne(&survey.Confirm{
				Message: "Delete " + name + "?",
				Default: false,
			}, &confirmed); err != nil || !confirmed {
				return util.ErrCanceled
			}
		}

		if err := client.DeleteConnector(cmd.Context(), name); err != nil {
			return fmt.Errorf("failed to delete connector: %w", err)
		}
		if util.IsJSONOutput(cmd) {
			printActionJSON(name, "delete")
			return nil
		}
		color.Green("Connector %s deleted\n", name)
		return nil
	},
}

func init() {
	DeleteCmd.Flags().BoolVarP(&deleteYes, "yes", "y", false, "Skip the confirmation prompt")
}
