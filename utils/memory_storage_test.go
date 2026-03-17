package utils

import (
	"testing"
	"time"

	"time-tracker/models"
)

func TestMemoryStorage_LoadProjectsIncludesProjectsFromEntries(t *testing.T) {
	storage := NewMemoryStorage()

	if err := storage.Save([]models.TimeEntry{
		{
			Start:   time.Date(2026, 3, 17, 9, 0, 0, 0, time.UTC),
			Project: "FromEntriesOnly",
			Title:   "Task A",
		},
		{
			Start:   time.Date(2026, 3, 17, 10, 0, 0, 0, time.UTC),
			Project: "WithMetadata",
			Title:   "Task B",
		},
	}); err != nil {
		t.Fatalf("Failed to save entries: %v", err)
	}

	if err := storage.SaveProjects([]models.Project{{Name: "WithMetadata", Code: "WM-1", Category: "Client"}}); err != nil {
		t.Fatalf("Failed to save projects: %v", err)
	}

	projects, err := storage.LoadProjects()
	if err != nil {
		t.Fatalf("Failed to load projects: %v", err)
	}

	if len(projects) != 2 {
		t.Fatalf("Expected 2 projects (metadata + inferred), got %d: %+v", len(projects), projects)
	}

	var foundInferred, foundMetadata bool
	for _, project := range projects {
		if project.Name == "FromEntriesOnly" {
			foundInferred = true
			if project.Code != "" || project.Category != "" {
				t.Fatalf("Expected inferred project metadata to be empty, got %+v", project)
			}
		}
		if project.Name == "WithMetadata" {
			foundMetadata = true
			if project.Code != "WM-1" || project.Category != "Client" {
				t.Fatalf("Expected stored metadata to be preserved, got %+v", project)
			}
		}
	}

	if !foundInferred || !foundMetadata {
		t.Fatalf("Expected both inferred and metadata-backed projects, got %+v", projects)
	}
}
