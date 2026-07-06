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

// collectProjects gathers all non-empty project names from stats totals and returns them sorted
func collectProjects(projectMaps []map[string]time.Duration) []string {
	projectSet := make(map[string]bool)
	for _, projects := range projectMaps {
		for project := range projects {
			if project != "" {
				projectSet[project] = true
			}
		}
	}
	var result []string
	for p := range projectSet {
		result = append(result, p)
	}
	sort.Strings(result)
	return result
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

		monthlyFlag, err := cmd.Flags().GetBool("monthly")
		if err != nil {
			return fmt.Errorf("failed to parse monthly flag: %w", err)
		}

		rows, err := cmd.Flags().GetInt("rows")
		if err != nil {
			return fmt.Errorf("failed to parse rows flag: %w", err)
		}

		// Set default row counts based on view mode if user didn't specify rows
		if !cmd.Flags().Changed("rows") {
			switch {
			case weeklyFlag:
				rows = 4
			case monthlyFlag:
				rows = 6
			}
		}

		// Validate rows value
		if rows <= 0 {
			return fmt.Errorf("rows must be a positive integer")
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
		if monthlyFlag {
			monthlyTotals := utils.CalculateMonthlyTotals(entries, rows)
			if len(monthlyTotals) == 0 {
				fmt.Println("No data available")
				return nil
			}
			projectMaps := make([]map[string]time.Duration, len(monthlyTotals))
			for i, total := range monthlyTotals {
				projectMaps[i] = total.Projects
			}
			projects := collectProjects(projectMaps)

			headers := []string{"Month", "Total"}
			headers = append(headers, projects...)

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader(headers)
			table.SetBorder(true)
			table.SetRowLine(true)
			table.SetAutoWrapText(false)
			for _, total := range monthlyTotals {
				row := []string{
					total.Month.Format("2006-01"),
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
		} else if weeklyFlag {
			weeklyTotals := utils.CalculateWeeklyTotals(entries, rows)
			if len(weeklyTotals) == 0 {
				fmt.Println("No data available")
				return nil
			}
			// Collect all projects (excluding empty project name)
			projectMaps := make([]map[string]time.Duration, len(weeklyTotals))
			for i, total := range weeklyTotals {
				projectMaps[i] = total.Projects
			}
			projects := collectProjects(projectMaps)

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
			dailyTotals := utils.CalculateDailyTotals(entries, rows)
			if len(dailyTotals) == 0 {
				fmt.Println("No data available")
				return nil
			}
			// Collect all projects (excluding empty project name)
			projectMaps := make([]map[string]time.Duration, len(dailyTotals))
			for i, total := range dailyTotals {
				projectMaps[i] = total.Projects
			}
			projects := collectProjects(projectMaps)

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
	statsCmd.Flags().BoolP("monthly", "m", false, "Show monthly totals")
	statsCmd.Flags().BoolP("weekly", "w", false, "Show weekly totals")
	statsCmd.Flags().IntP("rows", "r", 14, "Number of rows to display")

	rootCmd.AddCommand(statsCmd)
}
