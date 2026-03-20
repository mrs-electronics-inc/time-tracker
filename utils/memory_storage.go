package utils

import (
	"sort"
	"strings"
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

	byName := make(map[string]struct{}, len(projects))
	for _, project := range projects {
		byName[project.Name] = struct{}{}
	}

	for _, entry := range ms.data {
		if entry.Project == "" {
			continue
		}
		projectName := strings.TrimSpace(entry.Project)
		if projectName == "" {
			continue
		}
		if _, ok := byName[projectName]; ok {
			continue
		}
		projects = append(projects, models.Project{Name: projectName})
		byName[projectName] = struct{}{}
	}

	return normalizeProjects(projects), nil
}

func normalizeProjects(projects []models.Project) []models.Project {
	type projectGroup struct {
		name string
		key  string
	}

	seen := make(map[string]models.Project, len(projects))
	groups := make([]projectGroup, 0, len(projects))

	for _, project := range projects {
		name := strings.TrimSpace(project.Name)
		if name == "" {
			continue
		}
		key := strings.ToLower(name)
		if _, ok := seen[key]; ok {
			continue
		}
		project.Name = name
		seen[key] = project
		groups = append(groups, projectGroup{name: name, key: key})
	}

	sort.SliceStable(groups, func(i, j int) bool {
		if groups[i].key != groups[j].key {
			return groups[i].key < groups[j].key
		}
		return groups[i].name < groups[j].name
	})

	normalized := make([]models.Project, 0, len(groups))
	for _, group := range groups {
		normalized = append(normalized, seen[group.key])
	}
	return normalized
}

func (ms *MemoryStorage) SaveProjects(projects []models.Project) error {
	ms.projects = make([]models.Project, len(projects))
	copy(ms.projects, projects)
	return nil
}
