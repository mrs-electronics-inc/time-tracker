package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"time-tracker/config"
	"time-tracker/models"
	"time-tracker/utils"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List time entries",
	Long: `List time entries from data.json in chronological order (oldest first).
By default, only the entries from the current day will be shown. Use --all to view all entries`,
	Aliases: []string{"l", "ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		displayAll, err := cmd.Flags().GetBool("all")
		if err != nil {
			return fmt.Errorf("failed to parse all flag")
		}

		storage, err := utils.NewFileStorage(config.DataFilePath())
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		taskManager := utils.NewTaskManager(storage)

		allEntries, err := taskManager.ListEntries()
		if err != nil {
			return fmt.Errorf("failed to load time entries: %w", err)
		}

		displayEntries := []models.TimeEntry{}
		if displayAll {
			// Filter out blank entries
			for _, entry := range allEntries {
				if !entry.IsBlank() {
					displayEntries = append(displayEntries, entry)
				}
			}
		} else {
			now := time.Now()
			year, month, day := now.Date()
			startOfToday := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
			for _, entry := range allEntries {
				// Filter out blank entries
				if entry.IsBlank() {
					continue
				}
				if entry.IsRunning() || entry.End.After(startOfToday) {
					displayEntries = append(displayEntries, entry)
				}
			}
		}

		if len(displayEntries) == 0 {
			fmt.Println("No time entries found")
			return nil
		}

		displayEntriesTable(displayEntries)
		return nil
	},
}

func displayEntriesTable(entries []models.TimeEntry) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Start", "End", "Project", "Title", "Duration"})
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetAutoWrapText(false)

	for _, entry := range entries {
		startTime := entry.Start.Format("2006-01-02 15:04")
		endTime := "\033[32mrunning\033[0m"
		if entry.End != nil {
			endTime = entry.End.Format("2006-01-02 15:04")
		}

		duration := utils.FormatDuration(entry.Duration())

		row := []string{
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

func init() {
	listCmd.Flags().BoolP("all", "a", false, "display all time entries")
	rootCmd.AddCommand(listCmd)
}
