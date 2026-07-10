package config

import (
	"encoding/json"
	"fmt"

	"github.com/Maksim-Gr/kkon/internal/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// ShowConfigCmd represents the showConfig command.
var ShowConfigCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show current configuration",
	Long:    `Print the config file path and its contents. The password is masked.`,
	Example: `  kkon config show`,
	RunE: func(_ *cobra.Command, _ []string) error {
		cfg, err := util.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if cfg.KafkaConnect.Password != "" {
			cfg.KafkaConnect.Password = "********"
		}

		if configPath, err := util.GetConfigPath(); err == nil {
			color.Cyan("Config file: %s\n", configPath)
		}

		color.Cyan("Current Configuration:")
		data, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format config: %w", err)
		}
		fmt.Printf("\n%s\n\n", string(data))
		return nil
	},
}
