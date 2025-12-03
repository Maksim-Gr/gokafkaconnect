/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"gokafkaconnect/internal/connector"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var outputFile string

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump connectors config from Kafka Connect API",
	Long:  `Dump connectors config from Kafka Connect API and save to file for future usage `,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := LoadConfig()
		if err != nil {
			color.Red("Failed to load config: %v\n", err)
			return
		}
		connectors, err := connector.ListConnectors(cfg.KafkaConnect.URL)
		if err != nil {
			color.Red("Failed to dump  connector config: %v\n", err)
			return
		}
		err = connector.DumpConnectorConfig(cfg.KafkaConnect.URL, connectors, outputFile)
		if err != nil {
			color.Red("Failed to dump connector config: %v\n", err)
			return
		}
		color.Green("Successfully dumped connector config to config.json")
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dumpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dumpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	dumpCmd.Flags().StringVarP(&outputFile, "output", "o", "current_config.json", "Output file for dumped configurations")
}
