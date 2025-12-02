package tui

import (
	"time-tracker/models"
	"time-tracker/utils"
)

// Screen represents the current screen being displayed
type Screen string

const (
	ScreenList Screen = "list"
	ScreenMenu Screen = "menu"
)

// Model represents the state of the TUI application
type Model struct {
	storage      models.Storage
	taskManager  *utils.TaskManager
	entries      []models.TimeEntry
	currentEntry *models.TimeEntry
	currentScreen Screen
	err          error
}

// NewModel creates a new TUI model
func NewModel(storage models.Storage, taskManager *utils.TaskManager) *Model {
	return &Model{
		storage:       storage,
		taskManager:   taskManager,
		currentScreen: ScreenList,
	}
}

// LoadEntries loads time entries from storage
func (m *Model) LoadEntries() error {
	entries, err := m.storage.Load()
	if err != nil {
		return err
	}
	m.entries = entries

	// Find currently running entry
	for i := range entries {
		if entries[i].IsRunning() {
			m.currentEntry = &entries[i]
			break
		}
	}

	return nil
}
