package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"os"
	"path/filepath"

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

//var configPath = filepath.Join(os.Getenv("HOME"), ".gokafkacon", "config.json")

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Set Kafka connect URL ",
	Long:  `Configure and set Kafka Connect REST API URL`,
	Run: func(cmd *cobra.Command, args []string) {
		color.Cyan("\n Configuring Kafka Connect URL...\n")
		configPath, err := getConfigPath()
		if err != nil {
			color.Red("Failed to determine config path: %v", err)
			return
		}

		var url string
		prompt := &survey.Input{
			Message: "Kafka Connect URL:",
		}
		err = survey.AskOne(prompt, &url, survey.WithValidator(survey.Required))
		if err != nil {
			fmt.Println("Failed: ", err)
			return
		}

		cfg := RestAPIConfig{
			KafkaConnectURL: url,
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
