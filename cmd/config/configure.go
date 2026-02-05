package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gokafkaconnect/internal/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"

	"github.com/spf13/cobra"
)

var dryRun bool

//var configPath = filepath.Join(os.Getenv("HOME"), ".gokafkacon", "config.json")

// ConfigureCmd represents the configure command
var ConfigureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Set Kafka connect URL",
	Long:  `Configure and set Kafka Connect REST API URL.`,
	Run: func(cmd *cobra.Command, args []string) {

		if dryRun {
			color.Cyan("Dry run mode")
		} else {
			color.Cyan("\nConfiguring Kafka Connect URL...\n")
		}

		configPath, err := util.GetConfigPath()
		if err != nil {
			color.Red("Failed to determine config path: %v", err)
			os.Exit(1)
		}
		var current util.RestAPIConfig
		currentURL := ""

		if loaded, err := util.LoadConfig(); err == nil {
			current = loaded
			currentURL = loaded.KafkaConnect.URL
			color.Yellow("Current Kafka Connect URL: %s", currentURL)
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

				if s == currentURL {
					return nil
				}

				if s == "" && currentURL == "" {
					return errors.New("URL cannot be empty")
				}

				if s == "" {
					return nil
				}

				return util.ValidateURL(s)
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

		cfg := util.RestAPIConfig{
			KafkaConnect: util.KafkaConnectConfig{
				URL:      inputURL,
				Username: current.KafkaConnect.Username,
				Password: current.KafkaConnect.Password,
			},
		}

		if dryRun {
			color.Cyan("Dry run mode — config will not be saved.")
			color.Cyan("Kafka Connect URL would be: %s", inputURL)
			return
		}

		err = util.SaveConfig(cfg, configPath)
		if err != nil {
			color.Red("Failed to save config file: %s", err)
			return
		}

		color.Green("Configuration saved successfully!")
		color.Green("Kafka Connect URL: %s", inputURL)
	},
}

func SetDryRun(value bool) {
	dryRun = value
}
