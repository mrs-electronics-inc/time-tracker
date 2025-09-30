package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"time-tracker/config"
	"time-tracker/models"
	"time-tracker/utils"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all time entries",
	Long:    `List all time entries from data.json in chronological order (newest first).`,
	Aliases: []string{"l", "ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		storage, err := utils.NewFileStorage(config.DataFilePath())
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		taskManager := utils.NewTaskManager(storage)

		entries, err := taskManager.ListEntries()
		if err != nil {
			return fmt.Errorf("failed to load time entries: %w", err)
		}

		if len(entries) == 0 {
			fmt.Println("No time entries found")
			return nil
		}

		displayEntriesTable(entries)
		return nil
	},
}

func displayEntriesTable(entries []models.TimeEntry) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Start", "End", "Project", "Title", "Duration"})
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetAutoWrapText(false)

	for _, entry := range entries {
		startTime := entry.Start.Format("2006-01-02 15:04")
		endTime := "\033[32mrunning\033[0m"
		if entry.End != nil {
			endTime = entry.End.Format("2006-01-02 15:04")
		}

		duration := formatDuration(entry.Duration())

		row := []string{
			strconv.Itoa(entry.ID),
			startTime,
			endTime,
			entry.Project,
			entry.Title,
			duration,
		}

		table.Append(row)
	}

	table.Render()
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

func init() {
	rootCmd.AddCommand(listCmd)
}
