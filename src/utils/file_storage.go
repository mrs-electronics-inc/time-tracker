package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"path/filepath"
	"time-tracker/models"
)

type fileData struct {
	Version     int                `json:"version"`
	TimeEntries []models.TimeEntry `json:"time-entries"`
}

// FileStorage implements Storage using JSON files
type FileStorage struct {
	FilePath string
}

// NewFileStorage creates a new file-based storage, initializing the file if needed
func NewFileStorage(filePath string) (*FileStorage, error) {
	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Ensure data file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		initialData := fileData{
			Version:     0,
			TimeEntries: []models.TimeEntry{},
		}
		jsonData, err := json.MarshalIndent(initialData, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal initial data: %w", err)
		}
		if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
			return nil, fmt.Errorf("failed to create data file: %w", err)
		}
	}

	return &FileStorage{FilePath: filePath}, nil
}

func (fs *FileStorage) Load() ([]models.TimeEntry, error) {
	jsonData, err := os.ReadFile(fs.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read data file: %w", err)
	}

	var data fileData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("failed to parse data: %w", err)
	}

	// TODO: apply migrations based on data.Version

	return data.TimeEntries, nil
}

func (fs *FileStorage) Save(entries []models.TimeEntry) error {
	data := fileData{
		Version:     models.CurrentVersion,
		TimeEntries: entries,
	}
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := os.WriteFile(fs.FilePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write data file: %w", err)
	}
	return nil
}
