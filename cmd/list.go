package cmd

import (
	"fmt"
	"gokafkaconnect/config"
	"gokafkaconnect/connector"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// listCmd represent command for retrieving connectors from API
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List running connector",
	Long:  `List current running connector`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			color.Red("Failed to load config: %v\n", err)
			return
		}
		connectors, err := connector.ListConnectors(cfg.KafkaConnectURL)
		if err != nil {
			color.Red("Failed to list connector: %v\n", err)
			return
		}

		if len(connectors) == 0 {
			color.Yellow("No connectors found")
			return
		}
		color.Cyan("🔗 Connectors:")
		for _, connector := range connectors {
			fmt.Printf("\t%s\n", connector)
		}

		var selected string
		prompt := &survey.Select{
			Message: "show connector config:",
			Options: connectors,
		}
		err = survey.AskOne(prompt, &selected)
		if err != nil {
			color.Red("canceled\n", err)
		}

		config, err := connector.GetConnectorConfig(cfg.KafkaConnectURL, selected)
		if err != nil {
			color.Red("Failed to get connector config: %v\n", err)
			return
		}
		color.Green("config for %s connector:\n", selected)
		fmt.Println(config)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

}
