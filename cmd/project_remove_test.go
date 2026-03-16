package cmd

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"time-tracker/models"
	"time-tracker/utils"
)

func TestRemoveProject_RemovesUnreferencedProjectAndPrintsSuccess(t *testing.T) {
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)

	if err := storage.SaveProjects([]models.Project{{Name: "Acme", Code: "A", Category: "Client"}}); err != nil {
		t.Fatalf("failed to seed projects: %v", err)
	}

	var out bytes.Buffer
	err := removeProject(tm, "  Acme  ", &out)
	if err != nil {
		t.Fatalf("removeProject returned error: %v", err)
	}

	projects, err := storage.LoadProjects()
	if err != nil {
		t.Fatalf("failed to load projects: %v", err)
	}
	if len(projects) != 0 {
		t.Fatalf("expected project to be removed, got %d project(s)", len(projects))
	}

	if !strings.Contains(out.String(), "Removed project \"Acme\"") {
		t.Fatalf("expected success output, got %q", out.String())
	}
}

func TestRemoveProject_BlockedWhenReferencedIncludesReferenceCount(t *testing.T) {
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)

	now := time.Now()
	later := now.Add(time.Hour)
	if err := storage.Save([]models.TimeEntry{
		{Start: now, End: &later, Project: "Acme", Title: "A"},
		{Start: later, End: &later, Project: "Acme", Title: "B"},
	}); err != nil {
		t.Fatalf("failed to seed entries: %v", err)
	}

	if err := storage.SaveProjects([]models.Project{{Name: "Acme", Code: "A", Category: "Client"}}); err != nil {
		t.Fatalf("failed to seed projects: %v", err)
	}

	var out bytes.Buffer
	err := removeProject(tm, "Acme", &out)
	if err == nil {
		t.Fatal("expected removeProject to fail while referenced")
	}

	if !strings.Contains(err.Error(), "failed to remove project") {
		t.Fatalf("expected remove wrapper error, got %v", err)
	}
	if !strings.Contains(err.Error(), "referenced by 2 time entries") {
		t.Fatalf("expected reference count in error, got %v", err)
	}
}
