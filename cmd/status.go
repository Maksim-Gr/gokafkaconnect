package cmd

import (
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
		statuses, err := connectorconfig.ListConnectorStatuses(cfg.KafkaConnectURL)
		if err != nil {
			color.Red("Failed to list  connector's statuses: %v\n", err)
			return
		}
		color.Cyan("🔗 Connector Statuses:")
		for name, data := range statuses {
			statusMap := data.(map[string]interface{})
			status := statusMap["status"].(map[string]interface{})
			connectorState := status["connector"].(map[string]interface{})["state"]
			fmt.Printf("\t%s - Status: %v\n", name, connectorState)
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
