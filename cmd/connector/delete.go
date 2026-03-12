package connector

import (
	"gokafkaconnect/internal/connector"
	"gokafkaconnect/internal/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// DeleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete connector",
	Long:  `Delete connector from Kafka Connect API`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := util.LoadConfig()
		if err != nil {
			color.Red("Failed to load config file: %v\n", err)
			return
		}
		client := connector.NewClient(cfg.KafkaConnect.URL)
		if cfg.KafkaConnect.Username != "" {
			client.SetBasicAuth(cfg.KafkaConnect.Username, cfg.KafkaConnect.Password)
		}

		connectors, err := client.ListConnectors()
		if err != nil {
			color.Red("Failed to list connectors: %v\n", err)
			return
		}
		if len(connectors) == 0 {
			color.Yellow("No connectors found")
			return
		}

		var selected string
		if err := survey.AskOne(&survey.Select{
			Message: "Select connector to delete:",
			Options: connectors,
		}, &selected); err != nil {
			color.Yellow("Canceled\n")
			return
		}

		var confirmed bool
		if err := survey.AskOne(&survey.Confirm{
			Message: "Delete " + selected + "?",
			Default: false,
		}, &confirmed); err != nil || !confirmed {
			color.Yellow("Canceled\n")
			return
		}

		if err := client.DeleteConnector(selected); err != nil {
			color.Red("Failed to delete connector: %v\n", err)
		} else {
			color.Green("Connector %s deleted\n", selected)
		}
	},
}

func init() {}
