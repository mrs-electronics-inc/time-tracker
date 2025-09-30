package cmd

import (
	"fmt"
	"os"
	"time"
	"time-tracker/config"
	"time-tracker/utils"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	weeklyFlag   bool
	projectsFlag bool
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
		if weeklyFlag && projectsFlag {
			return fmt.Errorf("cannot combine --weekly and --projects flags")
		}

		storage, err := utils.NewFileStorage(config.DataFilePath())
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}

		entries, err := storage.Load()
		if err != nil {
			return fmt.Errorf("failed to load entries: %w", err)
		}

		if projectsFlag {
			// Project-based stats
			projectTotals := utils.CalculateProjectTotals(entries)
			if len(projectTotals) == 0 {
				fmt.Println("No data available")
				return nil
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Project", "Total Time"})
			table.SetBorder(true)
			table.SetRowLine(true)
			table.SetAutoWrapText(false)
			for _, total := range projectTotals {
				row := []string{
					total.Project,
					formatTimeHHMM(total.Total),
				}
				table.Append(row)
			}
			table.Render()
		} else {
			// Time-based stats
			if weeklyFlag {
				weeklyTotals := utils.CalculateWeeklyTotals(entries)
				if len(weeklyTotals) == 0 {
					fmt.Println("No data available")
					return nil
				}
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Week Starting", "Total Time"})
				table.SetBorder(true)
				table.SetRowLine(true)
				table.SetAutoWrapText(false)
				for _, total := range weeklyTotals {
					row := []string{
						total.WeekStart.Format("2006-01-02"),
						formatTimeHHMM(total.Total),
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
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Date", "Total Time"})
				table.SetBorder(true)
				table.SetRowLine(true)
				table.SetAutoWrapText(false)
				for _, total := range dailyTotals {
					row := []string{
						total.Date.Format("2006-01-02"),
						formatTimeHHMM(total.Total),
					}
					table.Append(row)
				}
				table.Render()
			}
		}

		return nil
	},
}

func init() {
	statsCmd.Flags().BoolVar(&weeklyFlag, "weekly", false, "Show weekly totals for the past month")
	statsCmd.Flags().BoolVar(&projectsFlag, "projects", false, "Group totals by project")

	rootCmd.AddCommand(statsCmd)
}
