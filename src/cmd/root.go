package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "time-tracker",
	Short: "A simple time tracker",
	Long: `A simple time tracker TUI application.
See https://github.com/mrs-electronics-inc/time-tracker for more details.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
