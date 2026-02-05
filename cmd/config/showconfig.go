package config

import (
	"encoding/json"
	"fmt"

	"gokafkaconnect/internal/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// ShowConfigCmd represents the showConfig command
var ShowConfigCmd = &cobra.Command{
	Use:   "show-config",
	Short: "Display API endpoint",
	Long:  `Display Kafka Connect API endpoint.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := util.LoadConfig()
		if err != nil {
			color.Red("Failed to load config: %v\n", err)
			return
		}

		color.Cyan("Current Configuration:")
		data, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			color.Red("Failed to format config: %v\n", err)
			return
		}
		fmt.Printf("\n%s\n\n", string(data))
	},
}
