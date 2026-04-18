package connector

import (
	"fmt"
	"gokafkaconnect/internal/connector"
	"gokafkaconnect/internal/util"
	"sort"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// UpdateCmd interactively updates an existing connector's configuration.
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing connector configuration",
	Long:  `Fetch a connector's live config and edit fields interactively, then apply the changes.`,
	Run: func(cmd *cobra.Command, _ []string) {
		cfg, err := util.LoadConfig()
		if err != nil {
			color.Red("Failed to load config: %v\n", err)
			return
		}
		client := connector.NewClient(cfg.KafkaConnect.URL)
		if cfg.KafkaConnect.Username != "" {
			client.SetBasicAuth(cfg.KafkaConnect.Username, cfg.KafkaConnect.Password)
		}

		connectors, err := client.ListConnectors(cmd.Context())
		if err != nil {
			color.Red("Failed to list connectors: %v\n", err)
			return
		}
		if len(connectors) == 0 {
			color.Yellow("No connectors found\n")
			return
		}

		var selected string
		if err := survey.AskOne(&survey.Select{
			Message: "Select connector to update:",
			Options: connectors,
		}, &selected); err != nil {
			color.Yellow("Canceled\n")
			return
		}

		connectorConfig, err := client.GetConnectorConfigJSON(cmd.Context(), selected)
		if err != nil {
			color.Red("Failed to get connector config: %v\n", err)
			return
		}

		for {
			pretty, err := util.ToPrettyJSON(connectorConfig)
			if err != nil {
				color.Red("Failed to format config: %v\n", err)
				return
			}
			color.Cyan("\n Current config for %s:\n", selected)
			fmt.Println(pretty)

			fields := make([]string, 0, len(connectorConfig))
			for k := range connectorConfig {
				fields = append(fields, k)
			}
			sort.Strings(fields)

			var fieldToChange string
			if err := survey.AskOne(&survey.Select{
				Message: "Which field do you want to change?",
				Options: fields,
			}, &fieldToChange); err != nil {
				color.Yellow("Canceled\n")
				return
			}

			var newValue string
			if err := survey.AskOne(&survey.Input{
				Message: fmt.Sprintf("New value for %s (current: %v):", fieldToChange, connectorConfig[fieldToChange]),
			}, &newValue); err != nil {
				color.Yellow("Canceled\n")
				return
			}
			connectorConfig[fieldToChange] = newValue

			var more bool
			if err := survey.AskOne(&survey.Confirm{
				Message: "Change another field?",
				Default: false,
			}, &more); err != nil {
				color.Yellow("Canceled\n")
				return
			}
			if !more {
				break
			}
		}

		pretty, err := util.ToPrettyJSON(connectorConfig)
		if err != nil {
			color.Red("Failed to format config: %v\n", err)
			return
		}
		color.Cyan("\nFinal config:\n")
		fmt.Println(pretty)

		var confirm bool
		if err := survey.AskOne(&survey.Confirm{
			Message: "Apply this config to " + selected + "?",
			Default: true,
		}, &confirm); err != nil {
			color.Yellow("Canceled\n")
			return
		}
		if !confirm {
			color.Yellow("Canceled\n")
			return
		}

		if err := client.UpdateConnectorConfig(cmd.Context(), selected, connectorConfig); err != nil {
			color.Red("Failed to update connector: %v\n", err)
			return
		}
		color.Green("Connector %s updated successfully\n", selected)
	},
}
