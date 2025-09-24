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

func TestStartWithIDScenario(t *testing.T) {
	// Create task manager with memory storage
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)

	// Create an initial entry
	initialEntry, err := tm.StartEntry("existing-project", "Existing task")
	if err != nil {
		t.Fatalf("Start initial entry failed: %v", err)
	}

	// Stop it
	_, err = tm.StopEntry()
	if err != nil {
		t.Fatalf("Stop initial entry failed: %v", err)
	}

	// Now start a new entry using the ID
	newEntry, err := tm.GetEntry(initialEntry.ID)
	if err != nil {
		t.Fatalf("Get entry failed: %v", err)
	}

	// Start new entry with same project and title
	resumedEntry, err := tm.StartEntry(newEntry.Project, newEntry.Title)
	if err != nil {
		t.Fatalf("Start resumed entry failed: %v", err)
	}

	// Check the new entry
	if resumedEntry.Project != "existing-project" {
		t.Errorf("Expected project 'existing-project', got %s", resumedEntry.Project)
	}
	if resumedEntry.Title != "Existing task" {
		t.Errorf("Expected title 'Existing task', got %s", resumedEntry.Title)
	}
	if resumedEntry.End != nil {
		t.Errorf("New entry should be running")
	}

	// List entries
	entries, err := tm.ListEntries()
	if err != nil {
		t.Fatalf("List entries failed: %v", err)
	}

	if len(entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(entries))
	}

	// The initial should be stopped, new running
	var initialStopped, newRunning bool
	for _, e := range entries {
		if e.ID == initialEntry.ID && e.End != nil {
			initialStopped = true
		}
		if e.ID == resumedEntry.ID && e.End == nil {
			newRunning = true
		}
	}
	if !initialStopped {
		t.Errorf("Initial entry should be stopped")
	}
	if !newRunning {
		t.Errorf("New entry should be running")
	}
}
