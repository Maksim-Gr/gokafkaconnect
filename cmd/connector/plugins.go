package connector

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Maksim-Gr/kkon/internal/connector"
	"github.com/Maksim-Gr/kkon/internal/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var pluginType string

// PluginsCmd lists the connector plugins installed on the cluster.
var PluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "List connector plugins installed on the cluster",
	Long: `List the connector plugin classes available on the Kafka Connect cluster.
Useful for finding a class to pass to 'kkon connector create --plugin'.`,
	Example: `  # List all installed plugins
  kkon connector plugins

  # Show only source plugins
  kkon connector plugins --type source

  # Machine-readable output
  kkon connector plugins --output json`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := util.NewKafkaConnectClient()
		if err != nil {
			return err
		}

		stop := util.StartSpinner("Fetching connector plugins...")
		plugins, err := client.ListConnectorPlugins(cmd.Context())
		stop()
		if err != nil {
			return fmt.Errorf("failed to list connector plugins: %w", err)
		}

		if pluginType != "" {
			filtered := make([]connector.Plugin, 0, len(plugins))
			for _, p := range plugins {
				if strings.EqualFold(p.Type, pluginType) {
					filtered = append(filtered, p)
				}
			}
			plugins = filtered
		}

		if util.IsJSONOutput(cmd) {
			b, _ := json.MarshalIndent(plugins, "", "  ")
			fmt.Println(string(b))
			return nil
		}

		if len(plugins) == 0 {
			color.Yellow("No connector plugins found\n")
			return nil
		}

		maxLen := 0
		for _, p := range plugins {
			if len(p.Class) > maxLen {
				maxLen = len(p.Class)
			}
		}

		color.Cyan("Connector plugins:")
		for _, p := range plugins {
			version := p.Version
			if version == "" {
				version = "unknown"
			}
			fmt.Printf("  %-*s  %-6s  %s\n", maxLen, p.Class, p.Type, version)
		}
		return nil
	},
}

func init() {
	PluginsCmd.Flags().StringVar(&pluginType, "type", "", "Only show plugins of this type (source or sink)")
}
