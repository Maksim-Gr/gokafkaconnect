package connector

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	template "github.com/Maksim-Gr/kkon/internal/connector/kafka/templates"
	"github.com/Maksim-Gr/kkon/internal/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Available connectors.
var connectors = []string{
	"RabbitMQ Connector",
	"Debezium PostgreSQL CDC",
	"JDBC Source Connector",
	"JDBC Sink Connector",
	"S3 Sink Connector",
}

var connectorJSONPath string

// CreateCmd represents the create command.
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a connector from predefined configuration",
	Long:  `Browse predefined connector.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		if connectorJSONPath != "" {
			return submitConnectorFromFile(cmd, connectorJSONPath)
		}

		if util.IsJSONOutput(cmd) {
			return errors.New("--output json requires --file (the create wizard is interactive)")
		}

		var selected string
		color.Cyan("\n Available Kafka Connectors \n")
		prompt := &survey.Select{
			Message: "Pick a connector to work with:",
			Options: connectors,
		}
		if err := survey.AskOne(prompt, &selected); err != nil {
			return util.ErrCanceled
		}
		color.Green("\n You selected: %s\n", selected)
		switch selected {
		case "RabbitMQ Connector":
			return configureConnector(cmd, "RabbitMQ Connector", template.GetRabbitMQConnectorTemplate(), template.RabbitMQRequiredFields(), "rabbitmq.password")
		case "Debezium PostgreSQL CDC":
			return configureConnector(cmd, "Debezium PostgreSQL CDC", template.GetDebeziumPostgresConnectorTemplate(), template.DebeziumPostgresRequiredFields(), "database.password")
		case "JDBC Source Connector":
			return configureConnector(cmd, "JDBC Source Connector", template.GetJDBCSourceConnectorTemplate(), template.JDBCSourceRequiredFields(), "connection.password")
		case "JDBC Sink Connector":
			return configureConnector(cmd, "JDBC Sink Connector", template.GetJDBCSinkConnectorTemplate(), template.JDBCSinkRequiredFields(), "connection.password")
		case "S3 Sink Connector":
			return configureConnector(cmd, "S3 Sink Connector", template.GetS3SinkConnectorTemplate(), template.S3SinkRequiredFields(), "")
		}
		return nil
	},
}

func init() {
	CreateCmd.Flags().StringVarP(&connectorJSONPath, "file", "f", "", "Path to connector JSON config file")
}

func submitConnectorFromFile(cmd *cobra.Command, path string) error {
	ctx := cmd.Context()
	jsonMode := util.IsJSONOutput(cmd)

	b, err := os.ReadFile(path) //nolint:gosec
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	var js json.RawMessage
	if err := json.Unmarshal(b, &js); err != nil {
		return fmt.Errorf("invalid JSON in %s: %w", path, err)
	}

	client, err := util.NewKafkaConnectClient()
	if err != nil {
		return err
	}

	if configMap := configMapFromFile(b); configMap != nil {
		if jsonMode {
			if err := validateConfigStrict(ctx, client, configMap); err != nil {
				return err
			}
		} else if !validateConfigOrConfirm(ctx, client, configMap) {
			color.Yellow("Submission cancelled.\n")
			return util.ErrNothingToDo
		}
	}

	if util.IsDryRun(cmd) {
		color.Yellow("[dry-run] Would submit connector from file %s\n", path)
		return nil
	}

	if !jsonMode {
		color.Green("\n Submitting connector from file: %s ...\n", path)
	}
	if _, err := client.SubmitConnector(ctx, string(b)); err != nil {
		return fmt.Errorf("failed to submit connector: %w", err)
	}

	if jsonMode {
		var meta struct {
			Name string `json:"name"`
		}
		_ = json.Unmarshal(b, &meta)
		printActionJSON(meta.Name, "create")
		return nil
	}
	color.Green("Connector submitted successfully!\n")
	return nil
}

func configureConnector(cmd *cobra.Command, name string, connectorConfig map[string]string, required []string, passwordField string) error {
	ctx := cmd.Context()
	color.Yellow("\n  Starting configuration for %s...\n", name)

	questions := make([]*survey.Question, 0, len(required))
	for _, field := range required {
		var prompt survey.Prompt
		if field == passwordField {
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
	if err := survey.Ask(questions, &answers); err != nil {
		if util.IsSurveyInterrupt(err) {
			return util.ErrCanceled
		}
		return fmt.Errorf("failed to get input: %w", err)
	}

	for key, value := range answers {
		connectorConfig[key] = fmt.Sprintf("%v", value)
	}

	for {
		finalConfig, err := util.ToPrettyJSON(connectorConfig)
		if err != nil {
			return fmt.Errorf("failed to format config: %w", err)
		}
		color.Cyan("\n Current %s Configuration:\n", name)
		fmt.Println(finalConfig)

		var confirmChange bool
		changePrompt := &survey.Confirm{
			Message: "Do you want to change any field?",
			Default: false,
		}
		if err := survey.AskOne(changePrompt, &confirmChange); err != nil {
			if util.IsSurveyInterrupt(err) {
				return util.ErrCanceled
			}
			return fmt.Errorf("prompt failed: %w", err)
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
		if err := survey.AskOne(fieldPrompt, &fieldToChange); err != nil {
			if util.IsSurveyInterrupt(err) {
				return util.ErrCanceled
			}
			return fmt.Errorf("prompt failed: %w", err)
		}

		var newValue string
		valuePrompt := &survey.Input{
			Message: fmt.Sprintf("Enter new value for %s:", fieldToChange),
		}
		if err := survey.AskOne(valuePrompt, &newValue); err != nil {
			if util.IsSurveyInterrupt(err) {
				return util.ErrCanceled
			}
			return fmt.Errorf("prompt failed: %w", err)
		}

		connectorConfig[fieldToChange] = newValue
	}
	finalConfig, err := util.ToPrettyJSON(connectorConfig)
	if err != nil {
		return fmt.Errorf("failed to format config: %w", err)
	}
	color.Cyan("\nFinal %s Configuration:\n", name)
	fmt.Println(finalConfig)

	client, err := util.NewKafkaConnectClient()
	if err != nil {
		return err
	}

	if util.IsDryRun(cmd) {
		if !validateConfigOrConfirm(ctx, client, connectorConfig) {
			color.Yellow("\n Submission cancelled. Exiting.\n")
			return util.ErrNothingToDo
		}
		color.Yellow("[dry-run] Would submit connector %s\n", name)
		return nil
	}

	var submitConfirm bool
	submitPrompt := &survey.Confirm{
		Message: "Do you want to submit this connector to Kafka Connect?",
		Default: true,
	}
	if err := survey.AskOne(submitPrompt, &submitConfirm); err != nil {
		if util.IsSurveyInterrupt(err) {
			return util.ErrCanceled
		}
		return fmt.Errorf("prompt failed: %w", err)
	}

	if !submitConfirm {
		color.Yellow("\n Submission cancelled. Exiting.\n")
		return util.ErrNothingToDo
	}

	color.Green("\n Submitting connector...\n")
	if !validateConfigOrConfirm(ctx, client, connectorConfig) {
		color.Yellow("\n Submission cancelled. Exiting.\n")
		return util.ErrNothingToDo
	}

	if _, err := client.SubmitConnector(ctx, finalConfig); err != nil {
		return fmt.Errorf("failed to submit connector: %w", err)
	}
	color.Green("Connector submitted successfully!\n")
	return nil
}
