package cmd

import (
	"gokafkaconnect/cmd/config"
	"gokafkaconnect/cmd/connector"
	"gokafkaconnect/cmd/task"
	"os"

	"github.com/fatih/color"

	"gokafkaconnect/internal/util"

	"github.com/spf13/cobra"
)

var DryRun bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gk",
	Short: "CLI to manage Kafka connector fast and easy!",
	Long: `gk - cli tool for working  with Kafka Connect.
	Manage, create, and list predefined connector in seconds!`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Pass dryRun flag to subpackages
		config.SetDryRun(DryRun)

		color.Blue("\nChecking configuration...\n")
		cfg, err := util.LoadConfig()
		if err != nil || cfg.KafkaConnect.URL == "" {
			color.Yellow("No Kafka Connect URL configured.")
			color.Cyan("Running initial configuration...\n")
			config.ConfigureCmd.Run(cmd, args)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gokafkaconnect.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.PersistentFlags().BoolVarP(&DryRun, "dry-run", "d", false, "Dry run mode")

	// Bind global flags to subpackages
	task.BindGlobals(&DryRun)

	// Set up command tree
	RootCmd.AddCommand(task.Cmd)
	RootCmd.AddCommand(config.Cmd)
	RootCmd.AddCommand(connector.Cmd)
}
