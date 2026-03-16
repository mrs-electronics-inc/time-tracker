package utils

import (
	"errors"
	"strings"
	"testing"
	"time"

	"time-tracker/models"
)

func TestTaskManager_AddProjectTrimsAndEnforcesCaseInsensitiveUniqueness(t *testing.T) {
	storage := NewMemoryStorage()
	tm := NewTaskManager(storage)

	project, err := tm.AddProject("  Auth Refactor  ", " 12572 ", "  Infrastructure ")
	if err != nil {
		t.Fatalf("AddProject failed: %v", err)
	}

	if project.Name != "Auth Refactor" {
		t.Fatalf("expected trimmed name, got %q", project.Name)
	}
	if project.Code != "12572" {
		t.Fatalf("expected trimmed code, got %q", project.Code)
	}
	if project.Category != "Infrastructure" {
		t.Fatalf("expected trimmed category, got %q", project.Category)
	}

	if _, err := tm.AddProject("auth refactor", "", ""); err == nil {
		t.Fatal("expected duplicate project name error")
	} else if !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("expected duplicate error, got: %v", err)
	}
}

func TestTaskManager_EditProjectMergeRewritesEntriesAndKeepsTargetMetadata(t *testing.T) {
	storage := NewMemoryStorage()
	tm := NewTaskManager(storage)

	now := time.Now()
	later := now.Add(time.Hour)
	if err := storage.Save([]models.TimeEntry{
		{Start: now, End: &later, Project: "Legacy", Title: "A"},
		{Start: later, End: &later, Project: "Current", Title: "B"},
		{Start: later.Add(time.Hour), End: &later, Project: "Legacy", Title: "C"},
	}); err != nil {
		t.Fatalf("failed to seed entries: %v", err)
	}

	if err := storage.SaveProjects([]models.Project{
		{Name: "Legacy", Code: "OLD", Category: "Old"},
		{Name: "Current", Code: "NEW", Category: "Canonical"},
	}); err != nil {
		t.Fatalf("failed to seed projects: %v", err)
	}

	result, err := tm.EditProject("Legacy", "Current", "IGNORED", "IGNORED")
	if err != nil {
		t.Fatalf("EditProject failed: %v", err)
	}
	if result.RewrittenEntries != 2 {
		t.Fatalf("expected 2 rewritten entries, got %d", result.RewrittenEntries)
	}
	if !result.Merged {
		t.Fatal("expected merge result")
	}

	entries, err := storage.Load()
	if err != nil {
		t.Fatalf("failed to load entries: %v", err)
	}
	for _, entry := range entries {
		if entry.Project == "Legacy" {
			t.Fatalf("expected all Legacy entries rewritten, found %+v", entry)
		}
	}

	projects, err := storage.LoadProjects()
	if err != nil {
		t.Fatalf("failed to load projects: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("expected 1 project after merge, got %d", len(projects))
	}
	if projects[0].Name != "Current" || projects[0].Code != "NEW" || projects[0].Category != "Canonical" {
		t.Fatalf("expected target metadata to win, got %+v", projects[0])
	}
}

func TestTaskManager_RemoveProjectBlockedWhenReferenced(t *testing.T) {
	storage := NewMemoryStorage()
	tm := NewTaskManager(storage)

	now := time.Now()
	later := now.Add(time.Hour)
	if err := storage.Save([]models.TimeEntry{
		{Start: now, End: &later, Project: "Acme", Title: "Work 1"},
		{Start: later, End: &later, Project: "Acme", Title: "Work 2"},
	}); err != nil {
		t.Fatalf("failed to seed entries: %v", err)
	}

	if err := storage.SaveProjects([]models.Project{{Name: "Acme", Code: "A", Category: "Client"}}); err != nil {
		t.Fatalf("failed to seed projects: %v", err)
	}

	err := tm.RemoveProject("Acme")
	if err == nil {
		t.Fatal("expected remove to fail while referenced")
	}

	var inUseErr *ProjectInUseError
	if !errors.As(err, &inUseErr) {
		t.Fatalf("expected ProjectInUseError, got %T (%v)", err, err)
	}
	if inUseErr.ReferenceCount != 2 {
		t.Fatalf("expected 2 references, got %d", inUseErr.ReferenceCount)
	}
}
