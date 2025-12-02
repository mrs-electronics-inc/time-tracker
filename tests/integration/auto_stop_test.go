package integration

import (
	"testing"

	"time-tracker/utils"
)

func TestAutoStopScenario(t *testing.T) {
	// Create task manager with memory storage
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)

	// Start first task
	_, err := tm.StartEntry("project1", "Task 1")
	if err != nil {
		t.Fatalf("Start entry 1 failed: %v", err)
	}

	// Start second task (should auto-stop first)
	_, err = tm.StartEntry("project2", "Task 2")
	if err != nil {
		t.Fatalf("Start entry 2 failed: %v", err)
	}

	// List entries
	entries, err := tm.ListEntries()
	if err != nil {
		t.Fatalf("List entries failed: %v", err)
	}

	// Should have 2 entries, first stopped, second running
	if len(entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(entries))
	}
	
	// Verify entries by order (since sorted by start time ascending)
	if entries[0].Project != "project1" || entries[0].End == nil {
		t.Errorf("First entry should be project1 and stopped")
	}
	if entries[1].Project != "project2" || entries[1].End != nil {
		t.Errorf("Second entry should be project2 and running")
	}
}
