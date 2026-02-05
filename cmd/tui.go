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
	RootCmd.AddCommand(tuiCmd)
}
