package utils

import (
	"time-tracker/models"
)

// MemoryStorage implements Storage using in-memory storage for testing
type MemoryStorage struct {
	data    []models.TimeEntry
	version int
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data:    []models.TimeEntry{},
		version: 0,
	}
}

func (ms *MemoryStorage) Load() ([]models.TimeEntry, error) {
	// Return a copy to avoid external modification
	entries := make([]models.TimeEntry, len(ms.data))
	copy(entries, ms.data)
	return entries, nil
}

func (ms *MemoryStorage) Save(entries []models.TimeEntry) error {
	// Store a copy
	ms.data = make([]models.TimeEntry, len(entries))
	copy(ms.data, entries)
	return nil
}

func (ms *MemoryStorage) Version() int {
	return ms.version
}

func (ms *MemoryStorage) SetVersion(version int) {
	ms.version = version
}
