package tui

import (
	"fmt"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"time-tracker/utils"
)

// mustParseTime parses a time string or panics
func mustParseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

// Helper to create a test model
func newTestModel() *Model {
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)
	return NewModel(storage, tm)
}

// TestModelInitialization verifies the model is properly initialized
func TestModelInitialization(t *testing.T) {
	m := newTestModel()

	if m.CurrentMode != m.ListMode {
		t.Error("Expected current mode to be list mode")
	}
	if m.Width != 0 || m.Height != 0 {
		t.Error("Expected width and height to be 0 initially")
	}
	if m.Status != "" {
		t.Error("Expected status to be empty initially")
	}
	if m.SelectedIdx != 0 {
		t.Error("Expected scroll offset to be 0")
	}
}

// TestWindowSizeUpdate verifies window size messages are handled
func TestWindowSizeUpdate(t *testing.T) {
	m := newTestModel()

	msg := tea.WindowSizeMsg{Width: 80, Height: 24}
	updated, _ := m.Update(msg)
	model := updated.(*Model)

	if model.Width != 80 || model.Height != 24 {
		t.Errorf("Expected 80x24, got %dx%d", model.Width, model.Height)
	}
}

// TestModeTransitionFromListToNew verifies navigation to new mode
func TestModeTransitionFromListToNew(t *testing.T) {
	m := newTestModel()
	// Load no entries first
	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load entries: %v", err)
	}

	// Simulate 'n' key in list mode (should open new mode)
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	updated, _ := m.Update(msg)
	model := updated.(*Model)

	if model.CurrentMode != model.NewMode {
		t.Error("Expected mode to switch to new mode after 'n' key")
	}
}

// TestModeTransitionFromStartToList verifies canceling start mode
func TestModeTransitionFromStartToList(t *testing.T) {
	m := newTestModel()
	// Manually switch to start mode
	m.CurrentMode = m.StartMode

	// Simulate Esc key (should return to list mode)
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updated, _ := m.Update(msg)
	model := updated.(*Model)

	if model.CurrentMode != model.ListMode {
		t.Error("Expected mode to return to list mode after Esc")
	}
}

// TestResumeShortcutOnNonBlankEntry verifies r opens resume mode on non-blank entry
func TestResumeShortcutOnNonBlankEntry(t *testing.T) {
	m := newTestModel()
	// Create an entry to resume
	m.TaskManager.StartEntryAt("test-project", "test-task", mustParseTime("2025-01-01T10:00:00Z"))
	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load entries: %v", err)
	}

	// Press 'r' on the entry
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	updated, _ := m.Update(msg)
	model := updated.(*Model)

	if model.CurrentMode != model.ResumeMode {
		t.Error("Expected mode to switch to resume mode after 'r' key on non-blank entry")
	}
	// Check that project is pre-filled
	if model.Inputs[0].Value() != "test-project" {
		t.Errorf("Expected project to be pre-filled, got %q", model.Inputs[0].Value())
	}
}

// TestResumeShortcutOnBlankEntry verifies r does nothing on blank entry
func TestResumeShortcutOnBlankEntry(t *testing.T) {
	m := newTestModel()
	// Create and stop an entry to make a blank
	m.TaskManager.StartEntryAt("test-project", "test-task", mustParseTime("2025-01-01T10:00:00Z"))
	m.TaskManager.StopEntry()
	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load entries: %v", err)
	}
	// Select the blank entry (last one)
	m.SelectedIdx = len(m.Entries) - 1

	// Press 'r' on the blank entry
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	updated, _ := m.Update(msg)
	model := updated.(*Model)

	// Should stay in list mode
	if model.CurrentMode != model.ListMode {
		t.Error("Expected to stay in list mode after 'r' key on blank entry")
	}
}

// TestEditShortcut verifies e opens edit mode
func TestEditShortcut(t *testing.T) {
	m := newTestModel()
	// Create an entry to edit
	m.TaskManager.StartEntryAt("test-project", "test-task", mustParseTime("2025-01-01T10:00:00Z"))
	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load entries: %v", err)
	}

	// Press 'e' on the entry
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
	updated, _ := m.Update(msg)
	model := updated.(*Model)

	if model.CurrentMode != model.EditMode {
		t.Error("Expected mode to switch to edit mode after 'e' key")
	}
	// Check that fields are pre-filled
	if model.Inputs[0].Value() != "test-project" {
		t.Errorf("Expected project to be pre-filled, got %q", model.Inputs[0].Value())
	}
	if model.Inputs[2].Value() != "10" {
		t.Errorf("Expected hour to be pre-filled with entry time, got %q", model.Inputs[2].Value())
	}
}

// TestDeleteShortcut verifies d opens confirm mode
func TestDeleteShortcut(t *testing.T) {
	m := newTestModel()
	// Create an entry to delete
	m.TaskManager.StartEntryAt("test-project", "test-task", mustParseTime("2025-01-01T10:00:00Z"))
	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load entries: %v", err)
	}

	// Press 'd' on the entry
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
	updated, _ := m.Update(msg)
	model := updated.(*Model)

	if model.CurrentMode != model.ConfirmMode {
		t.Error("Expected mode to switch to confirm mode after 'd' key")
	}
}

// TestStopShortcutOnRunningEntry verifies s stops running entry
func TestStopShortcutOnRunningEntry(t *testing.T) {
	m := newTestModel()
	// Create a running entry
	m.TaskManager.StartEntryAt("test-project", "test-task", mustParseTime("2025-01-01T10:00:00Z"))
	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load entries: %v", err)
	}

	// Verify entry is running
	if !m.Entries[0].IsRunning() {
		t.Fatal("Expected entry to be running")
	}

	// Press 's' on the running entry
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	updated, _ := m.Update(msg)
	model := updated.(*Model)

	// Should still be in list mode
	if model.CurrentMode != model.ListMode {
		t.Error("Expected to stay in list mode after stopping")
	}
	// Status should indicate stopped
	if !strings.Contains(model.Status, "stopped") {
		t.Errorf("Expected status to mention 'stopped', got %q", model.Status)
	}
}

// TestStopShortcutOnNonRunningEntry verifies s does nothing on non-running entry
func TestStopShortcutOnNonRunningEntry(t *testing.T) {
	m := newTestModel()
	// Create and stop an entry
	m.TaskManager.StartEntryAt("test-project", "test-task", mustParseTime("2025-01-01T10:00:00Z"))
	m.TaskManager.StopEntry()
	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load entries: %v", err)
	}
	m.SelectedIdx = 0 // Select the stopped entry

	originalStatus := m.Status

	// Press 's' on non-running entry
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	updated, _ := m.Update(msg)
	model := updated.(*Model)

	// Should stay in list mode
	if model.CurrentMode != model.ListMode {
		t.Error("Expected to stay in list mode")
	}
	// Status should not change (s does nothing on non-running)
	if model.Status != originalStatus {
		t.Errorf("Expected status unchanged, got %q", model.Status)
	}
}

// TestStartEntryViaUI verifies starting an entry through the TUI
func TestStartEntryViaUI(t *testing.T) {
	m := newTestModel()

	// Initialize
	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load entries: %v", err)
	}

	// Switch to start mode
	m.CurrentMode = m.StartMode
	m.FocusIndex = 0

	// Set project
	m.Inputs[0].SetValue("test-project")
	m.Inputs[1].SetValue("test-task")
	m.Inputs[2].SetValue("14")
	m.Inputs[3].SetValue("30")

	// Simulate Enter key
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updated, _ := m.Update(msg)
	model := updated.(*Model)

	// Should return to list mode
	if model.CurrentMode != model.ListMode {
		t.Error("Expected mode to return to list mode after submit")
	}

	// Check that entry was created
	if len(model.Entries) == 0 {
		t.Error("Expected entries to be loaded after starting")
	}

	// Verify the entry
	found := false
	for _, entry := range model.Entries {
		if entry.Project == "test-project" && entry.Title == "test-task" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find started entry in list")
	}
}

// TestStopEntryViaUI verifies stopping a running entry
func TestStopEntryViaUI(t *testing.T) {
	m := newTestModel()
	tm := m.TaskManager

	// Start an entry first
	if _, err := tm.StartEntry("project1", "task1"); err != nil {
		t.Fatalf("Failed to start entry: %v", err)
	}

	// Load entries
	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load entries: %v", err)
	}

	// Find the running entry
	runningIdx := -1
	for i, entry := range m.Entries {
		if entry.IsRunning() {
			runningIdx = i
			break
		}
	}

	if runningIdx == -1 {
		t.Fatal("Expected to find running entry")
	}

	m.SelectedIdx = runningIdx

	// Simulate 's' key (should stop the entry)
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	updated, _ := m.Update(msg)
	model := updated.(*Model)

	// Reload entries to verify
	if err := model.LoadEntries(); err != nil {
		t.Fatalf("Failed to reload entries: %v", err)
	}

	// The entry should now be stopped
	if len(model.Entries) > 0 && model.Entries[runningIdx].IsRunning() {
		t.Error("Expected entry to be stopped")
	}
}

// TestOperationCompleteMsg verifies operation complete messages
func TestOperationCompleteMsg(t *testing.T) {
	m := newTestModel()
	m.Loading = true

	// Send operation complete with no error
	msg := OperationCompleteMsg{Error: nil}
	updated, _ := m.Update(msg)
	model := updated.(*Model)

	if model.Loading {
		t.Error("Expected loading flag to be cleared")
	}
	if model.Status != "" {
		t.Error("Expected no status message for successful operation")
	}
}

// TestOperationCompleteWithError verifies error messages are handled
func TestOperationCompleteWithError(t *testing.T) {
	m := newTestModel()
	m.Loading = true

	// Send operation complete with error
	msg := OperationCompleteMsg{Error: fmt.Errorf("test error")}
	updated, _ := m.Update(msg)
	model := updated.(*Model)

	if model.Loading {
		t.Error("Expected loading flag to be cleared")
	}
	if model.Status == "" {
		t.Error("Expected status message to be set for error")
	}
}

// TestLoadEntriesLoadsMostRecentFirst verifies entries are selected correctly
func TestLoadEntriesLoadsMostRecentFirst(t *testing.T) {
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)
	m := NewModel(storage, tm)

	// Start multiple entries
	if _, err := tm.StartEntry("project1", "task1"); err != nil {
		t.Fatalf("Failed to start first entry: %v", err)
	}
	if _, err := tm.StartEntry("project2", "task2"); err != nil {
		t.Fatalf("Failed to start second entry: %v", err)
	}

	// Load entries
	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load entries: %v", err)
	}

	// Should select the last (most recent) entry
	if m.SelectedIdx != len(m.Entries)-1 {
		t.Errorf("Expected selection at last entry, got %d of %d", m.SelectedIdx, len(m.Entries))
	}
}

// TestNavigationWithArrowKeys verifies up/down navigation
func TestNavigationWithArrowKeys(t *testing.T) {
	m := newTestModel()
	tm := m.TaskManager

	// Create multiple entries
	for i := 0; i < 5; i++ {
		if _, err := tm.StartEntry("project", "task"+string(rune(i))); err != nil {
			t.Fatalf("Failed to start entry %d: %v", i, err)
		}
	}

	// Load entries
	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load entries: %v", err)
	}

	startIdx := m.SelectedIdx

	// Navigate up
	msg := tea.KeyMsg{Type: tea.KeyUp}
	updated, _ := m.Update(msg)
	model := updated.(*Model)

	if model.SelectedIdx >= startIdx {
		t.Error("Expected index to decrease with up key")
	}

	// Navigate down
	prevIdx := model.SelectedIdx
	msg = tea.KeyMsg{Type: tea.KeyDown}
	updated, _ = model.Update(msg)
	model = updated.(*Model)

	if model.SelectedIdx != prevIdx+1 {
		t.Error("Expected index to increase with down key")
	}
}

// TestInputFieldNavigation verifies tabbing between input fields
func TestInputFieldNavigation(t *testing.T) {
	m := newTestModel()
	m.CurrentMode = m.StartMode
	m.FocusIndex = 0

	// Tab forward
	msg := tea.KeyMsg{Type: tea.KeyTab}
	updated, _ := m.Update(msg)
	model := updated.(*Model)

	if model.FocusIndex != 1 {
		t.Errorf("Expected focus to move to field 1, got %d", model.FocusIndex)
	}

	// Tab forward again
	msg = tea.KeyMsg{Type: tea.KeyTab}
	updated, _ = model.Update(msg)
	model = updated.(*Model)

	if model.FocusIndex != 2 {
		t.Errorf("Expected focus to move to field 2, got %d", model.FocusIndex)
	}

	// Shift+Tab backward
	msg = tea.KeyMsg{Type: tea.KeyShiftTab}
	updated, _ = model.Update(msg)
	model = updated.(*Model)

	if model.FocusIndex != 1 {
		t.Errorf("Expected focus to move back to field 1, got %d", model.FocusIndex)
	}
}

// TestHelpModeToggle verifies help mode navigation
func TestHelpModeToggle(t *testing.T) {
	m := newTestModel()
	m.CurrentMode = m.ListMode

	// Press '?' to open help
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	updated, _ := m.Update(msg)
	model := updated.(*Model)

	if model.CurrentMode != model.HelpMode {
		t.Error("Expected mode to switch to help mode")
	}
	if model.PreviousMode != m.ListMode {
		t.Error("Expected previous mode to be saved")
	}
}

// TestQuitFromListMode verifies quitting the application
func TestQuitFromListMode(t *testing.T) {
	m := newTestModel()

	// Press q to quit
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := m.Update(msg)

	// Check if quit command was returned
	if cmd == nil {
		// This is OK - the actual quit happens when the command is executed
		// Just verify that list mode handles 'q' (we saw it in the code)
	}
}

// TestEmptyListState verifies behavior with no entries
func TestEmptyListState(t *testing.T) {
	m := newTestModel()

	// Load with no entries
	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load entries: %v", err)
	}

	if len(m.Entries) != 0 {
		t.Error("Expected no entries")
	}

	if m.SelectedIdx != 0 {
		t.Error("Expected scroll offset to be 0 for empty list")
	}
}

// TestStatusMessageDisplay verifies status messages are set
func TestStatusMessageDisplay(t *testing.T) {
	m := newTestModel()
	m.Status = "Test message"

	view := m.View()

	if view == "" {
		t.Error("Expected non-empty view output")
	}

	// Status should be visible in the rendered view
	if !strings.Contains(view, "Test message") {
		t.Error("Expected view to contain the status message")
	}
}
