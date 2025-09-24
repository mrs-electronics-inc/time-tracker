package utils

import (
	"fmt"
	"sort"
	"time"

	"time-tracker/models"
)

// Storage interface for abstracting data persistence
type Storage interface {
	Load() ([]models.TimeEntry, error)
	Save([]models.TimeEntry) error
}

type TaskManager struct {
	storage Storage
}

func NewTaskManager(storage Storage) *TaskManager {
	return &TaskManager{storage: storage}
}

func (tm *TaskManager) GetNextID() (int, error) {
	entries, err := tm.storage.Load()
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
	entries, err := tm.storage.Load()
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

	if err := tm.storage.Save(entries); err != nil {
		return nil, err
	}

	return &newEntry, nil
}

func (tm *TaskManager) StopEntry() (*models.TimeEntry, error) {
	entries, err := tm.storage.Load()
	if err != nil {
		return nil, err
	}

	for i, entry := range entries {
		if entry.IsRunning() {
			now := time.Now()
			entries[i].End = &now
			if err := tm.storage.Save(entries); err != nil {
				return nil, err
			}
			return &entries[i], nil
		}
	}

	return nil, fmt.Errorf("no active time entry to stop")
}

func (tm *TaskManager) ListEntries() ([]models.TimeEntry, error) {
	entries, err := tm.storage.Load()
	if err != nil {
		return nil, err
	}

	// Sort by start time descending (newest first)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Start.After(entries[j].Start)
	})

	return entries, nil
}
