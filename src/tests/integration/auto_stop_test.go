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
	entry1, err := tm.StartEntry("project1", "Task 1")
	if err != nil {
		t.Fatalf("Start entry 1 failed: %v", err)
	}

	// Start second task (should auto-stop first)
	entry2, err := tm.StartEntry("project2", "Task 2")
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
	
	// Find the entries (since sorted by start desc, entry2 first)
	var found1, found2 bool
	for _, e := range entries {
		if e.ID == entry1.ID {
			found1 = true
			if e.End == nil {
				t.Errorf("Entry 1 should be stopped")
			}
		}
		if e.ID == entry2.ID {
			found2 = true
			if e.End != nil {
				t.Errorf("Entry 2 should be running")
			}
		}
	}
	if !found1 || !found2 {
		t.Errorf("Both entries not found in list")
	}
}
