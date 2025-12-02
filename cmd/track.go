package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"time-tracker/config"
	"time-tracker/utils"
)

var trackCmd = &cobra.Command{
	Use:     "s [project title]",
	Short:   "Start (or stop) time tracking",
	Long:    `Start a new time entry with project and title, or stop current entry. Can be called as 'start', 'stop', or 's'.`,
	Aliases: []string{"start", "stop"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 && len(args) != 2 {
			return fmt.Errorf("accepts 0 arguments (stop) or exactly 2 arguments (project title), received %d", len(args))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		calledAs := cmd.CalledAs()

		storage, err := utils.NewFileStorage(config.DataFilePath())
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		taskManager := utils.NewTaskManager(storage)

		// Determine if this is a start or stop operation
		isStop := calledAs == "stop" || (calledAs == "s" && len(args) == 0)
		isStart := calledAs == "start" || (calledAs == "s" && len(args) > 0)

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
			if len(args) != 2 {
				return fmt.Errorf("'start' requires project and title arguments")
			}

			project := args[0]
			title := args[1]

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
	rootCmd.AddCommand(trackCmd)
}
