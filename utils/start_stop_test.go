package utils

import (
	"testing"
	"time"
)

func TestStartStopScenario(t *testing.T) {
	// Create task manager with memory storage
	storage := NewMemoryStorage()
	tm := NewTaskManager(storage)

	// Start tracking
	_, err := tm.StartEntry("test-project", "Test task")
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

	// Should have 2 entries: 1 real entry and 1 blank entry after it
	if len(entries) != 2 {
		t.Errorf("Expected 2 entries (real + blank), got %d", len(entries))
	}
	if entries[0].Project != "test-project" {
		t.Errorf("Expected project 'test-project', got %s", entries[0].Project)
	}
	// Second entry should be blank
	if entries[1].Project != "" || entries[1].Title != "" {
		t.Errorf("Expected second entry to be blank, got project=%s title=%s", entries[1].Project, entries[1].Title)
	}
}

func TestMultipleEntriesScenario(t *testing.T) {
	// Create task manager with memory storage
	storage := NewMemoryStorage()
	tm := NewTaskManager(storage)

	// Create and stop an initial entry
	_, err := tm.StartEntry("project1", "task1")
	if err != nil {
		t.Fatalf("Start initial entry failed: %v", err)
	}

	// Stop it
	_, err = tm.StopEntry()
	if err != nil {
		t.Fatalf("Stop initial entry failed: %v", err)
	}

	// Start a second entry
	_, err = tm.StartEntry("project2", "task2")
	if err != nil {
		t.Fatalf("Start second entry failed: %v", err)
	}

	// List entries
	entries, err := tm.ListEntries()
	if err != nil {
		t.Fatalf("List entries failed: %v", err)
	}

	// Should have 3 entries: 1 initial (stopped) + 1 blank (from stop) + 1 second (running)
	if len(entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(entries))
	}

	// Verify entries
	if entries[0].Project != "project1" || entries[0].End == nil {
		t.Errorf("First entry should be project1 and stopped")
	}
	if entries[1].Project != "" || entries[1].Title != "" {
		t.Errorf("Second entry should be blank")
	}
	if entries[2].Project != "project2" || entries[2].End != nil {
		t.Errorf("Third entry should be project2 and running")
	}
}

func TestDuplicateEntryPrevention(t *testing.T) {
	// Create task manager with memory storage
	storage := NewMemoryStorage()
	tm := NewTaskManager(storage)

	// Start tracking a project/task
	_, err := tm.StartEntry("myproject", "mytask")
	if err != nil {
		t.Fatalf("Start entry failed: %v", err)
	}

	// Try to start the same project/task again
	_, err = tm.StartEntry("myproject", "mytask")
	if err == nil {
		t.Errorf("Expected error when starting duplicate entry, got nil")
	}

	// List entries to verify only one entry exists
	entries, err := tm.ListEntries()
	if err != nil {
		t.Fatalf("List entries failed: %v", err)
	}

	// Should have only 1 entry (the running one), no duplicate
	if len(entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(entries))
	}
	if entries[0].Project != "myproject" || entries[0].Title != "mytask" || entries[0].End != nil {
		t.Errorf("Expected single running entry for myproject/mytask")
	}
}

func TestDifferentTaskAllowed(t *testing.T) {
	// Create task manager with memory storage
	storage := NewMemoryStorage()
	tm := NewTaskManager(storage)

	// Start tracking a project/task
	_, err := tm.StartEntry("myproject", "task1")
	if err != nil {
		t.Fatalf("Start first entry failed: %v", err)
	}

	// Starting the same project but different task should succeed (stops previous without blank)
	_, err = tm.StartEntry("myproject", "task2")
	if err != nil {
		t.Fatalf("Start different task failed: %v", err)
	}

	// List entries
	entries, err := tm.ListEntries()
	if err != nil {
		t.Fatalf("List entries failed: %v", err)
	}

	// Should have 2 entries: stopped task1 + running task2 (no blank when starting new)
	if len(entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(entries))
	}
	if entries[0].Project != "myproject" || entries[0].Title != "task1" || entries[0].End == nil {
		t.Errorf("First entry should be stopped task1")
	}
	if entries[1].Project != "myproject" || entries[1].Title != "task2" || entries[1].End != nil {
		t.Errorf("Second entry should be running task2")
	}
}
