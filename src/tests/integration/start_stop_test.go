package integration

import (
	"testing"
	"time"

	"time-tracker/utils"
)

func TestStartStopScenario(t *testing.T) {
	// Create task manager with memory storage
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)

	// Start tracking
	entry, err := tm.StartEntry("test-project", "Test task")
	if err != nil {
		t.Fatalf("Start entry failed: %v", err)
	}

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// Stop tracking
	stoppedEntry, err := tm.StopEntry()
	if err != nil {
		t.Fatalf("Stop entry failed: %v", err)
	}

	// Check the stopped entry
	if stoppedEntry.ID != entry.ID {
		t.Errorf("Expected stopped entry ID %d, got %d", entry.ID, stoppedEntry.ID)
	}
	if stoppedEntry.End == nil {
		t.Errorf("Expected end time to be set")
	}
	if stoppedEntry.Title != "Test task" {
		t.Errorf("Expected title 'Test task', got %s", stoppedEntry.Title)
	}

	// List entries
	entries, err := tm.ListEntries()
	if err != nil {
		t.Fatalf("List entries failed: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(entries))
	}
	if entries[0].Project != "test-project" {
		t.Errorf("Expected project 'test-project', got %s", entries[0].Project)
	}
}
