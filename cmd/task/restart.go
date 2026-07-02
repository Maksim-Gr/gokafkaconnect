package task

import (
	"errors"
	"fmt"

	"github.com/Maksim-Gr/kkon/internal/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart a task",
	Long:  "Restarts a single Kafka Connect task (select interactively or use --connector and --id).",
	RunE: func(cmd *cobra.Command, _ []string) error {
		if util.IsJSONOutput(cmd) {
			return errors.New("--output json is not supported for task restart (interactive command)")
		}

		client, err := util.NewKafkaConnectClient()
		if err != nil {
			return err
		}

		name, err := util.ResolveConnectorName(cmd.Context(), client, connectorName)
		if err != nil {
			return err
		}

		isDryRun := util.IsDryRun(cmd)
		id, err := util.ResolveTaskID(cmd.Context(), client, name, taskID, isDryRun)
		if err != nil {
			return err
		}

		if isDryRun {
			color.Yellow("[dry-run] Would restart %s\n", util.FormatTaskRef(name, id))
			return nil
		}

		var confirm bool
		confirmPrompt := &survey.Confirm{
			Message: fmt.Sprintf("Restart %s?", util.FormatTaskRef(name, id)),
			Default: true,
		}
		if err := survey.AskOne(confirmPrompt, &confirm); err != nil || !confirm {
			return util.ErrCanceled
		}

		if err := client.RestartConnectorTask(cmd.Context(), name, id); err != nil {
			return fmt.Errorf("failed to restart %s: %w", util.FormatTaskRef(name, id), err)
		}

		color.Green("Restart requested for %s\n", util.FormatTaskRef(name, id))
		return nil
	},
}

func init() {
	Cmd.AddCommand(restartCmd)
}
