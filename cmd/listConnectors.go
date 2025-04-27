/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

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
		fmt.Println("Available connectors:")
		for i, connector := range connectors {
			fmt.Printf("%d) %s\n", i+1, connector)
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Select connector number: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(connectors) {
			fmt.Println("Invalid choice")
			return
		}
		selectedConnector := connectors[choice-1]
		fmt.Printf("You chose: %s \n", selectedConnector)
	},
}

func init() {
	rootCmd.AddCommand(listConnectorsCmd)
}
