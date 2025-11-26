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

type loadData struct {
	Version     int             `json:"version"`
	TimeEntries json.RawMessage `json:"time-entries"`
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
		// File does not exist, create it with initial data
		initialData := fileData{
			Version:     models.CurrentVersion,
			TimeEntries: []models.TimeEntry{},
		}
		jsonData, err := json.MarshalIndent(initialData, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal initial data: %w", err)
		}
		if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
			return nil, fmt.Errorf("failed to create data file: %w", err)
		}
	} else if err != nil {
		// Stat failed for reasons other than file not existing (e.g., permission error)
		return nil, fmt.Errorf("failed to stat data file: %w", err)
	} else if info.IsDir() {
		// File exists but is a directory, not a file
		return nil, errors.New("provided path must be a file, not a directory")
	}

	return &FileStorage{FilePath: filePath}, nil
}

func (fs *FileStorage) Load() ([]models.TimeEntry, error) {
	jsonData, err := os.ReadFile(fs.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read data file: %w", err)
	}

	var loadData loadData
	if err := json.Unmarshal(jsonData, &loadData); err != nil {
		return nil, fmt.Errorf("failed to parse data: %w", err)
	}

	// Apply migrations in-memory for older data versions to ensure compatibility.
	entriesJson := loadData.TimeEntries
	for v := loadData.Version; v < models.CurrentVersion; v++ {
		if transform, ok := Transformations[v]; ok {
			var err error
			entriesJson, err = callTransformWithMarshal(entriesJson, transform)
			if err != nil {
				return nil, fmt.Errorf("migration from version %d failed: %w", v, err)
			}
		}
		loadData.Version++
	}

	var entries []models.TimeEntry
	if err := json.Unmarshal(entriesJson, &entries); err != nil {
		return nil, fmt.Errorf("failed to unmarshal migrated data: %w", err)
	}

	// Reconstruct End times from next entry's Start time for all entries
	for i := 0; i < len(entries)-1; i++ {
		next := entries[i+1].Start
		entries[i].End = &next
	}
	// Last entry has no End time
	if len(entries) > 0 {
		entries[len(entries)-1].End = nil
	}

	return entries, nil
}

func (fs *FileStorage) Save(entries []models.TimeEntry) error {
	// Saves entries with the current version. If entries were loaded from an
	// older version and migrated, this will upgrade the on-disk format to
	// include migrated changes (e.g., blank entries).

	// Sort entries by start time before saving
	sorted := append([]models.TimeEntry(nil), entries...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Start.Before(sorted[j].Start)
	})

	saved := make([]models.V3Entry, len(sorted))
	for i, entry := range sorted {
		saved[i] = models.V3Entry{
			Start:   entry.Start,
			Project: entry.Project,
			Title:   entry.Title,
		}
	}

	data := map[string]any{
		"version":      models.CurrentVersion,
		"time-entries": saved,
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
