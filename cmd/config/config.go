// Package config provides CLI commands for managing kkon configuration.
package config

import (
	"github.com/spf13/cobra"
)

// Cmd is the root command for configuration management.
var Cmd = &cobra.Command{
	Use:   "config",
	Short: "Manage kkon configuration",
	Long: `Configure and inspect the Kafka Connect connection settings (URL and optional
basic auth) stored in kkon's config file.`,
	Example: `  # Set the Kafka Connect URL interactively
  kkon config set

  # Print the current configuration
  kkon config show`,
}

func init() {
	Cmd.AddCommand(ConfigureCmd)
	Cmd.AddCommand(ShowConfigCmd)
}
