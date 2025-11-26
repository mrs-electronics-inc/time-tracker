package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"time"

	"path/filepath"
	"time-tracker/models"
)

var migrations = map[int]func([]byte) ([]byte, error){
	0: MigrateToV1,
	1: MigrateToV2,
}

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
		if mig, ok := migrations[v]; ok {
			var err error
			entriesJson, err = mig(entriesJson)
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
	return entries, nil
}

func MigrateToV1(data []byte) ([]byte, error) {
	var entries []models.TimeEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data during migration to v1: %w", err)
	}
	if len(entries) == 0 {
		return data, nil
	}

	// Make a shallow copy to avoid mutating the input
	copied := append([]models.TimeEntry(nil), entries...)

	// Sort by start time
	sort.Slice(copied, func(i, j int) bool {
		return copied[i].Start.Before(copied[j].Start)
	})

	// Find max ID
	maxID := 0
	for _, e := range copied {
		if e.ID > maxID {
			maxID = e.ID
		}
	}

	var newEntries []models.TimeEntry
	for i, entry := range copied {
		newEntries = append(newEntries, entry)
		if i < len(copied)-1 && entry.End != nil && entry.End.Before(copied[i+1].Start) {
			end := copied[i+1].Start
			maxID++
			blank := models.TimeEntry{
				ID:      maxID,
				Start:   *entry.End,
				End:     &end,
				Project: "",
				Title:   "",
			}
			newEntries = append(newEntries, blank)
		}
	}
	result, err := json.Marshal(newEntries)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal migrated data: %w", err)
	}
	return result, nil
}

func MigrateToV2(data []byte) ([]byte, error) {
	var entries []models.TimeEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data during migration to v2: %w", err)
	}

	// Filter out blank entries (empty project and title) that are less than 5 seconds long
	var filtered []models.TimeEntry
	for _, entry := range entries {
		// Skip if it's a blank entry with duration < 5 seconds
		if entry.Project == "" && entry.Title == "" && entry.End != nil {
			duration := entry.End.Sub(entry.Start)
			if duration < 5*time.Second {
				continue
			}
		}
		filtered = append(filtered, entry)
	}

	// Remove End field from filtered entries (in v2, End is calculated from the next entry's Start)
	var newEntries []models.TimeEntry
	for _, entry := range filtered {
		newEntry := entry
		newEntry.End = nil
		newEntries = append(newEntries, newEntry)
	}

	result, err := json.Marshal(newEntries)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal migrated data: %w", err)
	}
	return result, nil
}

func (fs *FileStorage) Save(entries []models.TimeEntry) error {
	// Saves entries with the current version. If entries were loaded from an older version and migrated,
	// this will upgrade the on-disk format to include migrated changes (e.g., blank entries).
	entriesJson, err := json.Marshal(entries)
	if err != nil {
		return fmt.Errorf("failed to marshal entries: %w", err)
	}

	// For version 2+, remove the End field from the JSON
	if models.CurrentVersion >= 2 {
		var processedEntries []map[string]any
		if err := json.Unmarshal(entriesJson, &processedEntries); err != nil {
			return fmt.Errorf("failed to unmarshal entries for processing: %w", err)
		}
		for _, entry := range processedEntries {
			delete(entry, "end")
		}
		entriesJson, err = json.Marshal(processedEntries)
		if err != nil {
			return fmt.Errorf("failed to marshal processed entries: %w", err)
		}
	}

	// Manually reconstruct the fileData with the processed entries JSON
	fileDataJson := map[string]any{
		"version":       models.CurrentVersion,
		"time-entries":  json.RawMessage(entriesJson),
	}
	jsonData, err := json.MarshalIndent(fileDataJson, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := os.WriteFile(fs.FilePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write data file: %w", err)
	}
	return nil
}
