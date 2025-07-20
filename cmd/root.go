package cmd

import (
	"github.com/fatih/color"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gokafkaconnect",
	Short: "CLI to manage Kafka connector fast and easy!",
	Long: `gokafkaconnect - cli tool for working  with Kafka Connect.
	Manage, create, and list predefined connector in seconds!`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		color.Blue("\nChecking configuration...\n")
		cfg, err := LoadConfig()
		if err != nil || cfg.KafkaConnectURL == "" {
			color.Yellow("No Kafka Connect URL configured.")
			color.Cyan("Running initial configuration...\n")
			configureCmd.Run(cmd, args)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gokafkaconnect.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
