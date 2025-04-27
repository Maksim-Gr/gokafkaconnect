/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// Available connectors
// TODO replace with config files
var connectors = []string{"RabbitMQ", "RabbitMQStream"}

// listConnectorsCmd represents the list-connectors command
var listConnectorsCmd = &cobra.Command{
	Use:   "list-connectors",
	Short: "List available connectors",
	Long:  `List connectors provide a list of  available connectors configuration files.`,
	Run: func(cmd *cobra.Command, args []string) {
		var selected string
		prompt := &survey.Select{
			Message: "Choose a connector:",
			Options: connectors,
		}
		err := survey.AskOne(prompt, &selected)
		if err != nil {
			fmt.Println("please try again")
			return
		}
		fmt.Println("Selected ", selected)
	},
}

func init() {
	rootCmd.AddCommand(listConnectorsCmd)
}
