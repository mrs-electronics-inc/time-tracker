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

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop [task name or ID]",
	Short: "Stop a task and mark it as completed",
	Args: cobra.ExactArgs(1),
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

		if task.Status == models.StatusCompleted {
			return fmt.Errorf("task is already completed")
		}

		if task.Status == models.StatusNotStarted {
			return fmt.Errorf("cannot stop a task that hasn't been started")
		}

		endTime := time.Now()
		var elapsed time.Duration

		if task.Status == models.StatusActive {
			elapsed = endTime.Sub(task.StartTime)
		}else{
			existingDuration,_ := time.ParseDuration(task.Duration)
			elapsed = existingDuration
		}

		task.Status = models.StatusCompleted
		task.EndTime = endTime
		task.Duration = elapsed.Round(time.Second).String()

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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
