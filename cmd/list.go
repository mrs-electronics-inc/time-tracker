/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/LeanMendez/time-tracker/models"
	"github.com/LeanMendez/time-tracker/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func calculateCurrentDuration(task models.Task) string {
	duration, err := utils.CalculateTaskDuration(task)
	if err != nil {
		return "error"
	}
	return duration.Round(time.Second).String()
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [task name]",
	Short: "List all tasks or a specific task",
	Long: `List command displays tasks in a table format.
If no task name is provided, it shows all tasks.
If a task name is provided, it shows only that specific task.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		configData, err := os.ReadFile(configFile)
		if err != nil {
			return fmt.Errorf("failed to read config file. Run 'timer-cli init' first: %w", err)
		}

		var config Config
		if err := json.Unmarshal(configData, &config); err != nil {
			return fmt.Errorf("failed to parse config: %w", err)
		}

		tasksFile := filepath.Join(config.StoragePath, "tasks.json")
		if _, err := os.Stat(tasksFile); err != nil {
			return fmt.Errorf("no tasks found. Create some tasks first")
		}

		tasksData, err := os.ReadFile(tasksFile)
		if err != nil {
			return fmt.Errorf("failed to read tasks file: %w", err)
		}

		var tasks []models.Task
		if err := json.Unmarshal(tasksData, &tasks); err != nil {
			return fmt.Errorf("failed to parse tasks: %w", err)
		}

		var filteredTasks []models.Task
		if len(args) > 0 {
			searchName := strings.ToLower(args[0])
			for _, task := range tasks {
				if strings.Contains(strings.ToLower(task.Name), searchName) {
					filteredTasks = append(filteredTasks, task)
				}
			}
			if len(filteredTasks) == 0 {
				return fmt.Errorf("no tasks found matching: %s", args[0])
			}
		} else {
			filteredTasks = tasks
		}

		// Create and configure table
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Status", "Start Time", "End Time", "Duration"})
		table.SetBorder(true)
		table.SetRowLine(true)
		table.SetAutoWrapText(false)

		// Add color coding based on status
		table.SetHeaderColor(
			tablewriter.Colors{tablewriter.FgHiWhiteColor},
			tablewriter.Colors{tablewriter.FgHiWhiteColor},
			tablewriter.Colors{tablewriter.FgHiWhiteColor},
			tablewriter.Colors{tablewriter.FgHiWhiteColor},
			tablewriter.Colors{tablewriter.FgHiWhiteColor},
			tablewriter.Colors{tablewriter.FgHiWhiteColor},
		)

		// Format and add rows
		for _, task := range filteredTasks {
			// Format times
			startTime := task.StartTime.Format("2006-01-02 15:04:05")
			endTime := "-"
			if !task.EndTime.IsZero() {
				endTime = task.EndTime.Format("2006-01-02 15:04:05")
			}

			// Calculate current duration for active tasks
			duration := calculateCurrentDuration(task)

			// Truncate ID for better display
			shortID := task.ID[:8]

			row := []string{
				shortID,
				task.Name,
				string(task.Status),
				startTime,
				endTime,
				duration,
			}

			// Add color coding based on status
			var colors []tablewriter.Colors
			switch task.Status {
			case models.StatusActive:
				colors = []tablewriter.Colors{
					{tablewriter.FgGreenColor},
					{tablewriter.FgGreenColor},
					{tablewriter.FgGreenColor},
					{tablewriter.FgGreenColor},
					{tablewriter.FgGreenColor},
					{tablewriter.FgGreenColor},
				}
			case models.StatusPaused:
				colors = []tablewriter.Colors{
					{tablewriter.FgYellowColor},
					{tablewriter.FgYellowColor},
					{tablewriter.FgYellowColor},
					{tablewriter.FgYellowColor},
					{tablewriter.FgYellowColor},
					{tablewriter.FgYellowColor},
				}
			case models.StatusCompleted:
				colors = []tablewriter.Colors{
					{tablewriter.FgBlueColor},
					{tablewriter.FgBlueColor},
					{tablewriter.FgBlueColor},
					{tablewriter.FgBlueColor},
					{tablewriter.FgBlueColor},
					{tablewriter.FgBlueColor},
				}
			}
			table.Rich(row, colors)
		}

		table.Render()
		return nil


	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
