package unit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time-tracker/utils"
)

func TestFileStorage_VersionField(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()
	dataFile := filepath.Join(tempDir, "data.json")

	// Create new file storage
	_, err := utils.NewFileStorage(dataFile)
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
