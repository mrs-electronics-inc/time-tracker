package cmd

import (
	"fmt"
	"os"
	"sort"
	"time"
	"time-tracker/config"
	"time-tracker/utils"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func formatTimeHHMM(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Display time tracking statistics",
	Long:  `Display various statistics about tracked time, including daily totals, weekly totals, and project breakdowns.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		weeklyFlag, err := cmd.Flags().GetBool("weekly")
		if err != nil {
			return fmt.Errorf("failed to parse weekly flag: %w", err)
		}

		storage, err := utils.NewFileStorage(config.DataFilePath())
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}

		entries, err := storage.Load()
		if err != nil {
			return fmt.Errorf("failed to load entries: %w", err)
		}

		// Time-based stats
		if weeklyFlag {
			weeklyTotals := utils.CalculateWeeklyTotals(entries)
			if len(weeklyTotals) == 0 {
				fmt.Println("No data available")
				return nil
			}
			// Collect all projects
			projectSet := make(map[string]bool)
			for _, total := range weeklyTotals {
				for project := range total.Projects {
					projectSet[project] = true
				}
			}
			var projects []string
			for p := range projectSet {
				projects = append(projects, p)
			}
			sort.Strings(projects)

			headers := []string{"Week Starting", "Total"}
			headers = append(headers, projects...)

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader(headers)
			table.SetBorder(true)
			table.SetRowLine(true)
			table.SetAutoWrapText(false)
			for _, total := range weeklyTotals {
				row := []string{
					total.WeekStart.Format("2006-01-02"),
					formatTimeHHMM(total.Total),
				}
				for _, p := range projects {
					dur := total.Projects[p]
					if dur == 0 {
						row = append(row, "")
					} else {
						row = append(row, formatTimeHHMM(dur))
					}
				}
				table.Append(row)
			}
			table.Render()
		} else {
			dailyTotals := utils.CalculateDailyTotals(entries)
			if len(dailyTotals) == 0 {
				fmt.Println("No data available")
				return nil
			}
			// Collect all projects
			projectSet := make(map[string]bool)
			for _, total := range dailyTotals {
				for project := range total.Projects {
					projectSet[project] = true
				}
			}
			var projects []string
			for p := range projectSet {
				projects = append(projects, p)
			}
			sort.Strings(projects)

			headers := []string{"Date", "Total"}
			headers = append(headers, projects...)

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader(headers)
			table.SetBorder(true)
			table.SetRowLine(true)
			table.SetAutoWrapText(false)
			for _, total := range dailyTotals {
				row := []string{
					total.Date.Format("2006-01-02"),
					formatTimeHHMM(total.Total),
				}
				for _, p := range projects {
					dur := total.Projects[p]
					if dur == 0 {
						row = append(row, "")
					} else {
						row = append(row, formatTimeHHMM(dur))
					}
				}
				table.Append(row)
			}
			table.Render()
		}

		return nil
	},
}

func init() {
	statsCmd.Flags().BoolP("weekly", "w", false, "Show weekly totals for the past month")

	rootCmd.AddCommand(statsCmd)
}
