package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"time-tracker/config"
	"time-tracker/utils"
)

var trackCmd = &cobra.Command{
	Use:     "s [project title | ID]",
	Short:   "Start (or stop) time tracking",
	Long:    `Start a new time entry with project and title, resume by ID, or stop current entry. Can be called as 'start', 'stop', or 's'.`,
	Aliases: []string{"start", "stop"},
	Args:    cobra.RangeArgs(0, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		calledAs := cmd.CalledAs()

		storage, err := utils.NewFileStorage(config.DataFilePath())
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		taskManager := utils.NewTaskManager(storage)

		// Determine if this is a start or stop operation
		isStop := calledAs == "stop" || (calledAs == "s" && len(args) == 0)
		isStart := calledAs == "start" || calledAs == "s"

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

			if len(args) == 1 {
				// Single argument: treat as ID to resume
				id, err := strconv.Atoi(args[0])
				if err != nil {
					return fmt.Errorf("invalid ID: %s", args[0])
				}
				entry, err := taskManager.GetEntry(id)
				if err != nil {
					return fmt.Errorf("failed to get entry: %w", err)
				}
				project = entry.Project
				title = entry.Title
			} else if len(args) == 2 {
				project = args[0]
				title = args[1]
			} else {
				return fmt.Errorf("provide project and title, or single ID to resume")
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
	rootCmd.AddCommand(trackCmd)
}
