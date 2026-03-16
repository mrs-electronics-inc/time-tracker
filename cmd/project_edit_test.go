package cmd

import (
	"bytes"
	"strings"
	"testing"

	"time-tracker/models"
	"time-tracker/utils"
)

func TestEditProject_RequiresAtLeastOneFlag(t *testing.T) {
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)

	var out bytes.Buffer
	err := editProject(storage, tm, "Acme", "", "", "", false, false, false, &out)
	if err == nil {
		t.Fatal("expected error when no edit flags are provided")
	}

	if !strings.Contains(err.Error(), "at least one flag") {
		t.Fatalf("expected missing flags error, got: %v", err)
	}
}

func TestEditProject_MetadataOnlyPreservesOmittedFieldsAndPrintsSimpleMessage(t *testing.T) {
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)

	if err := storage.SaveProjects([]models.Project{{Name: "Acme", Code: "A1", Category: "Client"}}); err != nil {
		t.Fatalf("failed to seed projects: %v", err)
	}

	var out bytes.Buffer
	err := editProject(storage, tm, "Acme", "", "", "Internal", false, false, true, &out)
	if err != nil {
		t.Fatalf("editProject returned error: %v", err)
	}

	projects, err := storage.LoadProjects()
	if err != nil {
		t.Fatalf("failed to load projects: %v", err)
	}

	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}
	if projects[0].Code != "A1" {
		t.Fatalf("expected omitted code to be preserved, got %q", projects[0].Code)
	}
	if projects[0].Category != "Internal" {
		t.Fatalf("expected category to be updated, got %q", projects[0].Category)
	}

	output := out.String()
	if !strings.Contains(output, "Updated project \"Acme\"") {
		t.Fatalf("expected simple success output, got %q", output)
	}
	if strings.Contains(output, "entries rewritten") {
		t.Fatalf("did not expect rewrite count for metadata-only edit, got %q", output)
	}
}

func TestEditProject_RenameReportsRewrittenCount(t *testing.T) {
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)

	if err := storage.Save([]models.TimeEntry{
		{Project: "Legacy", Title: "A"},
		{Project: "Other", Title: "B"},
		{Project: "Legacy", Title: "C"},
	}); err != nil {
		t.Fatalf("failed to seed entries: %v", err)
	}

	if err := storage.SaveProjects([]models.Project{{Name: "Legacy", Code: "OLD", Category: "Infra"}}); err != nil {
		t.Fatalf("failed to seed projects: %v", err)
	}

	var out bytes.Buffer
	err := editProject(storage, tm, "Legacy", "Current", "", "", true, false, false, &out)
	if err != nil {
		t.Fatalf("editProject returned error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Renamed project \"Legacy\" to \"Current\"") {
		t.Fatalf("expected rename output, got %q", output)
	}
	if !strings.Contains(output, "2 entries rewritten") {
		t.Fatalf("expected rewrite count in output, got %q", output)
	}
}

func TestEditProject_MergeReportsRewrittenCount(t *testing.T) {
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)

	if err := storage.Save([]models.TimeEntry{
		{Project: "Legacy", Title: "A"},
		{Project: "Legacy", Title: "B"},
	}); err != nil {
		t.Fatalf("failed to seed entries: %v", err)
	}

	if err := storage.SaveProjects([]models.Project{
		{Name: "Legacy", Code: "OLD", Category: "Infra"},
		{Name: "Current", Code: "NEW", Category: "Canonical"},
	}); err != nil {
		t.Fatalf("failed to seed projects: %v", err)
	}

	var out bytes.Buffer
	err := editProject(storage, tm, "Legacy", "Current", "", "", true, false, false, &out)
	if err != nil {
		t.Fatalf("editProject returned error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Merged project \"Legacy\" into \"Current\"") {
		t.Fatalf("expected merge output, got %q", output)
	}
	if !strings.Contains(output, "2 entries rewritten") {
		t.Fatalf("expected rewrite count in output, got %q", output)
	}
}
