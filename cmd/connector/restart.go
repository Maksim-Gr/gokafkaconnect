package connector

import (
	"fmt"

	"github.com/Maksim-Gr/kkon/internal/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	restartIncludeTasks bool
	restartOnlyFailed   bool
	restartYes          bool
)

// RestartCmd restarts a connector.
var RestartCmd = &cobra.Command{
	Use:   "restart [name]",
	Short: "Restart a connector",
	Long:  "Restarts a Kafka Connect connector (select interactively or pass the connector name).",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if util.IsJSONOutput(cmd) && (argOrEmpty(args) == "" || !restartYes) {
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
			color.Yellow("[dry-run] Would restart connector %s (includeTasks=%t, onlyFailed=%t)\n",
				name, restartIncludeTasks, restartOnlyFailed)
			return nil
		}

		if !restartYes {
			var confirm bool
			confirmPrompt := &survey.Confirm{
				Message: fmt.Sprintf("Restart connector %s?", name),
				Default: true,
			}
			if err := survey.AskOne(confirmPrompt, &confirm); err != nil || !confirm {
				return util.ErrCanceled
			}
		}

		if err := client.RestartConnector(cmd.Context(), name, restartIncludeTasks, restartOnlyFailed); err != nil {
			return fmt.Errorf("failed to restart %s: %w", name, err)
		}
		if util.IsJSONOutput(cmd) {
			printActionJSON(name, "restart")
			return nil
		}
		color.Green("Restart requested for %s\n", name)
		return nil
	},
}

func init() {
	RestartCmd.Flags().BoolVar(&restartIncludeTasks, "include-tasks", true, "Also restart the connector's tasks")
	RestartCmd.Flags().BoolVar(&restartOnlyFailed, "only-failed", false, "Restart only FAILED connector and tasks")
	RestartCmd.Flags().BoolVarP(&restartYes, "yes", "y", false, "Skip the confirmation prompt")
}
