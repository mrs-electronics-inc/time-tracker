package cmd

import (
	"bytes"
	"strings"
	"testing"

	"time-tracker/utils"
)

func TestAddProject_TrimsInputAndPersistsProject(t *testing.T) {
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)

	var out bytes.Buffer
	err := addProject(tm, "  Auth Refactor  ", " 12572 ", "  Infrastructure ", &out)
	if err != nil {
		t.Fatalf("addProject returned error: %v", err)
	}

	projects, err := storage.LoadProjects()
	if err != nil {
		t.Fatalf("LoadProjects returned error: %v", err)
	}

	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}

	if projects[0].Name != "Auth Refactor" {
		t.Fatalf("expected trimmed name, got %q", projects[0].Name)
	}
	if projects[0].Code != "12572" {
		t.Fatalf("expected trimmed code, got %q", projects[0].Code)
	}
	if projects[0].Category != "Infrastructure" {
		t.Fatalf("expected trimmed category, got %q", projects[0].Category)
	}

	if !strings.Contains(out.String(), "Added project \"Auth Refactor\"") {
		t.Fatalf("expected success output, got: %q", out.String())
	}
}

func TestAddProject_RejectsWhitespaceOnlyName(t *testing.T) {
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)

	var out bytes.Buffer
	err := addProject(tm, "   ", "", "", &out)
	if err == nil {
		t.Fatal("expected error for empty project name")
	}

	if !strings.Contains(err.Error(), "cannot be empty") {
		t.Fatalf("expected empty name error, got: %v", err)
	}
}

func TestAddProject_RejectsCaseInsensitiveDuplicateName(t *testing.T) {
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)

	var out bytes.Buffer
	if err := addProject(tm, "Auth Refactor", "", "", &out); err != nil {
		t.Fatalf("first addProject failed: %v", err)
	}

	err := addProject(tm, "auth refactor", "", "", &out)
	if err == nil {
		t.Fatal("expected duplicate project name error")
	}

	if !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("expected duplicate name error, got: %v", err)
	}
}
