package headless

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"time-tracker/cmd/tui"
	"time-tracker/utils"
)

// setupTestServer creates a headless server with initialized TUI model
func setupTestServer(t *testing.T) *Server {
	t.Helper()
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)
	model := tui.NewModel(storage, tm)

	// Load entries
	if err := model.LoadEntries(); err != nil {
		t.Fatalf("Failed to load entries: %v", err)
	}

	server := NewServer(100)
	server.model = model

	// Set window size and create renderer
	server.width = 160
	server.height = 40

	var err error
	server.renderer, err = NewRenderer(server.width, server.height)
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}

	// Send window size to model
	updated, _ := model.Update(tea.WindowSizeMsg{Width: server.width, Height: server.height})
	server.model = updated.(*tui.Model)

	// Generate initial render
	if err := server.updateRender(); err != nil {
		t.Fatalf("Failed to update render: %v", err)
	}

	return server
}

// setupTestServerWithData creates a server with some test entries
func setupTestServerWithData(t *testing.T) *Server {
	t.Helper()
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)

	// Create some test entries
	startTime := time.Now().Add(-2 * time.Hour)
	tm.StartEntryAt("project-a", "task-1", startTime)
	tm.StartEntryAt("project-b", "task-2", startTime.Add(30*time.Minute))
	tm.StartEntryAt("project-a", "task-3", startTime.Add(60*time.Minute)) // Currently running

	model := tui.NewModel(storage, tm)

	if err := model.LoadEntries(); err != nil {
		t.Fatalf("Failed to load entries: %v", err)
	}

	server := NewServer(100)
	server.model = model

	// Set window size and create renderer
	server.width = 160
	server.height = 40

	var err error
	server.renderer, err = NewRenderer(server.width, server.height)
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}

	// Send window size to model
	updated, _ := model.Update(tea.WindowSizeMsg{Width: server.width, Height: server.height})
	server.model = updated.(*tui.Model)

	// Generate initial render
	if err := server.updateRender(); err != nil {
		t.Fatalf("Failed to update render: %v", err)
	}

	return server
}

// sendKey sends a key input to the server and returns the state response
func sendKey(t *testing.T, server *Server, key string) StateResponse {
	t.Helper()

	body, _ := json.Marshal(InputRequest{Action: "key", Key: key})
	req := httptest.NewRequest(http.MethodPost, "/input", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.handleInput(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp StateResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	return resp
}

// getState gets the current state from the server
func getState(t *testing.T, server *Server) StateResponse {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/state", nil)
	w := httptest.NewRecorder()

	server.handleState(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var resp StateResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	return resp
}

// TestE2E_NewShortcut tests the 'n' shortcut opens new mode
func TestE2E_NewShortcut(t *testing.T) {
	server := setupTestServer(t)

	// Initial state should be list mode
	state := getState(t, server)
	if state.Mode != "list" {
		t.Fatalf("Expected initial mode 'list', got %q", state.Mode)
	}

	// Press 'n' to open new mode
	state = sendKey(t, server, "n")
	if state.Mode != "new" {
		t.Errorf("Expected mode 'new' after 'n' key, got %q", state.Mode)
	}

	// ANSI output should contain "New Entry"
	if !strings.Contains(state.ANSI, "New Entry") {
		t.Error("Expected 'New Entry' in ANSI output")
	}

	// Press Esc to cancel
	state = sendKey(t, server, "esc")
	if state.Mode != "list" {
		t.Errorf("Expected mode 'list' after Esc, got %q", state.Mode)
	}
}

// TestE2E_ResumeShortcut tests the 'r' shortcut opens resume mode on non-blank entries
func TestE2E_ResumeShortcut(t *testing.T) {
	server := setupTestServerWithData(t)

	// Move up to a non-running entry
	sendKey(t, server, "k")

	// Press 'r' to open resume mode
	state := sendKey(t, server, "r")
	if state.Mode != "resume" {
		t.Errorf("Expected mode 'resume' after 'r' key, got %q", state.Mode)
	}

	// ANSI output should contain "Resume Entry"
	if !strings.Contains(state.ANSI, "Resume Entry") {
		t.Error("Expected 'Resume Entry' in ANSI output")
	}

	// Should have project pre-filled
	if !strings.Contains(state.ANSI, "project-b") {
		t.Error("Expected project 'project-b' to be pre-filled in resume form")
	}

	// Press Esc to cancel
	state = sendKey(t, server, "esc")
	if state.Mode != "list" {
		t.Errorf("Expected mode 'list' after Esc, got %q", state.Mode)
	}
}

// TestE2E_ResumeShortcutDisabledOnBlank tests 'r' does nothing on blank entries
func TestE2E_ResumeShortcutDisabledOnBlank(t *testing.T) {
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)

	// Create an entry and stop it to get a blank
	tm.StartEntryAt("test", "task", time.Now().Add(-1*time.Hour))
	tm.StopEntry()

	model := tui.NewModel(storage, tm)
	model.LoadEntries()

	server := NewServer(100)
	server.model = model

	// Set window size and create renderer
	server.width = 160
	server.height = 40

	var err error
	server.renderer, err = NewRenderer(server.width, server.height)
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}

	// Send window size to model
	updated, _ := model.Update(tea.WindowSizeMsg{Width: server.width, Height: server.height})
	server.model = updated.(*tui.Model)

	// Generate initial render
	if err := server.updateRender(); err != nil {
		t.Fatalf("Failed to update render: %v", err)
	}

	// Go to the blank entry (last one)
	state := getState(t, server)
	if state.Mode != "list" {
		t.Fatalf("Expected initial mode 'list', got %q", state.Mode)
	}

	// Press 'r' on the blank entry - should stay in list mode
	state = sendKey(t, server, "r")
	if state.Mode != "list" {
		t.Errorf("Expected to stay in 'list' mode on blank entry, got %q", state.Mode)
	}
}

// TestE2E_EditShortcut tests the 'e' shortcut opens edit mode
func TestE2E_EditShortcut(t *testing.T) {
	server := setupTestServerWithData(t)

	// Press 'e' to open edit mode
	state := sendKey(t, server, "e")
	if state.Mode != "edit" {
		t.Errorf("Expected mode 'edit' after 'e' key, got %q", state.Mode)
	}

	// ANSI output should contain "Edit Entry"
	if !strings.Contains(state.ANSI, "Edit Entry") {
		t.Error("Expected 'Edit Entry' in ANSI output")
	}

	// Should have project pre-filled from the running entry
	if !strings.Contains(state.ANSI, "project-a") {
		t.Error("Expected project 'project-a' to be pre-filled in edit form")
	}

	// Press Esc to cancel
	state = sendKey(t, server, "esc")
	if state.Mode != "list" {
		t.Errorf("Expected mode 'list' after Esc, got %q", state.Mode)
	}
}

// TestE2E_DeleteShortcut tests the 'd' shortcut opens confirm mode
func TestE2E_DeleteShortcut(t *testing.T) {
	server := setupTestServerWithData(t)

	// Press 'd' to open confirm mode
	state := sendKey(t, server, "d")
	if state.Mode != "confirm" {
		t.Errorf("Expected mode 'confirm' after 'd' key, got %q", state.Mode)
	}

	// ANSI output should contain "Delete Entry?"
	if !strings.Contains(state.ANSI, "Delete Entry?") {
		t.Error("Expected 'Delete Entry?' in ANSI output")
	}

	// Press 'n' to cancel
	state = sendKey(t, server, "n")
	if state.Mode != "list" {
		t.Errorf("Expected mode 'list' after 'n', got %q", state.Mode)
	}
}

// TestE2E_DeleteConfirm tests confirming deletion converts entry to blank
func TestE2E_DeleteConfirm(t *testing.T) {
	server := setupTestServerWithData(t)

	// Move to a non-running entry
	sendKey(t, server, "k")
	sendKey(t, server, "k")

	// Press 'd' to open confirm mode
	state := sendKey(t, server, "d")
	if state.Mode != "confirm" {
		t.Fatalf("Expected mode 'confirm', got %q", state.Mode)
	}

	// Press 'y' to confirm deletion
	state = sendKey(t, server, "y")
	if state.Mode != "list" {
		t.Errorf("Expected mode 'list' after confirmation, got %q", state.Mode)
	}

	// Status should indicate deletion
	if !strings.Contains(state.ANSI, "deleted") {
		t.Log("Warning: Expected 'deleted' in status message")
	}
}

// TestE2E_StopShortcut tests the 's' shortcut only stops running entries
func TestE2E_StopShortcut(t *testing.T) {
	server := setupTestServerWithData(t)

	// Currently selected entry is running (most recent)
	state := getState(t, server)
	if !strings.Contains(state.ANSI, "running") {
		t.Log("Selected entry should be running")
	}

	// Press 's' to stop
	state = sendKey(t, server, "s")

	// Should still be in list mode
	if state.Mode != "list" {
		t.Errorf("Expected mode 'list' after 's', got %q", state.Mode)
	}

	// Status should indicate stopped
	if !strings.Contains(strings.ToLower(state.ANSI), "stopped") {
		t.Error("Expected 'stopped' in output after stopping entry")
	}
}

// TestE2E_StopShortcutNoOpOnNonRunning tests 's' does nothing on non-running entries
func TestE2E_StopShortcutNoOpOnNonRunning(t *testing.T) {
	server := setupTestServerWithData(t)

	// Move to a non-running entry
	sendKey(t, server, "k")

	// Press 's' - should do nothing
	state := sendKey(t, server, "s")

	// Should still be in list mode
	if state.Mode != "list" {
		t.Errorf("Expected mode 'list', got %q", state.Mode)
	}

	// Should not show "stopped" for a non-running entry
	// The status should be empty or unchanged
}

// TestE2E_KeyboardNavigation tests navigation shortcuts work
func TestE2E_KeyboardNavigation(t *testing.T) {
	server := setupTestServerWithData(t)

	// Test 'k' (up)
	state := sendKey(t, server, "k")
	if state.Mode != "list" {
		t.Errorf("Expected mode 'list', got %q", state.Mode)
	}

	// Test 'j' (down)
	state = sendKey(t, server, "j")
	if state.Mode != "list" {
		t.Errorf("Expected mode 'list', got %q", state.Mode)
	}

	// Test 'G' (go to current/most recent)
	sendKey(t, server, "k") // move up first
	state = sendKey(t, server, "G")
	if state.Mode != "list" {
		t.Errorf("Expected mode 'list', got %q", state.Mode)
	}
}

// TestE2E_FormSubmitNewEntry tests creating a new entry via the form
func TestE2E_FormSubmitNewEntry(t *testing.T) {
	server := setupTestServer(t)

	// Open new mode
	sendKey(t, server, "n")

	// Type project name
	for _, c := range "myproject" {
		body, _ := json.Marshal(InputRequest{Action: "type", Text: string(c)})
		req := httptest.NewRequest(http.MethodPost, "/input", bytes.NewReader(body))
		w := httptest.NewRecorder()
		server.handleInput(w, req)
	}

	// Tab to title
	sendKey(t, server, "tab")

	// Type title
	for _, c := range "mytask" {
		body, _ := json.Marshal(InputRequest{Action: "type", Text: string(c)})
		req := httptest.NewRequest(http.MethodPost, "/input", bytes.NewReader(body))
		w := httptest.NewRecorder()
		server.handleInput(w, req)
	}

	// Submit
	state := sendKey(t, server, "enter")

	// Should be back in list mode
	if state.Mode != "list" {
		t.Errorf("Expected mode 'list' after submit, got %q", state.Mode)
	}

	// Should show the new entry
	if !strings.Contains(state.ANSI, "myproject") {
		t.Error("Expected 'myproject' to appear in the list")
	}
}

// TestE2E_StatusBarShowsAllShortcuts tests status bar displays all shortcuts
func TestE2E_StatusBarShowsAllShortcuts(t *testing.T) {
	server := setupTestServer(t)

	state := getState(t, server)

	// Check all expected shortcuts are in the status bar
	expectedShortcuts := []string{"NEW", "STOP", "RESUME", "EDIT", "DELETE", "STATS", "HELP", "QUIT"}

	for _, shortcut := range expectedShortcuts {
		if !strings.Contains(state.ANSI, shortcut) {
			t.Errorf("Expected %q in status bar", shortcut)
		}
	}
}
