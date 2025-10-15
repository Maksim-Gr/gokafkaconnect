package cmd

import (
	"fmt"
	"gokafkaconnect/config"
	"gokafkaconnect/connector"

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

		color.Cyan("🔗 Connectors:")
		for _, connector := range connectors {
			fmt.Printf("\t%s\n", connector)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

}
