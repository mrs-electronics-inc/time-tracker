package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"time-tracker/config"
	"time-tracker/models"
	"time-tracker/utils"
)

var stopAll bool

var stopCmd = &cobra.Command{
	Use:   "stop [task name or ID]",
	Short: "Stop a task",
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

		if stopAll {
			return stopAllActiveTasks(taskManager, tasks)
		}

		return stopSingleTask(taskManager, tasks, args[0])

	},
}

func stopAllActiveTasks(taskManager *utils.TaskManager, tasks []models.Task) error {
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
		stoppedTask := stopTask(task)

		for i, t := range modifiedTasks {
			if t.ID == stoppedTask.ID {
				modifiedTasks[i] = stoppedTask
				break
			}
		}
	}

	if err := taskManager.SaveTasks(modifiedTasks); err != nil {
		return fmt.Errorf("failed to save tasks: %w", err)
	}

	fmt.Println("All active tasks have been stopped")
	return nil
}

func stopSingleTask(taskManager *utils.TaskManager, tasks []models.Task, taskIdentifier string) error {
	task, index, err := taskManager.FindTask(taskIdentifier)
	if err != nil {
		return fmt.Errorf("failed to find task: %w", err)
	}

	if task.Status != models.StatusActive {
		return fmt.Errorf("can only stop active tasks")
	}

	stoppedTask := stopTask(*task)
	tasks[index] = stoppedTask

	if err := taskManager.SaveTasks(tasks); err != nil {
		return fmt.Errorf("failed to save tasks: %w", err)
	}

	fmt.Printf("Stopped task: %s (Duration: %s)\n", stoppedTask.Name, stoppedTask.Duration)
	return nil
}

func stopTask(task models.Task) models.Task {
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
	stopCmd.Flags().BoolVarP(&stopAll, "all", "a", false, "Stop all the active tasks")
	rootCmd.AddCommand(stopCmd)
}
