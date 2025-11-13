package unit

import (
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
	fs, err := utils.NewFileStorage(dataFile)
	if err != nil {
		t.Fatalf("Failed to create file storage: %v", err)
	}

	// The file should have been created with version 0
	data, err := os.ReadFile(dataFile)
	if err != nil {
		t.Fatalf("Failed to read data file: %v", err)
	}

	// Check that the version field exists in the JSON
	expectedContent := `{
  "version": 0,
  "time-entries": []
}`
	
	if string(data) != expectedContent {
		t.Errorf("Expected data file to contain version field.\nExpected:\n%s\nGot:\n%s", expectedContent, string(data))
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