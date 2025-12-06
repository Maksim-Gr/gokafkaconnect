package cmd

import (
	"encoding/json"
	"fmt"

	"gokafkaconnect/internal/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// showConfigCmd represents the showConfig command
var showConfigCmd = &cobra.Command{
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

func init() {
	rootCmd.AddCommand(showConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showConfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showConfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
