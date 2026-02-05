package cmd

import (
	"gokafkaconnect/internal/connector"
	"gokafkaconnect/internal/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var connectorName string

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete connector",
	Long:  `Delete connector from Kafka Connect API`,
	Run: func(cmd *cobra.Command, args []string) {
		if connectorName == "" {
			color.Red("error:  provide connector name")
			return
		}

		color.Yellow("Deleting connector: %s", connectorName)

		cfg, err := util.LoadConfig()
		if err != nil {
			color.Red("Failed to load config file: %v\n", err)
			return
		}
		client := connector.NewClient(cfg.KafkaConnect.URL)
		err = client.DeleteConnector(connectorName)
		if err != nil {
			color.Red("Failed to delete connector: %v\n", err)
		} else {
			color.Green("Connector deleted successfully!")
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	deleteCmd.Flags().StringVarP(&connectorName, "connector", "c", "", "Connector name to delete")
	deleteCmd.MarkFlagRequired("connector")
}
