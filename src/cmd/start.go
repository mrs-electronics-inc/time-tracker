package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"time-tracker/config"
	"time-tracker/models"
	"time-tracker/utils"
)

var startCmd = &cobra.Command{
	Use:   "start [task name or ID]",
	Short: "Start a task (creates it if it doesn't exist)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskManager, err := utils.NewTaskManager(config.ConfigFile)
		if err != nil {
			return err
		}

		tasks, err := taskManager.LoadTasks()
		if err != nil {
			// If tasks file doesn't exist, start with empty slice
			tasks = []models.Task{}
		}

		var task *models.Task
		var index int
		for i, t := range tasks {
			if strings.ToLower(t.Name) == strings.ToLower(args[0]) || strings.HasPrefix(t.ID, args[0]) {
				task = &tasks[i]
				index = i
				break
			}
		}

		if task == nil {
			// Task not found, create it
			newTask, err := taskManager.CreateTask(args[0])
			if err != nil {
				return fmt.Errorf("failed to create task: %w", err)
			}
			task = newTask
			// Reload tasks to include the new one
			tasks, err = taskManager.LoadTasks()
			if err != nil {
				return err
			}
			// Find the index of the new task
			for i, t := range tasks {
				if t.ID == task.ID {
					index = i
					break
				}
			}
		}

		if task.Status == models.StatusActive {
			return fmt.Errorf("task is already active")
		}

		if task.Status == models.StatusCompleted {
			return fmt.Errorf("cannot start a completed task")
		}

		now := time.Now()
		task.Status = models.StatusActive
		task.LastResumeTime = now

		if task.StartTime.IsZero() {
			task.StartTime = now
			task.AccumulatedTime = 0
		}

		tasks[index] = *task
		if err := taskManager.SaveTasks(tasks); err != nil {
			return err
		}

		fmt.Printf("Started task: %s\n", task.Name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
