package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"

	"github.com/spf13/cobra"
)

type RestAPIConfig struct {
	KafkaConnectURL string `json:"kafka_connect_url"`
}

func getExecutablePath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}

func getConfigPath() (string, error) {
	exe, err := getExecutablePath()
	if err != nil {
		return "", err
	}
	return filepath.Join(exe, "config.json"), nil
}

func validateURL(input string) error {
	parsed, err := url.ParseRequestURI(input)
	if err != nil {
		return errors.New("invalid URL")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("url should starts with http:// or https:// ")
	}
	if parsed.Host == "" {
		return errors.New("url should contains host")
	}
	return nil
}

//var configPath = filepath.Join(os.Getenv("HOME"), ".gokafkacon", "config.json")

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Set Kafka connect URL ",
	Long:  `Configure and set Kafka Connect REST API URL`,
	Run: func(cmd *cobra.Command, args []string) {

		if dryRun {
			color.Cyan("Dry run mode")
		} else {
			color.Cyan("\n Configuring Kafka Connect URL...\n")
		}

		configPath, err := getConfigPath()
		if err != nil {
			color.Red("Failed to determine config path: %v", err)
			os.Exit(1)
		}
		if currentConfig, err := LoadConfig(); err == nil {
			color.Yellow("Current URL %s", currentConfig.KafkaConnectURL)
		}
		var url string
		prompt := &survey.Input{
			Message: "Kafka Connect URL:",
			Help:    "Enter the URL of your Kafka Connect REST API",
		}
		err = survey.AskOne(prompt, &url,
			survey.WithValidator(survey.ComposeValidators(
				survey.Required, func(ans interface{}) error {
					str := ans.(string)
					return validateURL(str)
				},
			)),
		)
		if err != nil {
			fmt.Println("Failed: ", err)
			os.Exit(1)
		}

		cfg := RestAPIConfig{
			KafkaConnectURL: url,
		}

		if dryRun {
			color.Cyan("Dry run mode")
		}

		err = saveConfig(cfg, configPath)
		if err != nil {
			color.Red("Failed to save config file: %s", err)
			return
		}
	},
}

// Save config to file
func saveConfig(cfg RestAPIConfig, configPath string) error {
	// Create a directory if not exists
	err := os.MkdirAll(filepath.Dir(configPath), os.ModePerm)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// LoadConfig Load config
func LoadConfig() (RestAPIConfig, error) {
	var cfg RestAPIConfig
	configPath, err := getConfigPath()
	if err != nil {
		return cfg, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}

func init() {
	rootCmd.AddCommand(configureCmd)

}
