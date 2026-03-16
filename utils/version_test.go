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
