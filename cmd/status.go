package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gokafkaconnect/config"
	"gokafkaconnect/connectorconfig"
)

// statCmd represents the stat command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Stat command return statuses for all connectors",
	Long:  `Stat command return statuses for all connectors for configured Kafka connect endpoint`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			color.Red("Failed to load config: %v\n", err)
			return
		}
		rawStatuses, err := connectorconfig.ListConnectorStatuses(cfg.KafkaConnectURL)
		if err != nil {
			color.Red("Failed to list connector statuses: %v", err)
			return
		}

		// Marshal map[string]interface{} back into JSON
		rawJSON, err := json.Marshal(rawStatuses)
		if err != nil {
			color.Red("Failed to marshal raw connector statuses: %v", err)
			return
		}

		// Unmarshal into typed ConnectorsStatusResponse
		var connectorStatuses connectorconfig.ConnectorsStatusResponse
		if err := json.Unmarshal(rawJSON, &connectorStatuses); err != nil {
			color.Red("Failed to unmarshal into typed connector statuses: %v", err)
			return
		}

		color.Cyan("🔗 Connector Statuses:")
		for name, status := range connectorStatuses {
			fmt.Printf("\t%s - Status: %s\n", name, status.Connector.State)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
