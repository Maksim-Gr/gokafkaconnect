// Package cmd contains the root command and CLI entry point.
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Maksim-Gr/kkon/cmd/config"
	"github.com/Maksim-Gr/kkon/cmd/connector"
	"github.com/Maksim-Gr/kkon/cmd/task"

	"github.com/fatih/color"

	"github.com/Maksim-Gr/kkon/internal/util"

	"github.com/spf13/cobra"
)

// DryRun is the global dry-run flag shared across subcommands.
var DryRun bool

// OutputFormat is the global output format flag ("text" or "json").
var OutputFormat string

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "kkon",
	Short: "CLI to manage Kafka connector fast and easy!",
	Long: `kkon - cli tool for working with Kafka Connect.
	Manage, create, and list predefined connector in seconds!`,
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config check for config subcommands (configure, show-config, backup)
		// so the user isn't double-prompted on first-time setup.
		if cmd.Parent() != nil && cmd.Parent().Use == "config" {
			return nil
		}

		// Commands that don't talk to Kafka Connect must work without any config.
		switch cmd.Name() {
		case "version", "help", "completion":
			return nil
		}

		cfg, err := util.LoadConfig()
		if err != nil || cfg.KafkaConnect.URL == "" {
			color.Yellow("No Kafka Connect URL configured.")
			color.Cyan("Running initial configuration...\n")
			if err := config.ConfigureCmd.RunE(cmd, args); err != nil {
				return fmt.Errorf("configuration failed: %w", err)
			}
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
// Failures exit with code 1; user cancels and no-op runs exit 0.
func Execute(ctx context.Context) {
	if err := RootCmd.ExecuteContext(ctx); err != nil {
		switch {
		case errors.Is(err, util.ErrCanceled):
			color.Yellow("Canceled")
		case errors.Is(err, util.ErrNothingToDo):
			// Message already printed at the source; nothing to do is not a failure.
		default:
			fmt.Fprintln(os.Stderr, color.RedString("Error: %v", err))
			os.Exit(1)
		}
	}
}

func init() {
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.PersistentFlags().BoolVarP(&DryRun, "dry-run", "d", false, "Preview mutating commands without executing them (read-only commands ignore it)")
	RootCmd.PersistentFlags().StringVarP(&OutputFormat, "output", "o", "text", "Output format: text or json")

	// Set up command tree
	RootCmd.AddCommand(task.Cmd)
	RootCmd.AddCommand(config.Cmd)
	RootCmd.AddCommand(connector.Cmd)
}
