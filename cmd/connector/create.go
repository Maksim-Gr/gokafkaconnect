package connector

import (
	"fmt"
	"gokafkaconnect/internal/connector"
	template "gokafkaconnect/internal/connector/kafka/templates"
	"gokafkaconnect/internal/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Available connectors
var connectors = []string{
	"RabbitMQ Connector",
	"RabbitMQ  Stream Connector",
	"️Iceberg Connector",
}

// CreateCmd represents the create command
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a connector from predefined configuration  ",
	Long:  `Browse predefined connector.`,
	Run: func(cmd *cobra.Command, args []string) {
		var selected string
		color.Cyan("\n Available Kafka Connectors \n")
		prompt := &survey.Select{
			Message: "Pick a connector to work with:",
			Options: connectors,
		}
		err := survey.AskOne(prompt, &selected)
		if err != nil {
			fmt.Println("please try again")
			return
		}
		color.Green("\n You selected: %s\n", selected)
		if selected == "RabbitMQ Connector" {
			configureRedisConnector()
		}
	},
}

func configureRedisConnector() {

	color.Yellow("\n  Starting configuration for Redis Connector...\n")

	connectorConfig := template.GetRedisConnectorTemplate()

	var questions []*survey.Question
	for _, field := range template.RequiredFields() {
		var prompt survey.Prompt
		if field == "rabbitmq.password" {
			prompt = &survey.Password{Message: fmt.Sprintf("Enter %s:", field)}
		} else {
			prompt = &survey.Input{Message: fmt.Sprintf("Enter %s:", field)}
		}
		questions = append(questions, &survey.Question{
			Name:     field,
			Prompt:   prompt,
			Validate: survey.Required,
		})
	}

	answers := make(map[string]interface{})
	err := survey.Ask(questions, &answers)
	if err != nil {
		fmt.Println("Failed to get input:", err)
		return
	}

	// Update config
	for key, value := range answers {
		connectorConfig[key] = fmt.Sprintf("%v", value)
	}

	for {
		finalConfig, _ := util.ToPrettyJSON(connectorConfig)
		color.Cyan("\n Current Redis Connector Configuration:\n")
		fmt.Println(finalConfig)

		var confirmChange bool
		changePrompt := &survey.Confirm{
			Message: "Do you want to change any field?",
			Default: false,
		}
		err := survey.AskOne(changePrompt, &confirmChange)
		if err != nil {
			fmt.Println("Prompt failed:", err)
			return
		}

		if !confirmChange {
			color.Green("\n Configuration complete!\n")
			break
		}

		var fieldToChange string
		fieldPrompt := &survey.Select{
			Message: "Which field do you want to change?",
			Options: util.KeysFromMap(connectorConfig),
		}
		err = survey.AskOne(fieldPrompt, &fieldToChange)
		if err != nil {
			fmt.Println("Prompt failed:", err)
			return
		}

		var newValue string
		valuePrompt := &survey.Input{
			Message: fmt.Sprintf("Enter new value for %s:", fieldToChange),
		}
		err = survey.AskOne(valuePrompt, &newValue)
		if err != nil {
			fmt.Println("Prompt failed:", err)
			return
		}

		connectorConfig[fieldToChange] = newValue

	}
	finalConfig, _ := util.ToPrettyJSON(connectorConfig)
	color.Cyan("\nFinal Redis Connector Configuration:\n")
	fmt.Println(finalConfig)

	var submitConfirm bool
	submitPrompt := &survey.Confirm{
		Message: "Do you want to submit this connector to Kafka Connect?",
		Default: true,
	}
	err = survey.AskOne(submitPrompt, &submitConfirm)
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	if submitConfirm {
		color.Green("\n Submitting connector...\n")
		cfg, err := util.LoadConfig()

		if err != nil {
			color.Red("Failed to load config file: %v\n", err)
			return
		}
		client := connector.NewClient(cfg.KafkaConnect.URL)
		err = client.SubmitConnector(finalConfig)
		if err != nil {
			color.Red("Failed to submit connector: %v\n", err)
		} else {
			color.Green("Connector submitted successfully!\n")
		}

	} else {
		color.Yellow("\n Submission cancelled. Exiting.\n")
	}

}
