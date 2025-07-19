/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gokafkaconnect/connector"
)

var outputFile string

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := LoadConfig()
		if err != nil {
			color.Red("Failed to load config: %v\n", err)
			return
		}
		connectors, err := connector.ListConnectors(cfg.KafkaConnectURL)
		if err != nil {
			color.Red("Failed to dump  connector config: %v\n", err)
			return
		}
		err = connector.DumpConnectorConfig(cfg.KafkaConnectURL, connectors, outputFile)
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
