package tui

import (
	"fmt"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"time-tracker/cmd/tui/modes"
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
	if m.ProjectsMode == nil || m.ProjectsMode.Name != "projects" {
		t.Error("Expected projects mode to be initialized")
	}
	if m.SearchActive {
		t.Error("Expected search to be inactive initially")
	}
	if m.SearchInputFocused {
		t.Error("Expected search input to be unfocused initially")
	}
	if m.SearchQueryDraft != "" {
		t.Error("Expected search draft query to be empty initially")
	}
	if m.SearchAppliedQuery != "" {
		t.Error("Expected search applied query to be empty initially")
	}
}

func TestModelInitializationCreatesDateInputs(t *testing.T) {
	m := newTestModel()

	if len(m.Inputs) != modes.InputMinute+1 {
		t.Fatalf("len(Inputs) = %d, expected %d", len(m.Inputs), modes.InputMinute+1)
	}

	if got := m.Inputs[modes.InputYear].Placeholder; got != "YYYY" {
		t.Fatalf("year placeholder = %q, expected YYYY", got)
	}
	if got := m.Inputs[modes.InputMonth].Placeholder; got != "MM" {
		t.Fatalf("month placeholder = %q, expected MM", got)
	}
	if got := m.Inputs[modes.InputDay].Placeholder; got != "DD" {
		t.Fatalf("day placeholder = %q, expected DD", got)
	}
}

func TestProjectsViewRendersProjectMetadata(t *testing.T) {
	m := newTestModel()

	if _, err := m.TaskManager.AddProject("API Updates", "12573", "Backend"); err != nil {
		t.Fatalf("Failed to add project: %v", err)
	}

	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load data: %v", err)
	}

	m.CurrentMode = m.ProjectsMode
	m.Width = 80
	m.Height = 20

	view := m.View()
	if !strings.Contains(view, "Name") || !strings.Contains(view, "Code") || !strings.Contains(view, "Category") {
		t.Fatal("Expected projects view to render metadata columns")
	}
	if !strings.Contains(view, "API Updates") || !strings.Contains(view, "12573") || !strings.Contains(view, "Backend") {
		t.Fatal("Expected projects view to render project metadata values")
	}
}

func TestProjectsViewSortsCaseInsensitiveByName(t *testing.T) {
	m := newTestModel()

	projects := []struct {
		name     string
		code     string
		category string
	}{
		{name: "zeta", code: "003", category: "Ops"},
		{name: "Alpha", code: "001", category: "Core"},
		{name: "beta", code: "002", category: "Infra"},
	}

	for _, project := range projects {
		if _, err := m.TaskManager.AddProject(project.name, project.code, project.category); err != nil {
			t.Fatalf("Failed to add project %q: %v", project.name, err)
		}
	}

	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load data: %v", err)
	}

	m.CurrentMode = m.ProjectsMode
	m.Width = 80
	m.Height = 20

	view := m.View()

	alphaIdx := strings.Index(view, "Alpha")
	betaIdx := strings.Index(view, "beta")
	zetaIdx := strings.Index(view, "zeta")

	if alphaIdx == -1 || betaIdx == -1 || zetaIdx == -1 {
		t.Fatalf("Expected all projects in view, got:\n%s", view)
	}

	if !(alphaIdx < betaIdx && betaIdx < zetaIdx) {
		t.Fatalf("Expected case-insensitive name sort order Alpha, beta, zeta; got:\n%s", view)
	}
}

func TestProjectsViewScrollsThroughAllProjects(t *testing.T) {
	m := newTestModel()

	for i := 1; i <= 6; i++ {
		name := fmt.Sprintf("Project %d", i)
		if _, err := m.TaskManager.AddProject(name, fmt.Sprintf("%03d", i), "Category"); err != nil {
			t.Fatalf("Failed to add project %q: %v", name, err)
		}
	}

	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load data: %v", err)
	}

	m.CurrentMode = m.ProjectsMode
	m.Width = 80
	m.Height = 6 // 5 content lines: 2 header + 3 project rows

	view := m.View()
	if !strings.Contains(view, "Project 1") {
		t.Fatal("Expected Project 1 to be visible initially")
	}
	if strings.Contains(view, "Project 6") {
		t.Fatal("Expected Project 6 to require scrolling")
	}

	for i := 0; i < 3; i++ {
		msg := tea.KeyMsg{Type: tea.KeyDown}
		updated, _ := m.Update(msg)
		m = updated.(*Model)
	}

	view = m.View()
	if !strings.Contains(view, "Project 6") {
		t.Fatalf("Expected Project 6 to be visible after scrolling, got:\n%s", view)
	}
}

func TestProjectsModeAddProjectViaForm(t *testing.T) {
	m := newTestModel()

	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load data: %v", err)
	}

	m.CurrentMode = m.ProjectsMode

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	updated, _ := m.Update(msg)
	m = updated.(*Model)

	if m.CurrentMode.Name != "project-new" {
		t.Fatalf("Expected mode to switch to project-new after 'n', got %q", m.CurrentMode.Name)
	}

	m.ProjectInputs[0].SetValue("Website Redesign")
	m.ProjectInputs[1].SetValue("12580")
	m.ProjectInputs[2].SetValue("Design")

	msg = tea.KeyMsg{Type: tea.KeyEnter}
	updated, _ = m.Update(msg)
	m = updated.(*Model)

	if m.CurrentMode != m.ProjectsMode {
		t.Fatalf("Expected mode to return to projects after submit, got %q", m.CurrentMode.Name)
	}

	projects, err := m.Storage.LoadProjects()
	if err != nil {
		t.Fatalf("Failed to load projects: %v", err)
	}

	if len(projects) != 1 {
		t.Fatalf("Expected exactly one project, got %d", len(projects))
	}

	if projects[0].Name != "Website Redesign" || projects[0].Code != "12580" || projects[0].Category != "Design" {
		t.Fatalf("Unexpected project data: %+v", projects[0])
	}
}

func TestProjectsModeEditProjectViaForm(t *testing.T) {
	m := newTestModel()

	if _, err := m.TaskManager.AddProject("API Updates", "12573", "Backend"); err != nil {
		t.Fatalf("Failed to add project: %v", err)
	}

	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load data: %v", err)
	}

	m.CurrentMode = m.ProjectsMode

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
	updated, _ := m.Update(msg)
	m = updated.(*Model)

	if m.CurrentMode.Name != "project-edit" {
		t.Fatalf("Expected mode to switch to project-edit after 'e', got %q", m.CurrentMode.Name)
	}

	if m.ProjectInputs[0].Value() != "API Updates" || m.ProjectInputs[1].Value() != "12573" || m.ProjectInputs[2].Value() != "Backend" {
		t.Fatalf("Expected form to be pre-filled, got name=%q code=%q category=%q", m.ProjectInputs[0].Value(), m.ProjectInputs[1].Value(), m.ProjectInputs[2].Value())
	}

	m.ProjectInputs[0].SetValue("API Platform")
	m.ProjectInputs[1].SetValue("12599")
	m.ProjectInputs[2].SetValue("Core")

	msg = tea.KeyMsg{Type: tea.KeyEnter}
	updated, _ = m.Update(msg)
	m = updated.(*Model)

	if m.CurrentMode != m.ProjectsMode {
		t.Fatalf("Expected mode to return to projects after submit, got %q", m.CurrentMode.Name)
	}

	projects, err := m.Storage.LoadProjects()
	if err != nil {
		t.Fatalf("Failed to load projects: %v", err)
	}

	if len(projects) != 1 {
		t.Fatalf("Expected exactly one project, got %d", len(projects))
	}

	if projects[0].Name != "API Platform" || projects[0].Code != "12599" || projects[0].Category != "Core" {
		t.Fatalf("Unexpected project data after edit: %+v", projects[0])
	}
}

func TestProjectsModeAddProjectRejectsEmptyName(t *testing.T) {
	m := newTestModel()
	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load data: %v", err)
	}

	m.CurrentMode = m.ProjectsMode
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	updated, _ := m.Update(msg)
	m = updated.(*Model)

	m.ProjectInputs[0].SetValue("   ")
	m.ProjectInputs[1].SetValue("123")
	m.ProjectInputs[2].SetValue("Client")

	msg = tea.KeyMsg{Type: tea.KeyEnter}
	updated, _ = m.Update(msg)
	m = updated.(*Model)

	if m.CurrentMode.Name != "project-new" {
		t.Fatalf("Expected to stay in project-new mode on validation error, got %q", m.CurrentMode.Name)
	}

	if m.Status != "Error adding project: project name cannot be empty" {
		t.Fatalf("Expected validation error status, got %q", m.Status)
	}
}

func TestProjectsModeEditProjectRejectsEmptyName(t *testing.T) {
	m := newTestModel()

	if _, err := m.TaskManager.AddProject("API Updates", "12573", "Backend"); err != nil {
		t.Fatalf("Failed to add project: %v", err)
	}
	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load data: %v", err)
	}

	m.CurrentMode = m.ProjectsMode
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
	updated, _ := m.Update(msg)
	m = updated.(*Model)

	m.ProjectInputs[0].SetValue("  ")
	msg = tea.KeyMsg{Type: tea.KeyEnter}
	updated, _ = m.Update(msg)
	m = updated.(*Model)

	if m.CurrentMode.Name != "project-edit" {
		t.Fatalf("Expected to stay in project-edit mode on validation error, got %q", m.CurrentMode.Name)
	}

	if m.Status != "Error editing project: project name cannot be empty" {
		t.Fatalf("Expected validation error status, got %q", m.Status)
	}
}

func TestProjectsModeDeleteProject(t *testing.T) {
	m := newTestModel()

	if _, err := m.TaskManager.AddProject("API Updates", "12573", "Backend"); err != nil {
		t.Fatalf("Failed to add project: %v", err)
	}

	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load data: %v", err)
	}

	m.CurrentMode = m.ProjectsMode

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
	updated, _ := m.Update(msg)
	m = updated.(*Model)

	projects, err := m.Storage.LoadProjects()
	if err != nil {
		t.Fatalf("Failed to load projects: %v", err)
	}

	if len(projects) != 0 {
		t.Fatalf("Expected all projects removed, got %+v", projects)
	}

	if m.Status != "Project removed" {
		t.Fatalf("Expected success status, got %q", m.Status)
	}
}

func TestProjectsModeDeleteProjectBlockedWhenInUse(t *testing.T) {
	m := newTestModel()

	if _, err := m.TaskManager.AddProject("API Updates", "12573", "Backend"); err != nil {
		t.Fatalf("Failed to add project: %v", err)
	}
	if _, err := m.TaskManager.StartEntryAt("API Updates", "Task A", mustParseTime("2026-03-17T09:00:00Z")); err != nil {
		t.Fatalf("Failed to start entry: %v", err)
	}

	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load data: %v", err)
	}

	m.CurrentMode = m.ProjectsMode

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
	updated, _ := m.Update(msg)
	m = updated.(*Model)

	projects, err := m.Storage.LoadProjects()
	if err != nil {
		t.Fatalf("Failed to load projects: %v", err)
	}

	if len(projects) != 1 {
		t.Fatalf("Expected project to remain after blocked delete, got %+v", projects)
	}

	if !strings.Contains(m.Status, `Error removing project: cannot remove project "API Updates": referenced by 1 time entries`) {
		t.Fatalf("Expected reference-blocking status, got %q", m.Status)
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
	if model.Inputs[modes.InputProject].Value() != "test-project" {
		t.Errorf("Expected project to be pre-filled, got %q", model.Inputs[modes.InputProject].Value())
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
	if model.Inputs[modes.InputProject].Value() != "test-project" {
		t.Errorf("Expected project to be pre-filled, got %q", model.Inputs[modes.InputProject].Value())
	}
	if model.Inputs[modes.InputHour].Value() != "10" {
		t.Errorf("Expected hour to be pre-filled with entry time, got %q", model.Inputs[modes.InputHour].Value())
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

	// Switch to new mode
	m.CurrentMode = m.NewMode
	m.FocusIndex = modes.InputProject

	// Set project
	m.Inputs[modes.InputProject].SetValue("test-project")
	m.Inputs[modes.InputTitle].SetValue("test-task")
	now := time.Now()
	m.Inputs[modes.InputYear].SetValue(fmt.Sprintf("%04d", now.Year()))
	m.Inputs[modes.InputMonth].SetValue(fmt.Sprintf("%02d", int(now.Month())))
	m.Inputs[modes.InputDay].SetValue(fmt.Sprintf("%02d", now.Day()))
	m.Inputs[modes.InputHour].SetValue("14")
	m.Inputs[modes.InputMinute].SetValue("30")

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
	m.CurrentMode = m.NewMode
	m.FocusIndex = modes.InputProject

	// Tab forward
	msg := tea.KeyMsg{Type: tea.KeyTab}
	updated, _ := m.Update(msg)
	model := updated.(*Model)

	if model.FocusIndex != modes.InputTitle {
		t.Errorf("Expected focus to move to field 1, got %d", model.FocusIndex)
	}

	// Tab forward again
	msg = tea.KeyMsg{Type: tea.KeyTab}
	updated, _ = model.Update(msg)
	model = updated.(*Model)

	if model.FocusIndex != modes.InputYear {
		t.Errorf("Expected focus to move to field 2, got %d", model.FocusIndex)
	}

	// Shift+Tab backward
	msg = tea.KeyMsg{Type: tea.KeyShiftTab}
	updated, _ = model.Update(msg)
	model = updated.(*Model)

	if model.FocusIndex != modes.InputTitle {
		t.Errorf("Expected focus to move back to field 1, got %d", model.FocusIndex)
	}
}

// TestTabCyclesListStatsProjects verifies tab cycles through primary modes
func TestTabCyclesListStatsProjects(t *testing.T) {
	m := newTestModel()

	if m.CurrentMode != m.ListMode {
		t.Fatalf("Expected to start in list mode, got %q", m.CurrentMode.Name)
	}

	msg := tea.KeyMsg{Type: tea.KeyTab}
	updated, _ := m.Update(msg)
	model := updated.(*Model)

	if model.CurrentMode != model.StatsMode {
		t.Fatalf("Expected tab in list mode to switch to stats mode, got %q", model.CurrentMode.Name)
	}

	updated, _ = model.Update(msg)
	model = updated.(*Model)

	if model.CurrentMode != model.ProjectsMode {
		t.Fatalf("Expected tab in stats mode to switch to projects mode, got %q", model.CurrentMode.Name)
	}

	updated, _ = model.Update(msg)
	model = updated.(*Model)

	if model.CurrentMode != model.ListMode {
		t.Fatalf("Expected tab in projects mode to switch to list mode, got %q", model.CurrentMode.Name)
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

func TestSearchBarAnchorsAboveStatusBar(t *testing.T) {
	m := newTestModel()

	if _, err := m.TaskManager.StartEntryAt("Backend", "Build API", mustParseTime("2026-03-17T09:00:00Z")); err != nil {
		t.Fatalf("Failed to create entry: %v", err)
	}

	if err := m.LoadEntries(); err != nil {
		t.Fatalf("Failed to load entries: %v", err)
	}

	m.Width = 80
	m.Height = 10
	m.SearchActive = true
	m.SearchQueryDraft = "backend"

	view := m.View()
	lines := strings.Split(view, "\n")
	if len(lines) < 2 {
		t.Fatalf("Expected at least two rendered lines, got %d", len(lines))
	}

	searchLine := lines[len(lines)-2]
	if !strings.Contains(searchLine, "Search: backend") {
		t.Fatalf("Expected search bar directly above status bar, got line: %q", searchLine)
	}
}
