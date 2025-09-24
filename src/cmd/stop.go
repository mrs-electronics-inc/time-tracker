package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"time-tracker/config"
	"time-tracker/utils"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the currently running time entry",
	Long:  `Stops the currently running time entry by setting its end timestamp to the current time.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		storage, err := utils.NewFileStorage(config.DataFilePath())
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		taskManager := utils.NewTaskManager(storage)

		entry, err := taskManager.StopEntry()
		if err != nil {
			return fmt.Errorf("failed to stop time entry: %w", err)
		}

		duration := entry.Duration()
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60
		durationStr := fmt.Sprintf("%dh %dm", hours, minutes)
		if hours == 0 {
			durationStr = fmt.Sprintf("%dm", minutes)
		}

		fmt.Printf("Stopped tracking time for \"%s\" in project \"%s\" (duration: %s)\n", entry.Title, entry.Project, durationStr)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
