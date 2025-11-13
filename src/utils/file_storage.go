package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"path/filepath"
	"time-tracker/models"
)

// DataStore represents the JSON structure for time entries
type DataStore struct {
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
		initialData := DataStore{
			Version:     0,
			TimeEntries: []models.TimeEntry{},
		}
		data, err := json.MarshalIndent(initialData, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal initial data: %w", err)
		}
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return nil, fmt.Errorf("failed to create data file: %w", err)
		}
	}

	return &FileStorage{FilePath: filePath}, nil
}

func (fs *FileStorage) Load() ([]models.TimeEntry, error) {
	data, err := os.ReadFile(fs.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read data file: %w", err)
	}

	var store DataStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("failed to parse data: %w", err)
	}

	// If version is missing, assume version 0
	if store.Version == 0 && len(store.TimeEntries) > 0 {
		// This is likely an old format without version field
		store.Version = 0
	}

	return store.TimeEntries, nil
}

func (fs *FileStorage) Save(entries []models.TimeEntry) error {
	data, err := os.ReadFile(fs.FilePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read existing data file: %w", err)
	}

	var currentVersion int
	if err == nil {
		var existingStore DataStore
		if err := json.Unmarshal(data, &existingStore); err == nil {
			currentVersion = existingStore.Version
		}
	}

	store := DataStore{
		Version:     currentVersion,
		TimeEntries: entries,
	}
	data, err = json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := os.WriteFile(fs.FilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write data file: %w", err)
	}
	return nil
}
