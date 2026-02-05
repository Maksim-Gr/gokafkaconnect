package connector

import (
	"gokafkaconnect/internal/connector"
	"gokafkaconnect/internal/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var connectorName string

// DeleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
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
	DeleteCmd.Flags().StringVarP(&connectorName, "connector", "c", "", "Connector name to delete")
	DeleteCmd.MarkFlagRequired("connector")
}
