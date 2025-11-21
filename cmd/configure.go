package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

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
	if input == "" {
		return errors.New("URL cannot be empty")
	}

	testURL := input
	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		testURL = "http://" + input
	}

	parsed, err := url.ParseRequestURI(testURL)
	if err != nil {
		return errors.New("invalid URL format")
	}

	if parsed.Host == "" {
		return errors.New("URL must contain a host (e.g. localhost:8083 or example.com)")
	}

	host := parsed.Hostname()
	if !strings.Contains(host, ".") && !strings.Contains(host, ":") && host != "localhost" {
		return fmt.Errorf("invalid host: %s (must include '.' or ':' or be 'localhost')", host)
	}

	return nil
}

//var configPath = filepath.Join(os.Getenv("HOME"), ".gokafkacon", "config.json")

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Set Kafka connect URL",
	Long:  `Configure and set Kafka Connect REST API URL.`,
	Run: func(cmd *cobra.Command, args []string) {

		if dryRun {
			color.Cyan("Dry run mode")
		} else {
			color.Cyan("\nConfiguring Kafka Connect URL...\n")
		}

		configPath, err := getConfigPath()
		if err != nil {
			color.Red("Failed to determine config path: %v", err)
			os.Exit(1)
		}
		var currentURL string
		if currentConfig, err := LoadConfig(); err == nil {
			currentURL = currentConfig.KafkaConnectURL
			color.Yellow("Current URL: %s", currentURL)
		}

		var inputURL string
		prompt := &survey.Input{
			Message: "Kafka Connect URL:",
			Help:    "Enter the URL of your Kafka Connect REST API (e.g. http://localhost:8083)",
			Default: currentURL,
		}

		err = survey.AskOne(prompt, &inputURL, survey.WithValidator(
			func(ans interface{}) error {
				s := ans.(string)

				// If unchanged → OK
				if s == currentURL {
					return nil
				}

				// If empty and no current URL → reject
				if s == "" && currentURL == "" {
					return errors.New("URL cannot be empty")
				}

				// If empty but current URL exists → OK
				if s == "" {
					return nil
				}

				// Validate only new value
				return validateURL(s)
			},
		))

		if err != nil {
			fmt.Println("Failed:", err)
			os.Exit(1)
		}

		if !strings.HasPrefix(inputURL, "http://") && !strings.HasPrefix(inputURL, "https://") {
			color.Yellow("No scheme specified — assuming http://")
			inputURL = "http://" + inputURL
		}

		cfg := RestAPIConfig{
			KafkaConnectURL: inputURL,
		}

		if dryRun {
			color.Cyan("Dry run mode — config will not be saved.")
			color.Cyan("Kafka Connect URL would be: %s", inputURL)
			return
		}

		err = saveConfig(cfg, configPath)
		if err != nil {
			color.Red("Failed to save config file: %s", err)
			return
		}

		color.Green("Configuration saved successfully!")
		color.Green("Kafka Connect URL: %s", inputURL)
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
