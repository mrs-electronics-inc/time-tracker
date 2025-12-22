package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"time-tracker/config"
	"time-tracker/utils"
)

var currentCmd = &cobra.Command{
	Use:     "current",
	Short:   "Show the currently running task",
	Long:    `Display the currently running task without entering the TUI. Can be called as 'current', 'curr', or 'c'.`,
	Aliases: []string{"curr", "c"},
	RunE: func(cmd *cobra.Command, args []string) error {
		storage, err := utils.NewFileStorage(config.DataFilePath())
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		taskManager := utils.NewTaskManager(storage)

		entries, err := taskManager.ListEntries()
		if err != nil {
			return fmt.Errorf("failed to load entries: %w", err)
		}

		// Find the last entry
		if len(entries) == 0 {
			return nil
		}

		lastEntry := entries[len(entries)-1]

		// Check if it's running and not blank
		if !lastEntry.IsRunning() || lastEntry.IsBlank() {
			return nil
		}

		// Format duration
		duration := lastEntry.Duration()
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60
		durationStr := fmt.Sprintf("%dh %dm", hours, minutes)
		if hours == 0 {
			durationStr = fmt.Sprintf("%dm", minutes)
		}

		fmt.Printf("%s %s, duration %s\n", lastEntry.Project, lastEntry.Title, durationStr)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
}
