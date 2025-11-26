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
	newEntry := models.TimeEntry{
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

	if len(entries) == 0 {
		return nil, fmt.Errorf("no active time entry to stop")
	}

	lastEntry := &entries[len(entries)-1]
	if !lastEntry.IsRunning() || lastEntry.IsBlank() {
		return nil, fmt.Errorf("no active time entry to stop")
	}

	now := time.Now()
	entries[len(entries)-1].End = &now
	
	// Add a blank entry to represent the gap after the stopped entry
	blankEntry := models.TimeEntry{
		Start:   now,
		End:     nil,
		Project: "",
		Title:   "",
	}
	entries = append(entries, blankEntry)
	
	if err := tm.storage.Save(entries); err != nil {
		return nil, err
	}

	return lastEntry, nil
}

func (tm *TaskManager) ListEntries() ([]models.TimeEntry, error) {
	entries, err := tm.storage.Load()
	if err != nil {
		return nil, err
	}

	// Sort by start time ascending (oldest first)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Start.Before(entries[j].Start)
	})

	return entries, nil
}
