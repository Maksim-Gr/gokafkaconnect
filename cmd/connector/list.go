package connector

import (
	"encoding/json"
	"fmt"
	"gokafkaconnect/internal/connector"
	"gokafkaconnect/internal/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var listConfigName string

// ListCmd represent command for retrieving connectors from API.
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List running connector",
	Long:  `List current running connector`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := util.LoadConfig()
		if err != nil {
			color.Red("Failed to load config: %v\n", err)
			return
		}
		client := connector.NewClient(cfg.KafkaConnect.URL)
		if cfg.KafkaConnect.Username != "" {
			client.SetBasicAuth(cfg.KafkaConnect.Username, cfg.KafkaConnect.Password)
		}
		connectors, err := client.ListConnectors(cmd.Context())
		if err != nil {
			color.Red("Failed to list connector: %v\n", err)
			return
		}

		if len(connectors) == 0 {
			color.Yellow("No connectors found\n")
			return
		}
		color.Cyan("Connectors:")
		for _, connector := range connectors {
			fmt.Printf("\t%s\n", connector)
		}

		selected := listConfigName
		if selected == "" {
			prompt := &survey.Select{
				Message: "Show connector config:",
				Options: connectors,
			}
			if err := survey.AskOne(prompt, &selected); err != nil {
				color.Yellow("Canceled\n")
				return
			}
		}

		config, err := client.GetConnectorConfig(cmd.Context(), selected)
		if err != nil {
			color.Red("Failed to get connector config: %v\n", err)
			return
		}
		color.Green("config for %s connector:\n", selected)
		var raw map[string]interface{}
		if err := json.Unmarshal([]byte(config), &raw); err != nil {
			fmt.Println(config)
			return
		}
		pretty, err := util.ToPrettyJSON(raw)
		if err != nil {
			fmt.Println(config)
			return
		}
		fmt.Println(pretty)
	},
}

func init() {
	ListCmd.Flags().StringVarP(&listConfigName, "config", "c", "", "Print config for the named connector (skips interactive prompt)")
}
