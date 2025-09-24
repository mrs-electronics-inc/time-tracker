package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/LeanMendez/time-tracker/config"
	"github.com/LeanMendez/time-tracker/models"
	"github.com/LeanMendez/time-tracker/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [task name]",
	Short: "List all tasks or a specific task",
	Long: `List command displays tasks in a table format.
If no task name is provided, it shows all tasks.
If a task name is provided, it shows only that specific task.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		taskManager, err := utils.NewTaskManager(config.ConfigFile)
		if err != nil {
			return fmt.Errorf("failed to initialize task manager: %w", err)
		}

		tasks, err := taskManager.LoadTasks()
		if err != nil {
			return fmt.Errorf("no tasks found: %w", err)
		}

		if len(tasks) == 0 {
			fmt.Println("no tasks")
			return nil
		}

		filteredTasks := filterTasks(tasks, args)

		if len(filteredTasks) == 0 {
			fmt.Printf("no tasks found matching \"%s\"\n", args[0])
			return nil
		}

		displayTasksTable(filteredTasks)
		return nil
	},
}

func filterTasks(tasks []models.Task, args []string) []models.Task {
	if len(args) == 0 {
		return tasks
	}

	searchName := strings.ToLower(args[0])
	var filteredTasks []models.Task
	for _, task := range tasks {
		if strings.Contains(strings.ToLower(task.Name), searchName) {
			filteredTasks = append(filteredTasks, task)
		}
	}
	return filteredTasks
}

func displayTasksTable(tasks []models.Task) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Status", "Start Time", "End Time", "Duration"})
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetAutoWrapText(false)

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
	)

	// Añadir filas
	for _, task := range tasks {
		startTime := task.StartTime.Format("2006-01-02 15:04:05")
		endTime := "-"
		if !task.EndTime.IsZero() {
			endTime = task.EndTime.Format("2006-01-02 15:04:05")
		}

		duration, err := utils.CalculateTaskDurationString(task)
		if err != nil {
			duration = "error"
		}
		shortID := task.ID[:8]

		row := []string{
			shortID,
			task.Name,
			string(task.Status),
			startTime,
			endTime,
			duration,
		}

		colors := getStatusColor(task.Status)
		table.Rich(row, colors)
	}

	table.Render()
}

func getStatusColor(status models.TaskStatus) []tablewriter.Colors {
	switch status {
	case models.StatusActive:
		return []tablewriter.Colors{
			{tablewriter.FgGreenColor},
			{tablewriter.FgGreenColor},
			{tablewriter.FgGreenColor},
			{tablewriter.FgGreenColor},
			{tablewriter.FgGreenColor},
			{tablewriter.FgGreenColor},
		}
	case models.StatusPaused:
		return []tablewriter.Colors{
			{tablewriter.FgYellowColor},
			{tablewriter.FgYellowColor},
			{tablewriter.FgYellowColor},
			{tablewriter.FgYellowColor},
			{tablewriter.FgYellowColor},
			{tablewriter.FgYellowColor},
		}
	case models.StatusCompleted:
		return []tablewriter.Colors{
			{tablewriter.FgBlueColor},
			{tablewriter.FgBlueColor},
			{tablewriter.FgBlueColor},
			{tablewriter.FgBlueColor},
			{tablewriter.FgBlueColor},
			{tablewriter.FgBlueColor},
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}
