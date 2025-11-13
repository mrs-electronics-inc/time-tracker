package unit

import (
	"os"
	"path/filepath"
	"strings"
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

	// The file should have been created with version 0
	data, err := os.ReadFile(dataFile)
	if err != nil {
		t.Fatalf("Failed to read data file: %v", err)
	}

	// Check that the version field exists in the JSON - just check it contains version
	if !strings.Contains(string(data), `"version": 0`) {
		t.Errorf("Expected data file to contain version field. Got: %s", string(data))
	}
}

func TestMemoryStorage_Version(t *testing.T) {
	ms := utils.NewMemoryStorage()
	
	// Test initial version
	if ms.Version() != 0 {
		t.Errorf("Expected initial version to be 0, got %d", ms.Version())
	}
	
	// Test setting version
	ms.SetVersion(1)
	if ms.Version() != 1 {
		t.Errorf("Expected version to be 1 after setting, got %d", ms.Version())
	}
}