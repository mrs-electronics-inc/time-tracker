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

var migrations = map[int]func([]byte) []byte{
	0: MigrateToV1,
}

type fileData struct {
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
		initialData := fileData{
			Version:     models.CurrentVersion,
			TimeEntries: json.RawMessage(`[]`),
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

	// Apply migrations in-memory for older data versions to ensure compatibility.
	// Note: This may add blank entries or other changes that will be persisted if Save() is called.
	for v := data.Version; v < models.CurrentVersion; v++ {
		if mig, ok := migrations[v]; ok {
			data.TimeEntries = mig(data.TimeEntries)
		}
		data.Version++
	}

	var entries []models.TimeEntry
	if err := json.Unmarshal(data.TimeEntries, &entries); err != nil {
		return nil, fmt.Errorf("failed to unmarshal migrated data: %w", err)
	}
	return entries, nil
}

func MigrateToV1(data []byte) []byte {
	var entries []models.TimeEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return data // or handle error
	}
	if len(entries) == 0 {
		return data
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
			blank := models.TimeEntry{
				ID:      maxID + 1,
				Start:   *entry.End,
				End:     &end,
				Project: "",
				Title:   "",
			}
			newEntries = append(newEntries, blank)
			maxID++
		}
	}
	result, _ := json.Marshal(newEntries)
	return result
}

func (fs *FileStorage) Save(entries []models.TimeEntry) error {
	// Saves entries with the current version. If entries were loaded from an older version and migrated,
	// this will upgrade the on-disk format to include migrated changes (e.g., blank entries).
	entriesJson, err := json.Marshal(entries)
	if err != nil {
		return fmt.Errorf("failed to marshal entries: %w", err)
	}
	data := fileData{
		Version:     models.CurrentVersion,
		TimeEntries: entriesJson,
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
