package tui

import (
	"fmt"
	"testing"
	"time"

	"time-tracker/models"
	"time-tracker/utils"
)

// TestStartEntryViaUICreatesFileEntry verifies starting via TUI saves to storage
func TestStartEntryViaUICreatesFileEntry(t *testing.T) {
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)
	m := NewModel(storage, tm)

	// Load initial state
	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load entries: %v", err)
	}

	// Simulate starting an entry via UI
	if _, err := tm.StartEntry("integration-project", "integration-task"); err != nil {
		t.Fatalf("Failed to start entry: %v", err)
	}

	// Load entries again
	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to reload entries: %v", err)
	}

	// Verify entry exists in storage
	entries, err := storage.Load()
	if err != nil {
		t.Fatalf("Failed to load from storage: %v", err)
	}

	found := false
	for _, entry := range entries {
		if entry.Project == "integration-project" && entry.Title == "integration-task" && entry.IsRunning() {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find started entry in storage")
	}
}

// TestStopEntryViaUIUpdatesFileEntry verifies stopping via TUI updates storage
func TestStopEntryViaUIUpdatesFileEntry(t *testing.T) {
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)
	m := NewModel(storage, tm)

	// Start an entry
	if _, err := tm.StartEntry("stop-project", "stop-task"); err != nil {
		t.Fatalf("Failed to start entry: %v", err)
	}

	// Stop the entry
	stoppedEntry, err := tm.StopEntry()
	if err != nil {
		t.Fatalf("Failed to stop entry: %v", err)
	}

	// Verify the entry is stopped
	if stoppedEntry.End == nil {
		t.Error("Expected entry to have end time")
	}

	// Load and verify in storage
	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load entries: %v", err)
	}

	// Find the stopped entry
	found := false
	for _, entry := range m.Entries {
		if entry.Project == "stop-project" && entry.Title == "stop-task" && entry.End != nil {
			found = true
			// Verify duration is reasonable
			duration := entry.Duration()
			if duration < 0 {
				t.Error("Expected positive duration")
			}
			break
		}
	}

	if !found {
		t.Error("Expected to find stopped entry in list")
	}

	// Verify persistence
	entries, err := storage.Load()
	if err != nil {
		t.Fatalf("Failed to load from storage: %v", err)
	}

	persistFound := false
	for _, entry := range entries {
		if entry.Project == "stop-project" && entry.Title == "stop-task" && entry.End != nil {
			persistFound = true
			break
		}
	}

	if !persistFound {
		t.Error("Expected stopped entry to persist in storage")
	}
}

// TestLoadRecentEntriesMatchesCLIList verifies TUI list matches CLI list output
func TestLoadRecentEntriesMatchesCLIList(t *testing.T) {
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)
	m := NewModel(storage, tm)

	// Create several entries
	entries := []struct {
		project string
		title   string
	}{
		{"proj1", "task1"},
		{"proj2", "task2"},
		{"proj3", "task3"},
	}

	for _, e := range entries {
		if _, err := tm.StartEntry(e.project, e.title); err != nil {
			t.Fatalf("Failed to start entry %s/%s: %v", e.project, e.title, err)
		}
	}

	// Load via TUI
	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load entries: %v", err)
	}

	// Get via CLI (direct storage)
	cliEntries, err := tm.ListEntries()
	if err != nil {
		t.Fatalf("Failed to list entries via CLI: %v", err)
	}

	// Verify counts match
	if len(m.Entries) != len(cliEntries) {
		t.Errorf("Expected same entry count: TUI=%d, CLI=%d", len(m.Entries), len(cliEntries))
	}

	// Verify entries match
	for i, tuiEntry := range m.Entries {
		if i >= len(cliEntries) {
			t.Errorf("TUI has more entries than CLI at index %d", i)
			break
		}
		cliEntry := cliEntries[i]

		if tuiEntry.Project != cliEntry.Project {
			t.Errorf("Project mismatch at %d: TUI=%s, CLI=%s", i, tuiEntry.Project, cliEntry.Project)
		}
		if tuiEntry.Title != cliEntry.Title {
			t.Errorf("Title mismatch at %d: TUI=%s, CLI=%s", i, tuiEntry.Title, cliEntry.Title)
		}
		if tuiEntry.Start.Unix() != cliEntry.Start.Unix() {
			t.Errorf("Start time mismatch at %d", i)
		}
	}
}

// TestDataConsistencyAfterMultipleOperations verifies data consistency
func TestDataConsistencyAfterMultipleOperations(t *testing.T) {
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)
	m := NewModel(storage, tm)

	// Sequence of operations
	// 1. Start entry A
	if _, err := tm.StartEntry("projectA", "taskA"); err != nil {
		t.Fatalf("Failed to start A: %v", err)
	}

	// 2. Start entry B (stops A)
	if _, err := tm.StartEntry("projectB", "taskB"); err != nil {
		t.Fatalf("Failed to start B: %v", err)
	}

	// 3. Stop entry B
	if _, err := tm.StopEntry(); err != nil {
		t.Fatalf("Failed to stop B: %v", err)
	}

	// 4. Start entry C
	if _, err := tm.StartEntry("projectC", "taskC"); err != nil {
		t.Fatalf("Failed to start C: %v", err)
	}

	// Load via TUI
	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load entries: %v", err)
	}

	// Verify consistency
	if len(m.Entries) < 4 {
		t.Errorf("Expected at least 4 entries (A, B, blank, C), got %d", len(m.Entries))
	}

	// Find and verify each entry
	entryMap := make(map[string]*models.TimeEntry)
	for i := range m.Entries {
		if m.Entries[i].Project != "" {
			entryMap[m.Entries[i].Project] = &m.Entries[i]
		}
	}

	// A should be stopped
	if a, ok := entryMap["projectA"]; !ok || a.End == nil {
		t.Error("Expected projectA to be stopped")
	}

	// B should be stopped
	if b, ok := entryMap["projectB"]; !ok || b.End == nil {
		t.Error("Expected projectB to be stopped")
	}

	// C should be running
	if c, ok := entryMap["projectC"]; !ok || c.End != nil {
		t.Error("Expected projectC to be running")
	}

	// Verify all have positive durations or are running
	for _, entry := range m.Entries {
		if entry.Project != "" { // Skip blank entries
			if !entry.IsRunning() {
				duration := entry.Duration()
				if duration < 0 {
					t.Errorf("Entry %s/%s has negative duration: %v", entry.Project, entry.Title, duration)
				}
			}
		}
	}
}

// TestEdgeCaseNoData verifies behavior with no data
func TestEdgeCaseNoData(t *testing.T) {
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)
	m := NewModel(storage, tm)

	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load: %v", err)
	}

	if len(m.Entries) != 0 {
		t.Error("Expected empty list with no data")
	}

	if m.SelectedIdx != 0 {
		t.Error("Expected selection at 0 with empty list")
	}

	// View should not crash
	view := m.View()
	if view == "" {
		t.Error("Expected non-empty view even with no data")
	}
}

// TestEdgeCaseRapidStartStop verifies handling of rapid sequential operations
func TestEdgeCaseRapidStartStop(t *testing.T) {
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)
	m := NewModel(storage, tm)

	// Rapid sequential start/stop cycles
	for i := 0; i < 3; i++ {
		projName := fmt.Sprintf("proj%d", i)
		taskName := fmt.Sprintf("task%d", i)

		if _, err := tm.StartEntry(projName, taskName); err != nil {
			t.Fatalf("Failed to start %d: %v", i, err)
		}

		if _, err := tm.StopEntry(); err != nil {
			t.Fatalf("Failed to stop %d: %v", i, err)
		}
	}

	// Verify all operations succeeded
	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load entries: %v", err)
	}

	if len(m.Entries) < 6 {
		t.Errorf("Expected at least 6 entries (3 real + 3 blanks), got %d", len(m.Entries))
	}

	// Verify no data corruption
	for _, entry := range m.Entries {
		if entry.Project != "" && entry.Title != "" {
			// This is a real entry, should have start time
			if entry.Start.IsZero() {
				t.Errorf("Entry %s/%s has zero start time", entry.Project, entry.Title)
			}
		}
	}
}

// TestDurationCalculationAfterStop verifies duration calculation is correct
func TestDurationCalculationAfterStop(t *testing.T) {
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)
	m := NewModel(storage, tm)

	// Start entry
	if _, err := tm.StartEntry("duration-test", "task"); err != nil {
		t.Fatalf("Failed to start: %v", err)
	}

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// Stop entry
	stoppedEntry, err := tm.StopEntry()
	if err != nil {
		t.Fatalf("Failed to stop: %v", err)
	}

	// Reload to get from storage
	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load: %v", err)
	}

	// Find the stopped entry
	var found *models.TimeEntry
	for i := range m.Entries {
		if m.Entries[i].Project == "duration-test" && m.Entries[i].Title == "task" && m.Entries[i].End != nil {
			found = &m.Entries[i]
			break
		}
	}

	if found == nil {
		t.Fatal("Expected to find stopped entry")
	}

	duration := found.Duration()

	// Duration should be at least 100ms (what we slept), but allow some overhead
	if duration < 50*time.Millisecond {
		t.Errorf("Duration too short: %v", duration)
	}

	// Verify it matches direct calculation
	expectedDuration := stoppedEntry.Duration()
	if duration != expectedDuration {
		t.Errorf("Duration mismatch: loaded=%v, returned=%v", duration, expectedDuration)
	}
}

// TestBlankEntriesAreTracked verifies blank entries are properly handled
func TestBlankEntriesAreTracked(t *testing.T) {
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)
	m := NewModel(storage, tm)

	// Start and stop an entry
	if _, err := tm.StartEntry("test", "task"); err != nil {
		t.Fatalf("Failed to start: %v", err)
	}

	if _, err := tm.StopEntry(); err != nil {
		t.Fatalf("Failed to stop: %v", err)
	}

	// Load entries
	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load: %v", err)
	}

	// Should have real entry + blank entry at end
	if len(m.Entries) < 2 {
		t.Errorf("Expected at least 2 entries (real + blank), got %d", len(m.Entries))
	}

	// Last entry should be blank
	lastEntry := m.Entries[len(m.Entries)-1]
	if !lastEntry.IsBlank() {
		t.Error("Expected last entry to be blank (marked end of day)")
	}
}

// TestCLIAndTUIMakeConsistentData verifies CLI and TUI produce same results
func TestCLIAndTUIMakeConsistentData(t *testing.T) {
	// Create two identical scenarios - one via CLI, one via TUI model
	storageCLI := utils.NewMemoryStorage()
	tmCLI := utils.NewTaskManager(storageCLI)

	storageTUI := utils.NewMemoryStorage()
	tmTUI := utils.NewTaskManager(storageTUI)
	modelTUI := NewModel(storageTUI, tmTUI)

	// Identical operations
	operations := []struct {
		op    string // "start" or "stop"
		proj  string
		title string
	}{
		{"start", "proj1", "task1"},
		{"start", "proj2", "task2"},
		{"stop", "", ""},
		{"start", "proj3", "task3"},
	}

	for _, op := range operations {
		if op.op == "start" {
			if _, err := tmCLI.StartEntry(op.proj, op.title); err != nil {
				t.Fatalf("CLI start failed: %v", err)
			}
			if _, err := tmTUI.StartEntry(op.proj, op.title); err != nil {
				t.Fatalf("TUI start failed: %v", err)
			}
		} else {
			if _, err := tmCLI.StopEntry(); err != nil {
				t.Fatalf("CLI stop failed: %v", err)
			}
			if _, err := tmTUI.StopEntry(); err != nil {
				t.Fatalf("TUI stop failed: %v", err)
			}
		}
	}

	// Load both
	cliEntries, err := tmCLI.ListEntries()
	if err != nil {
		t.Fatalf("CLI list failed: %v", err)
	}

	if err := modelTUI.LoadEntries(); err != nil {
		t.Fatalf("TUI load failed: %v", err)
	}

	// Compare
	if len(cliEntries) != len(modelTUI.Entries) {
		t.Errorf("Entry count mismatch: CLI=%d, TUI=%d", len(cliEntries), len(modelTUI.Entries))
	}

	for i := range cliEntries {
		if i >= len(modelTUI.Entries) {
			break
		}
		cli := cliEntries[i]
		tui := modelTUI.Entries[i]

		if cli.Project != tui.Project || cli.Title != tui.Title {
			t.Errorf("Entry %d mismatch: CLI=(%s,%s) TUI=(%s,%s)", i, cli.Project, cli.Title, tui.Project, tui.Title)
		}
	}
}
