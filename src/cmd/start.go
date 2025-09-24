package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"time-tracker/config"
	"time-tracker/utils"
)

var startCmd = &cobra.Command{
	Use:   "start <project> <title>",
	Short: "Start a new time entry",
	Long:  `Starts a new time entry with the specified project and title. If another time entry is currently running, it will be automatically stopped first.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		project := args[0]
		title := args[1]

		storage, err := utils.NewFileStorage(config.DataFilePath())
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		taskManager := utils.NewTaskManager(storage)

		entry, err := taskManager.StartEntry(project, title)
		if err != nil {
			return fmt.Errorf("failed to start time entry: %w", err)
		}

		fmt.Printf("Started tracking time for \"%s\" in project \"%s\"\n", entry.Title, entry.Project)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
