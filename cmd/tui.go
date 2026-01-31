/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"gokafkaconnect/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// tuiCmd represents the tui command
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "start tui",
	Long:  `Start TUI for managing connectors`,
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(tui.New("./backup"))
		if err := p.Start(); err != nil {
			color.Red("Failed to start TUI: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tuiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tuiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
