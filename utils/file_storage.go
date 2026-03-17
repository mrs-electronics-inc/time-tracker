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
	TimeEntries []models.V4Entry   `json:"time-entries"`
	Projects    []models.V4Project `json:"projects"`
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
			TimeEntries: []models.V4Entry{},
			Projects:    []models.V4Project{},
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
	var v0Entries []models.V0Entry
	var v1Entries []models.V1Entry
	var v2Entries []models.V2Entry
	var v3Entries []models.V3Entry
	var v4Entries []models.V4Entry

	// Step 1: Unmarshal based on version
	switch loadData.Version {
	case 0:
		if err := json.Unmarshal(loadData.TimeEntries, &v0Entries); err != nil {
			return nil, fmt.Errorf("failed to unmarshal v0 data: %w", err)
		}
	case 1:
		if err := json.Unmarshal(loadData.TimeEntries, &v1Entries); err != nil {
			return nil, fmt.Errorf("failed to unmarshal v1 data: %w", err)
		}
	case 2:
		if err := json.Unmarshal(loadData.TimeEntries, &v2Entries); err != nil {
			return nil, fmt.Errorf("failed to unmarshal v2 data: %w", err)
		}
	case 3:
		if err := json.Unmarshal(loadData.TimeEntries, &v3Entries); err != nil {
			return nil, fmt.Errorf("failed to unmarshal v3 data: %w", err)
		}
	case 4:
		if err := json.Unmarshal(loadData.TimeEntries, &v4Entries); err != nil {
			return nil, fmt.Errorf("failed to unmarshal v4 data: %w", err)
		}
	default:
		if loadData.Version > 4 {
			return nil, fmt.Errorf("unknown version: %d", loadData.Version)
		}
	}

	// Step 2: Migrate sequentially to reach CurrentVersion
	for v := loadData.Version; v < models.CurrentVersion; v++ {
		var err error
		switch v {
		case 0:
			v1Entries, err = TransformV0ToV1(v0Entries)
		case 1:
			v2Entries, err = TransformV1ToV2(v1Entries)
		case 2:
			v3Entries, err = TransformV2ToV3(v2Entries)
		case 3:
			v4Entries, err = TransformV3ToV4(v3Entries)
		}
		if err != nil {
			return nil, fmt.Errorf("migration from version %d failed: %w", v, err)
		}
	}

	var entries []models.TimeEntry
	for _, v4 := range v4Entries {
		entries = append(entries, models.TimeEntry{
			Start:   v4.Start,
			Project: v4.Project,
			Title:   v4.Title,
		})
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

	projects, err := fs.LoadProjects()
	if err != nil {
		return err
	}

	return fs.saveEntriesAndProjects(entries, projects)
}

func (fs *FileStorage) saveEntriesAndProjects(entries []models.TimeEntry, projects []models.Project) error {
	saved := toSortedV4Entries(entries)

	data := fileData{
		Version:     models.CurrentVersion,
		TimeEntries: saved,
		Projects:    toV4Projects(projects),
	}

	return fs.writeDataAtomic(data)
}

func toSortedV4Entries(entries []models.TimeEntry) []models.V4Entry {
	// Sort entries by start time before saving
	sorted := append([]models.TimeEntry(nil), entries...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Start.Before(sorted[j].Start)
	})

	saved := make([]models.V4Entry, len(sorted))
	for i, entry := range sorted {
		saved[i] = models.V4Entry{
			Start:   entry.Start,
			Project: entry.Project,
			Title:   entry.Title,
		}
	}

	return saved
}

func toV4Projects(projects []models.Project) []models.V4Project {
	out := make([]models.V4Project, len(projects))
	for i, project := range projects {
		out[i] = models.V4Project{
			Name:     project.Name,
			Code:     project.Code,
			Category: project.Category,
		}
	}
	return out
}

func fromV4Projects(projects []models.V4Project) []models.Project {
	out := make([]models.Project, len(projects))
	for i, project := range projects {
		out[i] = models.Project{
			Name:     project.Name,
			Code:     project.Code,
			Category: project.Category,
		}
	}
	return out
}

func (fs *FileStorage) writeDataAtomic(data fileData) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	dir := filepath.Dir(fs.FilePath)
	tmpFile, err := os.CreateTemp(dir, ".time-tracker-*.json")
	if err != nil {
		return fmt.Errorf("failed to create temp data file: %w", err)
	}

	tmpPath := tmpFile.Name()
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(tmpPath)
		}
	}()

	if err := tmpFile.Chmod(0644); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("failed to set temp file permissions: %w", err)
	}

	if _, err := tmpFile.Write(jsonData); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("failed to write temp data file: %w", err)
	}

	if err := tmpFile.Sync(); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("failed to sync temp data file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp data file: %w", err)
	}

	if err := os.Rename(tmpPath, fs.FilePath); err != nil {
		return fmt.Errorf("failed to atomically replace data file: %w", err)
	}

	cleanup = false
	return nil
}

func (fs *FileStorage) LoadProjects() ([]models.Project, error) {
	jsonData, err := os.ReadFile(fs.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read data file: %w", err)
	}

	var data struct {
		Projects []models.V4Project `json:"projects"`
	}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("failed to parse data: %w", err)
	}
	if data.Projects == nil {
		data.Projects = []models.V4Project{}
	}

	projects := fromV4Projects(data.Projects)
	byName := make(map[string]struct{}, len(projects))
	for _, project := range projects {
		byName[project.Name] = struct{}{}
	}

	entries, err := fs.Load()
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.Project == "" {
			continue
		}
		if _, ok := byName[entry.Project]; ok {
			continue
		}
		projects = append(projects, models.Project{Name: entry.Project})
		byName[entry.Project] = struct{}{}
	}

	return projects, nil
}

func (fs *FileStorage) SaveProjects(projects []models.Project) error {
	entries, err := fs.Load()
	if err != nil {
		return err
	}

	return fs.saveEntriesAndProjects(entries, projects)
}
