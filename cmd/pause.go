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

var pauseAll bool

var pauseCmd = &cobra.Command{
	Use:   "pause [task name or ID]",
	Short: "Pause a task",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskManager, err := utils.NewTaskManager(config.ConfigFile)
		if err != nil {
			return err
		}

		tasks, err := taskManager.LoadTasks()
		if err != nil {
			return err
		}

		if pauseAll {
			return pauseAllActiveTasks(taskManager, tasks)
		}

		return pauseSingleTask(taskManager, tasks, args[0])

	},
}

func pauseAllActiveTasks(taskManager *utils.TaskManager, tasks []models.Task) error {
	var activeTasks []models.Task
	var modifiedTasks []models.Task = make([]models.Task, len(tasks))
	copy(modifiedTasks, tasks)

	for _, task := range tasks {
		if task.Status == models.StatusActive {
			activeTasks = append(activeTasks, task)
		}
	}

	if len(activeTasks) == 0 {
		return fmt.Errorf("no active tasks found")
	}

	for _, task := range activeTasks {
		pausedTask := pauseTask(task)

		for i, t := range modifiedTasks {
			if t.ID == pausedTask.ID {
				modifiedTasks[i] = pausedTask
				break
			}
		}
	}

	if err := taskManager.SaveTasks(modifiedTasks); err != nil {
		return fmt.Errorf("failed to save tasks: %w", err)
	}

	fmt.Println("All active tasks have been paused")
	return nil
}

func pauseSingleTask(taskManager *utils.TaskManager, tasks []models.Task, taskIdentifier string) error {
	task, index, err := taskManager.FindTask(taskIdentifier)
	if err != nil {
		return fmt.Errorf("failed to find task: %w", err)
	}

	if task.Status != models.StatusActive {
		return fmt.Errorf("can only pause active tasks")
	}

	pausedTask := pauseTask(*task)
	tasks[index] = pausedTask

	if err := taskManager.SaveTasks(tasks); err != nil {
		return fmt.Errorf("failed to save tasks: %w", err)
	}

	fmt.Printf("Paused task: %s (Duration: %s)\n", pausedTask.Name, pausedTask.Duration)
	return nil
}

func pauseTask(task models.Task) models.Task {
	now := time.Now()
	var currentPeriodDuration time.Duration

	if task.LastResumeTime.IsZero() {
		currentPeriodDuration = now.Sub(task.StartTime)
	} else {
		currentPeriodDuration = now.Sub(task.LastResumeTime)
	}

	task.AccumulatedTime += currentPeriodDuration
	task.Status = models.StatusPaused
	task.Duration = task.AccumulatedTime.Round(time.Second).String()
	task.LastResumeTime = time.Time{}

	return task
}

func init() {
	pauseCmd.Flags().BoolVarP(&pauseAll, "all", "a", false, "Pause all the active tasks")
	rootCmd.AddCommand(pauseCmd)
}
