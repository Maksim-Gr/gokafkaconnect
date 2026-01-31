package cmd

import (
	"gokafkaconnect/internal/connector"
	"gokafkaconnect/internal/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var backupDir string

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup connectors config from Kafka Connect API",
	Long:  `Backup connectors config from Kafka Connect API and save to file for future usage `,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := util.LoadConfig()
		if err != nil {
			color.Red("Failed to load config: %v\n", err)
			return
		}
		connectors, err := connector.ListConnectors(cfg.KafkaConnect.URL)
		if err != nil {
			color.Red("Failed to dump  connector config: %v\n", err)
			return
		}
		backupFile, err := connector.BackupConnectorConfig(cfg.KafkaConnect.URL, connectors, backupDir)
		if err != nil {
			color.Red("Failed to back up  connectors config: %v\n", err)
			return
		}
		color.Green("Successfully created backup: %s", backupFile)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags, which will work for this command
	// and all subcommands, e.g.:
	// dumpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dumpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	backupCmd.Flags().StringVarP(&backupDir, "dir", "d", "./backup", "Directory to save backup files")
}
