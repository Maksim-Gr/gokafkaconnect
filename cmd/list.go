package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gokafkaconnect/config"
	"gokafkaconnect/connectorconfig"
)

// listCmd represent command for retrieving connectors from API
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list running connectors",
	Long:  `List current running connectors`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			color.Red("Failed to load config: %v\n", err)
			return
		}
		connectors, err := connectorconfig.ListConnectors(cfg.KafkaConnectURL)
		if err != nil {
			color.Red("Failed to list connectors: %v\n", err)
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
