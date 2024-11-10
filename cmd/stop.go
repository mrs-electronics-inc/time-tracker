/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/LeanMendez/time-tracker/config"
	"github.com/LeanMendez/time-tracker/models"
	"github.com/LeanMendez/time-tracker/utils"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop [task name or ID]",
	Short: "Stop a task and mark it as completed",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskManager, err := utils.NewTaskManager(config.ConfigFile)
		if err != nil {
			return err
		}

		tasks, err := taskManager.LoadTasks()
		if err != nil {
			return err
		}

		task, index, err := taskManager.FindTask(args[0])
		if err != nil {
			return err
		}

		if task.Status == models.StatusCompleted {
			return fmt.Errorf("task is already completed")
		}

		if task.Status == models.StatusNotStarted {
			return fmt.Errorf("cannot stop a task that hasn't been started")
		}

		now := time.Now()
		task.EndTime = now

		if task.Status == models.StatusActive {
			var currentPeriodDuration time.Duration
			if task.LastResumeTime.IsZero() {
				currentPeriodDuration = now.Sub(task.StartTime)
			} else {
				currentPeriodDuration = now.Sub(task.LastResumeTime)
			}
			task.AccumulatedTime += currentPeriodDuration
		}

		task.Status = models.StatusCompleted
		task.Duration = task.AccumulatedTime.Round(time.Second).String()
		task.LastResumeTime = time.Time{}

		tasks[index] = *task
		if err := taskManager.SaveTasks(tasks); err != nil {
			return err
		}

		fmt.Printf("Completed task: %s (Total duration: %s)\n", task.Name, task.Duration)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
