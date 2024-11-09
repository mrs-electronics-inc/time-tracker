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

// pauseCmd represents the pause command
var pauseCmd = &cobra.Command{
	Use:   "pause [task name or ID]",
	Short: "Pause a task",
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

		if task.Status != models.StatusActive {
			return fmt.Errorf("can only pause active tasks")
		}

		now := time.Now()
		var currenPeriodDuration time.Duration

		if task.LastResumeTime.IsZero(){
			currenPeriodDuration = now.Sub(task.StartTime)
		}else{
			currenPeriodDuration = now.Sub(task.LastResumeTime)
		}

		task.AccumulatedTime += currenPeriodDuration
		task.Status = models.StatusPaused
		task.Duration = task.AccumulatedTime.Round(time.Second).String()
		task.LastResumeTime = time.Time{}

		tasks[index] = *task
		if err := taskManager.SaveTasks(tasks); err != nil {
			return err
		}

		fmt.Printf("Paused task: %s (Duration: %s)\n", task.Name, task.Duration)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pauseCmd)
}
