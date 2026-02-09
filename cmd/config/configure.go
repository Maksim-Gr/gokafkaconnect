package config

import (
	"errors"
	"os"
	"strings"

	"gokafkaconnect/internal/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var dryRun bool

// ConfigureCmd represents the configure command
var ConfigureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure Kafka Connect REST API",
	Long:  `Configure Kafka Connect REST API URL and authentication.`,
	Run: func(cmd *cobra.Command, args []string) {

		if dryRun {
			color.Cyan("Dry run mode")
		} else {
			color.Cyan("\nConfiguring Kafka Connect...\n")
		}

		configPath, err := util.GetConfigPath()
		if err != nil {
			color.Red("Failed to determine config path: %v", err)
			os.Exit(1)
		}

		currentURL := ""
		currentUser := ""
		currentPass := ""

		if loaded, err := util.LoadConfig(); err == nil {
			currentURL = loaded.KafkaConnect.URL
			currentUser = loaded.KafkaConnect.Username
			currentPass = loaded.KafkaConnect.Password
			color.Yellow("Current Kafka Connect URL: %s", currentURL)
		}

		// --- URL prompt ---
		var inputURL string
		urlPrompt := &survey.Input{
			Message: "Kafka Connect URL:",
			Help:    "Enter the URL of your Kafka Connect REST API (e.g. http://localhost:8083)",
			Default: currentURL,
		}

		err = survey.AskOne(urlPrompt, &inputURL, survey.WithValidator(
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
			color.Red("Failed: %v", err)
			os.Exit(1)
		}

		if !strings.HasPrefix(inputURL, "http://") && !strings.HasPrefix(inputURL, "https://") {
			color.Yellow("No scheme specified — assuming http://")
			inputURL = "http://" + inputURL
		}

		var inputUser string
		userPrompt := &survey.Input{
			Message: "Kafka Connect username (leave empty for no auth):",
			Default: currentUser,
		}

		if err := survey.AskOne(userPrompt, &inputUser); err != nil {
			color.Red("Failed to read username: %v", err)
			os.Exit(1)
		}

		inputPass := currentPass
		if inputUser != "" {
			passPrompt := &survey.Password{
				Message: "Kafka Connect password:",
				Help:    "Password will be stored in your local config file",
			}
			if err := survey.AskOne(passPrompt, &inputPass); err != nil {
				color.Red("Failed to read password: %v", err)
				os.Exit(1)
			}
		} else {
			inputPass = ""
		}

		cfg := util.RestAPIConfig{
			KafkaConnect: util.KafkaConnectConfig{
				URL:      inputURL,
				Username: inputUser,
				Password: inputPass,
			},
		}

		if dryRun {
			color.Cyan("Dry run mode — config will not be saved.")
			color.Cyan("Kafka Connect URL: %s", inputURL)
			if inputUser != "" {
				color.Cyan("Authentication: enabled (username: %s)", inputUser)
			} else {
				color.Cyan("Authentication: disabled")
			}
			return
		}

		if err := util.SaveConfig(cfg, configPath); err != nil {
			color.Red("Failed to save config file: %v", err)
			os.Exit(1)
		}

		color.Green("Configuration saved successfully!")
		color.Green("Kafka Connect URL: %s", inputURL)
		if inputUser != "" {
			color.Green("Authentication enabled for user: %s", inputUser)
		} else {
			color.Green("Authentication disabled")
		}
	},
}

func SetDryRun(value bool) {
	dryRun = value
}
