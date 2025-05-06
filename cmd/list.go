package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gokafkaconnect/config"
	"gokafkaconnect/connectorconfig"
)

// listCmd represent command for retrieving connectors from API
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list available connectors",
	Long:  `List running connectors`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			color.Red("Failed to load config: %v\n", err)
			return
		}
		data, err := connectorconfig.ListConnectors(cfg.KafkaConnectURL)
		if err != nil {
			color.Red("Failed to list connectors: %v\n", err)
			return
		}
		var connectors []string
		err = json.Unmarshal(data, &connectors)

		color.Cyan("🔗 Connectors:")
		for _, connector := range connectors {
			fmt.Printf("\t%s\n", connector)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

}
