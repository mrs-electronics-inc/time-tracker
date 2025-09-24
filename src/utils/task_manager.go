package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"time-tracker/config"
	"time-tracker/models"
)

type TaskManager struct {
	StoragePath string
}

func NewTaskManager(configFile string) (*TaskManager, error) {
	configData, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file. Run 'time-tracker init' first: %w", err)
	}

	var config struct {
		StoragePath string `json:"storagePath"`
	}
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &TaskManager{
		StoragePath: config.StoragePath,
	}, nil
}

func (tm *TaskManager) LoadTasks() ([]models.Task, error) {
	tasksFile := filepath.Join(tm.StoragePath, "tasks.json")
	tasksData, err := os.ReadFile(tasksFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks file: %w", err)
	}

	var tasks []models.Task
	if err := json.Unmarshal(tasksData, &tasks); err != nil {
		return nil, fmt.Errorf("failed to parse tasks: %w", err)
	}

	return tasks, nil
}

func (tm *TaskManager) SaveTasks(tasks []models.Task) error {
	tasksData, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}

	tasksFile := filepath.Join(tm.StoragePath, "tasks.json")
	if err := os.WriteFile(tasksFile, tasksData, 0644); err != nil {
		return fmt.Errorf("failed to writing the file: %w", err)
	}
	return nil
}

func (tm *TaskManager) FindTask(nameOrID string) (*models.Task, int, error) {
	tasks, err := tm.LoadTasks()
	if err != nil {
		return nil, -1, err
	}

	searchTerm := strings.ToLower(nameOrID)
	for i, task := range tasks {
		if strings.ToLower(task.Name) == searchTerm || strings.HasPrefix(task.ID, searchTerm) {
			return &tasks[i], i, nil
		}
	}

	return nil, -1, fmt.Errorf("task not found: %s", nameOrID)
}

func CalculateTaskDuration(task models.Task) (time.Duration, error) {
	switch task.Status {
	case models.StatusNotStarted:
		return 0, nil
	case models.StatusPaused, models.StatusCompleted:
		return task.AccumulatedTime, nil
	case models.StatusActive:
		if task.LastResumeTime.IsZero() {
			return task.AccumulatedTime + time.Since(task.StartTime), nil
		}
		return task.AccumulatedTime + time.Since(task.LastResumeTime), nil
	default:
		return 0, fmt.Errorf("unknow task status")
	}
}

func CalculateTaskDurationString(task models.Task) (string, error) {
	duration, err := CalculateTaskDuration(task)
	if err != nil {
		return "", err
	}
	return duration.Round(time.Second).String(), nil
}

func RetrieveTaskFile(configFile string) (string, error) {
	configData, err := os.ReadFile(configFile)
	if err != nil {
		return "", fmt.Errorf("failed to read config file. Run 'time-tracker init' first: %w", err)
	}

	var config config.Config
	if err := json.Unmarshal(configData, &config); err != nil {
		return "", fmt.Errorf("failed to parse config: %w", err)
	}

	taskFile := filepath.Join(config.StoragePath, "tasks.json")
	if _, err := os.Stat(taskFile); err != nil {
		return "", fmt.Errorf("no tasks found. Create some tasks first")
	}
	return taskFile, nil
}
