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

		elapse := time.Since(task.StartTime)
		task.Status = models.StatusPaused
		task.Duration = elapse.Round(time.Second).String()

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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pauseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pauseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
