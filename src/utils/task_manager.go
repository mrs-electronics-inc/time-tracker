package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"time-tracker/config"
	"time-tracker/models"
)

type TaskManager struct {
	StoragePath string
}

func NewTaskManager(configFile string) (*TaskManager, error) {
	// Check if config file exists
	configData, err := os.ReadFile(configFile)
	if err != nil {
		// Config file doesn't exist, auto-initialize
		if os.IsNotExist(err) {
			if err := autoInitialize(configFile); err != nil {
				return nil, fmt.Errorf("failed to auto-initialize: %w", err)
			}
			// Try reading again after initialization
			configData, err = os.ReadFile(configFile)
			if err != nil {
				return nil, fmt.Errorf("failed to read config file after initialization: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
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

func autoInitialize(configFile string) error {
	// Get the storage path (same directory as config file)
	storagePath := filepath.Dir(configFile)

	// Create the storage path
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Save the storage path to the configuration file
	configLocal := config.Config{
		StoragePath: storagePath,
	}
	configData, err := json.Marshal(configLocal)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configFile, configData, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Create an empty tasks.json file
	tasksFile := filepath.Join(storagePath, "tasks.json")
	if err := os.WriteFile(tasksFile, []byte("[]"), 0644); err != nil {
		return fmt.Errorf("failed to create tasks file: %w", err)
	}

	return nil
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

func (tm *TaskManager) CreateTask(taskName string) (*models.Task, error) {
	tasks, err := tm.LoadTasks()
	if err != nil {
		// If tasks file doesn't exist, start with empty slice
		if os.IsNotExist(err) {
			tasks = []models.Task{}
		} else {
			return nil, err
		}
	}

	task := models.Task{
		ID:              uuid.New().String(),
		Name:            taskName,
		Status:          models.StatusNotStarted,
		AccumulatedTime: 0,
		Duration:        "0s",
	}

	tasks = append(tasks, task)

	if err := tm.SaveTasks(tasks); err != nil {
		return nil, err
	}

	return &task, nil
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
	// Ensure initialization
	_, err := NewTaskManager(configFile)
	if err != nil {
		return "", fmt.Errorf("failed to initialize task manager: %w", err)
	}

	configData, err := os.ReadFile(configFile)
	if err != nil {
		return "", fmt.Errorf("failed to read config file: %w", err)
	}

	var config config.Config
	if err := json.Unmarshal(configData, &config); err != nil {
		return "", fmt.Errorf("failed to parse config: %w", err)
	}

	taskFile := filepath.Join(config.StoragePath, "tasks.json")
	return taskFile, nil
}
