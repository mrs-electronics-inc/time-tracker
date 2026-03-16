package utils

import (
	"time-tracker/models"
)

// MemoryStorage implements Storage using in-memory storage for testing
type MemoryStorage struct {
	data     []models.TimeEntry
	projects []models.Project
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data:     []models.TimeEntry{},
		projects: []models.Project{},
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

func (ms *MemoryStorage) LoadProjects() ([]models.Project, error) {
	projects := make([]models.Project, len(ms.projects))
	copy(projects, ms.projects)
	return projects, nil
}

func (ms *MemoryStorage) SaveProjects(projects []models.Project) error {
	ms.projects = make([]models.Project, len(projects))
	copy(ms.projects, projects)
	return nil
}
