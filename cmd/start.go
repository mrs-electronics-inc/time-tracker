/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/LeanMendez/time-tracker/models"
	"github.com/LeanMendez/time-tracker/utils"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start [task name or ID]",
	Short: "Start a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskManager, err := utils.NewTaskManager(configFile)
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

		task.Status = models.StatusActive
		if task.StartTime.IsZero() {
			task.StartTime = time.Now()
		}
		task.Duration = "0s"

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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
