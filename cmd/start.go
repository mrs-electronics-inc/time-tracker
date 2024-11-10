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

var startCmd = &cobra.Command{
	Use:   "start [task name or ID]",
	Short: "Start a task",
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
