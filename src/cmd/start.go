package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"time-tracker/config"
	"time-tracker/utils"
)

var startCmd = &cobra.Command{
	Use:   "start [project title] [--id ID]",
	Short: "Start a new time entry",
	Long:  `Starts a new time entry with the specified project and title, or resumes an existing entry by ID. If another time entry is currently running, it will be automatically stopped first.`,
	Args:  cobra.RangeArgs(0, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetInt("id")
		var project, title string

		storage, err := utils.NewFileStorage(config.DataFilePath())
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		taskManager := utils.NewTaskManager(storage)

		if id != 0 {
			if len(args) > 0 {
				return fmt.Errorf("--id flag cannot be used with positional arguments")
			}
			entry, err := taskManager.GetEntry(id)
			if err != nil {
				return fmt.Errorf("failed to get entry: %w", err)
			}
			project = entry.Project
			title = entry.Title
		} else {
			if len(args) != 2 {
				return fmt.Errorf("project and title required when not using --id")
			}
			project = args[0]
			title = args[1]
		}

		entry, err := taskManager.StartEntry(project, title)
		if err != nil {
			return fmt.Errorf("failed to start time entry: %w", err)
		}

		fmt.Printf("Started tracking time for \"%s\" in project \"%s\"\n", entry.Title, entry.Project)
		return nil
	},
}

func init() {
	startCmd.Flags().Int("id", 0, "ID of existing time entry to resume")
	rootCmd.AddCommand(startCmd)
}
