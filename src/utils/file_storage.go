package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sort"

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
	info, err := os.Stat(filePath)
	if errors.Is(err, fs.ErrNotExist) {
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
	} else if info.IsDir() {
		return nil, errors.New("provided path must be a file")
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

	// Apply migrations based on data.Version
	if data.Version < models.CurrentVersion {
		if data.Version == 0 {
			data.TimeEntries = migrateToV1(data.TimeEntries)
			data.Version = 1
		}
	}

	return data.TimeEntries, nil
}

func migrateToV1(entries []models.TimeEntry) []models.TimeEntry {
	if len(entries) == 0 {
		return entries
	}

	// Sort by start time
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Start.Before(entries[j].Start)
	})

	// Find max ID
	maxID := 0
	for _, e := range entries {
		if e.ID > maxID {
			maxID = e.ID
		}
	}

	var newEntries []models.TimeEntry
	for i, entry := range entries {
		newEntries = append(newEntries, entry)
		if i < len(entries)-1 && entry.End != nil && entry.End.Before(entries[i+1].Start) {
			blank := models.TimeEntry{
				ID:      maxID + 1,
				Start:   *entry.End,
				End:     &entries[i+1].Start,
				Project: "",
				Title:   "",
			}
			newEntries = append(newEntries, blank)
			maxID++
		}
	}
	return newEntries
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
