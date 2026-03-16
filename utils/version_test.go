package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"time-tracker/models"
)

func TestFileStorage_VersionField(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()
	dataFile := filepath.Join(tempDir, "data.json")

	// Create new file storage
	_, err := NewFileStorage(dataFile)
	if err != nil {
		t.Fatalf("Failed to create file storage: %v", err)
	}

	// The file should have been created with the current version
	data, err := os.ReadFile(dataFile)
	if err != nil {
		t.Fatalf("Failed to read data file: %v", err)
	}

	// Parse the JSON and verify the version field exists and is positive
	var payload struct {
		Version int `json:"version"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("Failed to unmarshal data file: %v", err)
	}
	if payload.Version < 1 {
		t.Errorf("Expected data file to contain a positive version field, got %d", payload.Version)
	}
}

func TestFileStorage_InitialSchemaIncludesProjectsKey(t *testing.T) {
	tempDir := t.TempDir()
	dataFile := filepath.Join(tempDir, "data.json")

	_, err := NewFileStorage(dataFile)
	if err != nil {
		t.Fatalf("Failed to create file storage: %v", err)
	}

	data, err := os.ReadFile(dataFile)
	if err != nil {
		t.Fatalf("Failed to read data file: %v", err)
	}

	var payload map[string]json.RawMessage
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("Failed to unmarshal data file: %v", err)
	}

	rawProjects, ok := payload["projects"]
	if !ok {
		t.Fatalf("Expected data file schema to include projects key")
	}

	var projects []models.Project
	if err := json.Unmarshal(rawProjects, &projects); err != nil {
		t.Fatalf("Expected projects key to contain a JSON array: %v", err)
	}
	if len(projects) != 0 {
		t.Fatalf("Expected initial projects to be empty, got %d", len(projects))
	}
}

func TestFileStorage_SavePreservesProjectsKeyInSchema(t *testing.T) {
	tempDir := t.TempDir()
	dataFile := filepath.Join(tempDir, "data.json")

	storage, err := NewFileStorage(dataFile)
	if err != nil {
		t.Fatalf("Failed to create file storage: %v", err)
	}

	entries := []models.TimeEntry{{
		Start:   time.Date(2026, 3, 15, 9, 0, 0, 0, time.UTC),
		Project: "Demo",
		Title:   "Task",
	}}

	if err := storage.Save(entries); err != nil {
		t.Fatalf("Failed to save entries: %v", err)
	}

	data, err := os.ReadFile(dataFile)
	if err != nil {
		t.Fatalf("Failed to read data file: %v", err)
	}

	var payload map[string]json.RawMessage
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("Failed to unmarshal data file: %v", err)
	}

	rawProjects, ok := payload["projects"]
	if !ok {
		t.Fatalf("Expected saved data schema to include projects key")
	}

	var projects []models.Project
	if err := json.Unmarshal(rawProjects, &projects); err != nil {
		t.Fatalf("Expected projects key to contain a JSON array: %v", err)
	}
	if len(projects) != 0 {
		t.Fatalf("Expected saved projects to be empty, got %d", len(projects))
	}
}

func TestFileStorage_SaveAndLoadProjects(t *testing.T) {
	tempDir := t.TempDir()
	dataFile := filepath.Join(tempDir, "data.json")

	storage, err := NewFileStorage(dataFile)
	if err != nil {
		t.Fatalf("Failed to create file storage: %v", err)
	}

	wantProjects := []models.Project{
		{Name: "Acme", Code: "ACM", Category: "Client"},
		{Name: "Internal", Code: "INT", Category: "Ops"},
	}

	if err := storage.SaveProjects(wantProjects); err != nil {
		t.Fatalf("Failed to save projects: %v", err)
	}

	gotProjects, err := storage.LoadProjects()
	if err != nil {
		t.Fatalf("Failed to load projects: %v", err)
	}

	if len(gotProjects) != len(wantProjects) {
		t.Fatalf("Expected %d projects, got %d", len(wantProjects), len(gotProjects))
	}
	for i := range wantProjects {
		if gotProjects[i] != wantProjects[i] {
			t.Fatalf("Project %d mismatch: expected %+v, got %+v", i, wantProjects[i], gotProjects[i])
		}
	}
}

func TestFileStorage_LoadProjectsMissingKeyReturnsEmptySlice(t *testing.T) {
	tempDir := t.TempDir()
	dataFile := filepath.Join(tempDir, "data.json")

	initialData := `{"version":3,"time-entries":[]}`
	if err := os.WriteFile(dataFile, []byte(initialData), 0644); err != nil {
		t.Fatalf("Failed to write data file: %v", err)
	}

	storage, err := NewFileStorage(dataFile)
	if err != nil {
		t.Fatalf("Failed to create file storage: %v", err)
	}

	projects, err := storage.LoadProjects()
	if err != nil {
		t.Fatalf("Failed to load projects: %v", err)
	}
	if projects == nil {
		t.Fatalf("Expected missing projects key to load as empty slice, got nil")
	}
	if len(projects) != 0 {
		t.Fatalf("Expected no projects, got %d", len(projects))
	}
}

func TestFileStorage_SaveEntriesPreservesProjects(t *testing.T) {
	tempDir := t.TempDir()
	dataFile := filepath.Join(tempDir, "data.json")

	storage, err := NewFileStorage(dataFile)
	if err != nil {
		t.Fatalf("Failed to create file storage: %v", err)
	}

	wantProjects := []models.Project{{Name: "Acme", Code: "ACM", Category: "Client"}}
	if err := storage.SaveProjects(wantProjects); err != nil {
		t.Fatalf("Failed to save projects: %v", err)
	}

	entries := []models.TimeEntry{{
		Start:   time.Date(2026, 3, 16, 9, 0, 0, 0, time.UTC),
		Project: "Acme",
		Title:   "Build",
	}}
	if err := storage.Save(entries); err != nil {
		t.Fatalf("Failed to save entries: %v", err)
	}

	gotProjects, err := storage.LoadProjects()
	if err != nil {
		t.Fatalf("Failed to load projects: %v", err)
	}
	if len(gotProjects) != len(wantProjects) || gotProjects[0] != wantProjects[0] {
		t.Fatalf("Expected projects to be preserved, got %+v", gotProjects)
	}
}

func TestFileStorage_SaveProjectsPreservesEntries(t *testing.T) {
	tempDir := t.TempDir()
	dataFile := filepath.Join(tempDir, "data.json")

	storage, err := NewFileStorage(dataFile)
	if err != nil {
		t.Fatalf("Failed to create file storage: %v", err)
	}

	wantEntries := []models.TimeEntry{{
		Start:   time.Date(2026, 3, 16, 9, 0, 0, 0, time.UTC),
		Project: "Acme",
		Title:   "Build",
	}}
	if err := storage.Save(wantEntries); err != nil {
		t.Fatalf("Failed to save entries: %v", err)
	}

	projects := []models.Project{{Name: "Acme", Code: "ACM", Category: "Client"}}
	if err := storage.SaveProjects(projects); err != nil {
		t.Fatalf("Failed to save projects: %v", err)
	}

	gotEntries, err := storage.Load()
	if err != nil {
		t.Fatalf("Failed to load entries: %v", err)
	}
	if len(gotEntries) != len(wantEntries) {
		t.Fatalf("Expected %d entries, got %d", len(wantEntries), len(gotEntries))
	}
	if !gotEntries[0].Start.Equal(wantEntries[0].Start) ||
		gotEntries[0].Project != wantEntries[0].Project ||
		gotEntries[0].Title != wantEntries[0].Title {
		t.Fatalf("Expected entries to be preserved, got %+v", gotEntries)
	}
}

func TestCurrentVersionIsV4(t *testing.T) {
	if models.CurrentVersion != 4 {
		t.Fatalf("Expected CurrentVersion to be 4, got %d", models.CurrentVersion)
	}
}

func TestFileStorage_LoadSupportsV4Data(t *testing.T) {
	tempDir := t.TempDir()
	dataFile := filepath.Join(tempDir, "data.json")

	initialData := `{"version":4,"time-entries":[{"start":"2026-03-16T09:00:00Z","project":"Acme","title":"Build"}],"projects":[{"name":"Acme","code":"ACM","category":"Client"}]}`
	if err := os.WriteFile(dataFile, []byte(initialData), 0644); err != nil {
		t.Fatalf("Failed to write data file: %v", err)
	}

	storage, err := NewFileStorage(dataFile)
	if err != nil {
		t.Fatalf("Failed to create file storage: %v", err)
	}

	entries, err := storage.Load()
	if err != nil {
		t.Fatalf("Failed to load v4 entries: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(entries))
	}
	if entries[0].Project != "Acme" || entries[0].Title != "Build" {
		t.Fatalf("Expected loaded entry fields to match, got %+v", entries[0])
	}
}
