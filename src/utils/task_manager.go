package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"time-tracker/models"
)

type TaskManager struct {
	DataFile string
}

type DataStore struct {
	TimeEntries []models.TimeEntry `json:"time-entries"`
}

func NewTaskManager(dataFile string) (*TaskManager, error) {
	// Ensure data.json exists
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		initialData := DataStore{TimeEntries: []models.TimeEntry{}}
		data, err := json.MarshalIndent(initialData, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal initial data: %w", err)
		}
		if err := os.WriteFile(dataFile, data, 0644); err != nil {
			return nil, fmt.Errorf("failed to create data file: %w", err)
		}
	}

	return &TaskManager{
		DataFile: dataFile,
	}, nil
}

func (tm *TaskManager) LoadTimeEntries() ([]models.TimeEntry, error) {
	data, err := os.ReadFile(tm.DataFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read data file: %w", err)
	}

	var store DataStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("failed to parse data: %w", err)
	}

	return store.TimeEntries, nil
}

func (tm *TaskManager) SaveTimeEntries(entries []models.TimeEntry) error {
	store := DataStore{TimeEntries: entries}
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := os.WriteFile(tm.DataFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write data file: %w", err)
	}
	return nil
}

func (tm *TaskManager) GetNextID() (int, error) {
	entries, err := tm.LoadTimeEntries()
	if err != nil {
		return 0, err
	}

	maxID := 0
	for _, entry := range entries {
		if entry.ID > maxID {
			maxID = entry.ID
		}
	}
	return maxID + 1, nil
}

func (tm *TaskManager) StartEntry(project, title string) (*models.TimeEntry, error) {
	entries, err := tm.LoadTimeEntries()
	if err != nil {
		return nil, err
	}

	// Stop any running entry
	for i, entry := range entries {
		if entry.IsRunning() {
			now := time.Now()
			entries[i].End = &now
			break
		}
	}

	// Create new entry
	nextID, err := tm.GetNextID()
	if err != nil {
		return nil, err
	}

	newEntry := models.TimeEntry{
		ID:      nextID,
		Start:   time.Now(),
		End:     nil,
		Project: project,
		Title:   title,
	}

	entries = append(entries, newEntry)

	if err := tm.SaveTimeEntries(entries); err != nil {
		return nil, err
	}

	return &newEntry, nil
}

func (tm *TaskManager) StopEntry() (*models.TimeEntry, error) {
	entries, err := tm.LoadTimeEntries()
	if err != nil {
		return nil, err
	}

	for i, entry := range entries {
		if entry.IsRunning() {
			now := time.Now()
			entries[i].End = &now
			if err := tm.SaveTimeEntries(entries); err != nil {
				return nil, err
			}
			return &entries[i], nil
		}
	}

	return nil, fmt.Errorf("no active time entry to stop")
}

func (tm *TaskManager) ListEntries() ([]models.TimeEntry, error) {
	entries, err := tm.LoadTimeEntries()
	if err != nil {
		return nil, err
	}

	// Sort by start time descending (newest first)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Start.After(entries[j].Start)
	})

	return entries, nil
}
