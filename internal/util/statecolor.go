package util

import "github.com/fatih/color"

// ColorState returns s formatted with the appropriate color for its Kafka Connect state.
func ColorState(s string) string {
	switch s {
	case "RUNNING":
		return color.GreenString(s)
	case "FAILED":
		return color.RedString(s)
	default:
		return color.YellowString(s)
	}
}
