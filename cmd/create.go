/*
Copyright (c) 2024 Leandro Méndez <leandroa.mendez@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/LeanMendez/time-tracker/config"
	"github.com/LeanMendez/time-tracker/models"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var startOnCreate bool

var createCmd = &cobra.Command{
	Use:   "create [task name]",
	Short: "Create a new task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskName := args[0]
		now := time.Now()

		configData, err := os.ReadFile(config.ConfigFile)
		if err != nil {
			return fmt.Errorf("failed to read config file. Run 'time-tracker init' first: %w", err)
		}

		var config config.Config
		if err := json.Unmarshal(configData, &config); err != nil {
			return fmt.Errorf("failed to parse config: %w", err)
		}

		task := models.Task{
			ID:              uuid.New().String(),
			Name:            taskName,
			Status:          models.StatusNotStarted,
			AccumulatedTime: 0,
			Duration:        "0s",
		}

		if startOnCreate {
			task.Status = models.StatusActive
			task.StartTime = now
			task.LastResumeTime = now
		}

		taskFile := filepath.Join(config.StoragePath, "tasks.json")

		var tasks []models.Task

		if _, err := os.Stat(taskFile); err == nil {
			tasksData, err := os.ReadFile(taskFile)
			if err != nil {
				return fmt.Errorf("failed to read tasks file: %w", err)
			}

			if err := json.Unmarshal(tasksData, &tasks); err != nil {
				return fmt.Errorf("failed to parse tasks: %w", err)
			}
		}

		tasks = append(tasks, task)

		tasksData, err := json.MarshalIndent(tasks, "", " ")
		if err != nil {
			return fmt.Errorf("failed to marshal tasks: %w", err)
		}

		if err := os.WriteFile(taskFile, tasksData, 0644); err != nil {
			return fmt.Errorf("failed to write tasks: %w", err)
		}

		status := "not started"
		if startOnCreate {
			status = "active"
		}

		fmt.Printf("Created new task: %s (Status: %s)\n", taskName, status)
		return nil
	},
}

func init() {
	createCmd.Flags().BoolVarP(&startOnCreate, "start", "s", false, "Start the task immediately after creation")
	rootCmd.AddCommand(createCmd)
}
