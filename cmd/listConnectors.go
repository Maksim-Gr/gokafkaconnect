/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Available connectors
// TODO replace with config files
var connectors = []string{
	"🚀 PostgreSQL Connector",
	"💾 MySQL Connector",
	"📚 MongoDB Connector",
	"⚡ Redis Connector",
	"📨 Kafka Source Connector",
}

// listConnectorsCmd represents the list-connectors command
var listConnectorsCmd = &cobra.Command{
	Use:   "list-connectors",
	Short: "List and select available connectors 🔥",
	Long:  `Browse predefined connectors.`,
	Run: func(cmd *cobra.Command, args []string) {
		var selected string
		color.Cyan("\n✨ Available Kafka Connectors ✨\n")
		prompt := &survey.Select{
			Message: "Pick a connector to work with:",
			Options: connectors,
		}
		err := survey.AskOne(prompt, &selected)
		if err != nil {
			fmt.Println("please try again")
			return
		}
		color.Green("\n✅ You selected: %s\n", selected)
	},
}

func init() {
	rootCmd.AddCommand(listConnectorsCmd)
}
