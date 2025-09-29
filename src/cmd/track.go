package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"time-tracker/config"
	"time-tracker/utils"
)

var trackCmd = &cobra.Command{
	Use:     "start [project title] [--id ID]",
	Short:   "Start or stop time tracking",
	Long:    `Start a new time entry with the specified project and title, or stop the current tracking session. Can be called as 'start', 'stop', or 's'.`,
	Aliases: []string{"stop", "s"},
	Args:    cobra.RangeArgs(0, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		calledAs := cmd.CalledAs()
		id, _ := cmd.Flags().GetInt("id")

		storage, err := utils.NewFileStorage(config.DataFilePath())
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		taskManager := utils.NewTaskManager(storage)

		// Determine if this is a start or stop operation
		isStop := calledAs == "stop" || (calledAs == "s" && len(args) == 0 && id == 0)
		isStart := calledAs == "start" || (calledAs == "s" && len(args) > 0) || id != 0

		if isStop {
			// Stop operation
			if len(args) > 0 {
				return fmt.Errorf("'stop' command does not accept arguments. Use 'start' to begin tracking")
			}

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

		} else if isStart {
			// Start operation
			var project, title string

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

		} else {
			return fmt.Errorf("unknown command alias: %s", calledAs)
		}
	},
}

func init() {
	trackCmd.Flags().Int("id", 0, "ID of existing time entry to resume")
	rootCmd.AddCommand(trackCmd)
}
