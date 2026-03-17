package modes

import (
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"time-tracker/models"
	"time-tracker/utils"
)

func newSearchTestInputs() []textinput.Model {
	inputs := make([]textinput.Model, 4)
	for i := range inputs {
		inputs[i] = textinput.New()
	}
	return inputs
}

func TestRenderContentShowsSearchInputBarWhenSearchModeIsActive(t *testing.T) {
	m := &Model{
		Entries: []models.TimeEntry{
			{Project: "Backend", Title: "Build API"},
		},
		SelectedIdx:      0,
		SearchActive:     true,
		SearchQueryDraft: "backend",
		Width:            120,
	}

	content := ListMode.RenderContent(m, 6)

	if !strings.Contains(content, "Search: backend") {
		t.Fatalf("expected search input bar to be rendered, got:\n%s", content)
	}

	rowIndex := strings.Index(content, "Build API")
	searchIndex := strings.Index(content, "Search: backend")
	if rowIndex == -1 || searchIndex == -1 || searchIndex < rowIndex {
		t.Fatalf("expected search input bar after rows, got:\n%s", content)
	}
}

func TestEntryMatchesSearchQuery(t *testing.T) {
	tests := []struct {
		name     string
		entry    models.TimeEntry
		query    string
		expected bool
	}{
		{
			name:     "matches project case-insensitive substring",
			entry:    models.TimeEntry{Project: "ClientPortal", Title: "Daily sync"},
			query:    "portal",
			expected: true,
		},
		{
			name:     "matches title case-insensitive substring",
			entry:    models.TimeEntry{Project: "Ops", Title: "Incident Review"},
			query:    "review",
			expected: true,
		},
		{
			name:     "does not match when query missing from both fields",
			entry:    models.TimeEntry{Project: "Backend", Title: "Bugfix"},
			query:    "frontend",
			expected: false,
		},
		{
			name:     "empty query matches all entries",
			entry:    models.TimeEntry{Project: "Any", Title: "Task"},
			query:    "",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := entryMatchesSearchQuery(tt.entry, tt.query)
			if actual != tt.expected {
				t.Fatalf("entryMatchesSearchQuery() = %v, expected %v", actual, tt.expected)
			}
		})
	}
}

func TestFilterVisibleEntries(t *testing.T) {
	entries := []models.TimeEntry{
		{Project: "Backend", Title: "Build API"},
		{Project: "Frontend", Title: "Polish search bar"},
		{Project: "Backend", Title: "Review logs"},
	}

	tests := []struct {
		name                string
		query               string
		expectedSourceIndex []int
		expectedProjects    []string
	}{
		{
			name:                "empty query includes all entries",
			query:               "",
			expectedSourceIndex: []int{0, 1, 2},
			expectedProjects:    []string{"Backend", "Frontend", "Backend"},
		},
		{
			name:                "query preserves source indices for filtered subset",
			query:               "backend",
			expectedSourceIndex: []int{0, 2},
			expectedProjects:    []string{"Backend", "Backend"},
		},
		{
			name:                "query with no matches returns empty result",
			query:               "mobile",
			expectedSourceIndex: []int{},
			expectedProjects:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := filterVisibleEntries(entries, tt.query)

			if len(actual) != len(tt.expectedSourceIndex) {
				t.Fatalf("filterVisibleEntries() length = %d, expected %d", len(actual), len(tt.expectedSourceIndex))
			}

			for i, visible := range actual {
				if visible.SourceIndex != tt.expectedSourceIndex[i] {
					t.Fatalf("filterVisibleEntries()[%d].SourceIndex = %d, expected %d", i, visible.SourceIndex, tt.expectedSourceIndex[i])
				}

				if visible.Entry.Project != tt.expectedProjects[i] {
					t.Fatalf("filterVisibleEntries()[%d].Entry.Project = %q, expected %q", i, visible.Entry.Project, tt.expectedProjects[i])
				}
			}
		})
	}
}

func TestApplySearchOnEnterUpdatesAppliedQueryAndFilteredEntries(t *testing.T) {
	m := &Model{
		Entries: []models.TimeEntry{
			{Project: "Backend", Title: "Build API"},
			{Project: "Frontend", Title: "Polish search bar"},
			{Project: "Backend", Title: "Review logs"},
		},
		SearchActive:       true,
		SearchQueryDraft:   "backend",
		SearchAppliedQuery: "",
	}

	updatedModel, _ := ListMode.HandleKeyMsg(m, tea.KeyMsg{Type: tea.KeyEnter})

	if updatedModel.SearchAppliedQuery != "backend" {
		t.Fatalf("SearchAppliedQuery = %q, expected %q", updatedModel.SearchAppliedQuery, "backend")
	}

	if len(updatedModel.FilteredEntries) != 2 {
		t.Fatalf("FilteredEntries length = %d, expected %d", len(updatedModel.FilteredEntries), 2)
	}

	if updatedModel.FilteredEntries[0].SourceIndex != 0 || updatedModel.FilteredEntries[1].SourceIndex != 2 {
		t.Fatalf("FilteredEntries source indices = [%d, %d], expected [0, 2]",
			updatedModel.FilteredEntries[0].SourceIndex,
			updatedModel.FilteredEntries[1].SourceIndex,
		)
	}
}

func TestSearchInputBarRemainsVisibleAfterApplyingFilter(t *testing.T) {
	m := &Model{
		Entries: []models.TimeEntry{
			{Project: "Backend", Title: "Build API"},
			{Project: "Frontend", Title: "Polish search bar"},
			{Project: "Backend", Title: "Review logs"},
		},
		SelectedIdx:      2,
		SearchActive:     true,
		SearchQueryDraft: "backend",
		Width:            120,
	}

	updatedModel, _ := ListMode.HandleKeyMsg(m, tea.KeyMsg{Type: tea.KeyEnter})

	if !updatedModel.SearchActive {
		t.Fatal("SearchActive = false, expected true after applying filter")
	}

	content := ListMode.RenderContent(updatedModel, 6)
	if !strings.Contains(content, "Search: backend") {
		t.Fatalf("expected search input bar to remain visible after apply, got:\n%s", content)
	}
}

func TestRenderContentShowsSearchSpecificEmptyMessageWhenNoMatches(t *testing.T) {
	m := &Model{
		Entries: []models.TimeEntry{
			{Project: "Backend", Title: "Build API"},
		},
		SearchActive:       true,
		SearchAppliedQuery: "frontend",
		FilteredEntries:    []VisibleEntry{},
		Width:              120,
	}

	content := ListMode.RenderContent(m, 6)

	if !strings.Contains(content, "No matching time entries found") {
		t.Fatalf("expected search-specific empty message, got:\n%s", content)
	}

	if strings.Contains(content, "No time entries found. Press 'n' to start tracking.") {
		t.Fatalf("expected search-specific empty message to differ from no-data message, got:\n%s", content)
	}
}

func TestRenderContentUsesFilteredRowsForViewportRendering(t *testing.T) {
	m := &Model{
		Entries: []models.TimeEntry{
			{Project: "Backend", Title: "Build API"},
			{Project: "Frontend", Title: "Polish search bar"},
			{Project: "Backend", Title: "Review logs"},
		},
		SelectedIdx:        0,
		SearchAppliedQuery: "backend",
		FilteredEntries: []VisibleEntry{
			{Entry: models.TimeEntry{Project: "Backend", Title: "Build API"}, SourceIndex: 0},
			{Entry: models.TimeEntry{Project: "Backend", Title: "Review logs"}, SourceIndex: 2},
		},
		Width: 120,
	}

	content := ListMode.RenderContent(m, 4)

	if !strings.Contains(content, "Build API") {
		t.Fatalf("expected first filtered row to render, got:\n%s", content)
	}

	if !strings.Contains(content, "Review logs") {
		t.Fatalf("expected second filtered row to render, got:\n%s", content)
	}

	if strings.Contains(content, "Polish search bar") {
		t.Fatalf("expected non-matching row to be excluded from filtered render, got:\n%s", content)
	}
}

func TestApplySearchOnEnterPreservesSelectionWhenStillMatched(t *testing.T) {
	m := &Model{
		Entries: []models.TimeEntry{
			{Project: "Backend", Title: "Build API"},
			{Project: "Frontend", Title: "Polish search bar"},
			{Project: "Backend", Title: "Review logs"},
		},
		SelectedIdx:        2,
		SearchActive:       true,
		SearchQueryDraft:   "backend",
		SearchAppliedQuery: "",
	}

	updatedModel, _ := ListMode.HandleKeyMsg(m, tea.KeyMsg{Type: tea.KeyEnter})

	if updatedModel.SelectedIdx != 2 {
		t.Fatalf("SelectedIdx = %d, expected %d", updatedModel.SelectedIdx, 2)
	}
}

func TestApplySearchOnEnterSelectsLastFilteredResultWhenSelectionNotMatched(t *testing.T) {
	m := &Model{
		Entries: []models.TimeEntry{
			{Project: "Backend", Title: "Build API"},
			{Project: "Frontend", Title: "Polish search bar"},
			{Project: "Backend", Title: "Review logs"},
		},
		SelectedIdx:        1,
		SearchActive:       true,
		SearchQueryDraft:   "backend",
		SearchAppliedQuery: "",
	}

	updatedModel, _ := ListMode.HandleKeyMsg(m, tea.KeyMsg{Type: tea.KeyEnter})

	if updatedModel.SelectedIdx != 2 {
		t.Fatalf("SelectedIdx = %d, expected %d", updatedModel.SelectedIdx, 2)
	}
}

func TestEscWhileSearchActiveClearsSearchAndRestoresFullList(t *testing.T) {
	m := &Model{
		Entries: []models.TimeEntry{
			{Project: "Backend", Title: "Build API"},
			{Project: "Frontend", Title: "Polish search bar"},
			{Project: "Backend", Title: "Review logs"},
		},
		SearchActive:       true,
		SearchInputFocused: true,
		SearchQueryDraft:   "backend",
		SearchAppliedQuery: "backend",
		FilteredEntries: []VisibleEntry{
			{Entry: models.TimeEntry{Project: "Backend", Title: "Build API"}, SourceIndex: 0},
			{Entry: models.TimeEntry{Project: "Backend", Title: "Review logs"}, SourceIndex: 2},
		},
	}

	updatedModel, _ := ListMode.HandleKeyMsg(m, tea.KeyMsg{Type: tea.KeyEsc})

	if updatedModel.SearchActive {
		t.Fatal("SearchActive = true, expected false")
	}

	if updatedModel.SearchQueryDraft != "" {
		t.Fatalf("SearchQueryDraft = %q, expected empty string", updatedModel.SearchQueryDraft)
	}

	if updatedModel.SearchAppliedQuery != "" {
		t.Fatalf("SearchAppliedQuery = %q, expected empty string", updatedModel.SearchAppliedQuery)
	}

	if len(updatedModel.FilteredEntries) != len(updatedModel.Entries) {
		t.Fatalf("FilteredEntries length = %d, expected %d", len(updatedModel.FilteredEntries), len(updatedModel.Entries))
	}

	for i, visible := range updatedModel.FilteredEntries {
		if visible.SourceIndex != i {
			t.Fatalf("FilteredEntries[%d].SourceIndex = %d, expected %d", i, visible.SourceIndex, i)
		}
		if visible.Entry != updatedModel.Entries[i] {
			t.Fatalf("FilteredEntries[%d].Entry = %+v, expected %+v", i, visible.Entry, updatedModel.Entries[i])
		}
	}
}
func TestSlashInListModeActivatesSearchInput(t *testing.T) {
	m := &Model{
		SearchActive: false,
	}

	updatedModel, _ := ListMode.HandleKeyMsg(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})

	if !updatedModel.SearchActive {
		t.Fatal("SearchActive = false, expected true")
	}
	if !updatedModel.SearchInputFocused {
		t.Fatal("SearchInputFocused = false, expected true")
	}
}

func TestTypingInSearchModeUpdatesDraftQuery(t *testing.T) {
	m := &Model{
		SearchActive:       true,
		SearchInputFocused: true,
	}

	updatedModel, _ := ListMode.HandleKeyMsg(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	updatedModel, _ = ListMode.HandleKeyMsg(updatedModel, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	updatedModel, _ = ListMode.HandleKeyMsg(updatedModel, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

	if updatedModel.SearchQueryDraft != "bee" {
		t.Fatalf("SearchQueryDraft = %q, expected %q", updatedModel.SearchQueryDraft, "bee")
	}

	updatedModel, _ = ListMode.HandleKeyMsg(updatedModel, tea.KeyMsg{Type: tea.KeyBackspace})
	if updatedModel.SearchQueryDraft != "be" {
		t.Fatalf("SearchQueryDraft after backspace = %q, expected %q", updatedModel.SearchQueryDraft, "be")
	}
}

func TestEnterAppliesSearchAndExitsInputFocus(t *testing.T) {
	m := &Model{
		Entries: []models.TimeEntry{
			{Project: "Backend", Title: "Build API"},
			{Project: "Frontend", Title: "Polish search bar"},
		},
		SearchActive:       true,
		SearchInputFocused: true,
		SearchQueryDraft:   "backend",
	}

	updatedModel, _ := ListMode.HandleKeyMsg(m, tea.KeyMsg{Type: tea.KeyEnter})

	if updatedModel.SearchAppliedQuery != "backend" {
		t.Fatalf("SearchAppliedQuery = %q, expected %q", updatedModel.SearchAppliedQuery, "backend")
	}
	if updatedModel.SearchInputFocused {
		t.Fatal("SearchInputFocused = true, expected false after applying search")
	}
}

func TestListNavigationUsesFilteredEntriesWhenFilterApplied(t *testing.T) {
	m := &Model{
		Entries: []models.TimeEntry{
			{Project: "Backend", Title: "Build API"},
			{Project: "Frontend", Title: "Polish search bar"},
			{Project: "Backend", Title: "Review logs"},
		},
		SelectedIdx:        0,
		SearchAppliedQuery: "backend",
		FilteredEntries: []VisibleEntry{
			{Entry: models.TimeEntry{Project: "Backend", Title: "Build API"}, SourceIndex: 0},
			{Entry: models.TimeEntry{Project: "Backend", Title: "Review logs"}, SourceIndex: 2},
		},
	}

	updatedModel, _ := ListMode.HandleKeyMsg(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if updatedModel.SelectedIdx != 2 {
		t.Fatalf("SelectedIdx after j = %d, expected %d", updatedModel.SelectedIdx, 2)
	}

	updatedModel, _ = ListMode.HandleKeyMsg(updatedModel, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if updatedModel.SelectedIdx != 0 {
		t.Fatalf("SelectedIdx after k = %d, expected %d", updatedModel.SelectedIdx, 0)
	}

	updatedModel, _ = ListMode.HandleKeyMsg(updatedModel, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
	if updatedModel.SelectedIdx != 2 {
		t.Fatalf("SelectedIdx after G = %d, expected %d", updatedModel.SelectedIdx, 2)
	}
}

func TestIsValidSelectionFalseWhenSearchHasZeroFilteredResults(t *testing.T) {
	m := &Model{
		Entries: []models.TimeEntry{
			{Project: "Backend", Title: "Build API"},
		},
		SelectedIdx:        0,
		SearchAppliedQuery: "frontend",
		FilteredEntries:    []VisibleEntry{},
	}

	if isValidSelection(m) {
		t.Fatal("isValidSelection() = true, expected false when filtered result count is zero")
	}
}

func TestIsValidSelectionFalseWhenSelectionNotInFilteredResults(t *testing.T) {
	m := &Model{
		Entries: []models.TimeEntry{
			{Project: "Backend", Title: "Build API"},
			{Project: "Frontend", Title: "Polish search bar"},
			{Project: "Backend", Title: "Review logs"},
		},
		SelectedIdx:        1,
		SearchAppliedQuery: "backend",
		FilteredEntries: []VisibleEntry{
			{Entry: models.TimeEntry{Project: "Backend", Title: "Build API"}, SourceIndex: 0},
			{Entry: models.TimeEntry{Project: "Backend", Title: "Review logs"}, SourceIndex: 2},
		},
	}

	if isValidSelection(m) {
		t.Fatal("isValidSelection() = true, expected false when selected index is not in filtered rows")
	}
}

func TestNavigationDoesNotFallBackToUnfilteredSelectionWhenFilterApplied(t *testing.T) {
	m := &Model{
		Entries: []models.TimeEntry{
			{Project: "Backend", Title: "Build API"},
			{Project: "Frontend", Title: "Polish search bar"},
			{Project: "Backend", Title: "Review logs"},
		},
		SelectedIdx:        1,
		SearchAppliedQuery: "backend",
		FilteredEntries: []VisibleEntry{
			{Entry: models.TimeEntry{Project: "Backend", Title: "Build API"}, SourceIndex: 0},
			{Entry: models.TimeEntry{Project: "Backend", Title: "Review logs"}, SourceIndex: 2},
		},
	}

	updatedModel, _ := ListMode.HandleKeyMsg(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if updatedModel.SelectedIdx != 0 {
		t.Fatalf("SelectedIdx after j = %d, expected %d", updatedModel.SelectedIdx, 0)
	}

	updatedModel, _ = ListMode.HandleKeyMsg(updatedModel, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
	if updatedModel.SelectedIdx != 2 {
		t.Fatalf("SelectedIdx after G = %d, expected %d", updatedModel.SelectedIdx, 2)
	}
}

func TestEditAndDeleteUseUnderlyingSourceIndexWhenFilterApplied(t *testing.T) {
	m := &Model{
		Entries: []models.TimeEntry{
			{Project: "Backend", Title: "Build API"},
			{Project: "Frontend", Title: "Polish search bar"},
			{Project: "Backend", Title: "Review logs"},
		},
		SelectedIdx:        2,
		SearchAppliedQuery: "backend",
		FilteredEntries: []VisibleEntry{
			{Entry: models.TimeEntry{Project: "Backend", Title: "Build API"}, SourceIndex: 0},
			{Entry: models.TimeEntry{Project: "Backend", Title: "Review logs"}, SourceIndex: 2},
		},
		EditMode:    EditMode,
		ListMode:    ListMode,
		ConfirmMode: ConfirmMode,
		Inputs:      newSearchTestInputs(),
		Styles:      Styles{},
	}

	updatedModel, _ := ListMode.HandleKeyMsg(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	if updatedModel.FormState.EditingIdx != 2 {
		t.Fatalf("FormState.EditingIdx = %d, expected %d", updatedModel.FormState.EditingIdx, 2)
	}
	if updatedModel.Inputs[0].Value() != "Backend" || updatedModel.Inputs[1].Value() != "Review logs" {
		t.Fatalf("expected edit form to be prefilled from source index 2, got project=%q title=%q", updatedModel.Inputs[0].Value(), updatedModel.Inputs[1].Value())
	}

	updatedModel.SelectedIdx = 2
	updatedModel.ConfirmMode = ConfirmMode
	updatedModel, _ = ListMode.HandleKeyMsg(updatedModel, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	if updatedModel.ConfirmState.DeletingIdx != 2 {
		t.Fatalf("ConfirmState.DeletingIdx = %d, expected %d", updatedModel.ConfirmState.DeletingIdx, 2)
	}
}

func TestLoadEntriesReappliesFilterAndKeepsSearchState(t *testing.T) {
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)

	if _, err := tm.StartEntryAt("Backend", "Build API", time.Now().Add(-3*time.Hour)); err != nil {
		t.Fatalf("failed to create first entry: %v", err)
	}
	if _, err := tm.StartEntryAt("Frontend", "Polish search bar", time.Now().Add(-2*time.Hour)); err != nil {
		t.Fatalf("failed to create second entry: %v", err)
	}
	if _, err := tm.StartEntryAt("Backend", "Review logs", time.Now().Add(-1*time.Hour)); err != nil {
		t.Fatalf("failed to create third entry: %v", err)
	}

	m := &Model{
		Storage:            storage,
		TaskManager:        tm,
		SearchActive:       true,
		SearchQueryDraft:   "backend",
		SearchAppliedQuery: "backend",
		SelectedIdx:        999,
	}

	if err := m.LoadEntries(); err != nil {
		t.Fatalf("LoadEntries() error = %v", err)
	}

	if m.SearchAppliedQuery != "backend" {
		t.Fatalf("SearchAppliedQuery = %q, expected %q", m.SearchAppliedQuery, "backend")
	}
	if !m.SearchActive {
		t.Fatal("SearchActive = false, expected true")
	}
	if len(m.FilteredEntries) != 2 {
		t.Fatalf("FilteredEntries length = %d, expected %d", len(m.FilteredEntries), 2)
	}
	if m.SelectedIdx != m.FilteredEntries[len(m.FilteredEntries)-1].SourceIndex {
		t.Fatalf("SelectedIdx = %d, expected selection to snap to last filtered source index %d", m.SelectedIdx, m.FilteredEntries[len(m.FilteredEntries)-1].SourceIndex)
	}
}

func TestSwitchModeBackToListPreservesActiveFilter(t *testing.T) {
	m := &Model{
		Entries: []models.TimeEntry{
			{Project: "Backend", Title: "Build API"},
			{Project: "Frontend", Title: "Polish search bar"},
			{Project: "Backend", Title: "Review logs"},
		},
		SelectedIdx:        2,
		SearchActive:       true,
		SearchQueryDraft:   "backend",
		SearchAppliedQuery: "backend",
		FilteredEntries: []VisibleEntry{
			{Entry: models.TimeEntry{Project: "Backend", Title: "Build API"}, SourceIndex: 0},
			{Entry: models.TimeEntry{Project: "Backend", Title: "Review logs"}, SourceIndex: 2},
		},
		ListMode:     &Mode{Name: "list"},
		ProjectsMode: &Mode{Name: "projects"},
		CurrentMode:  &Mode{Name: "projects"},
	}

	m.SelectedIdx = 1
	m.SwitchMode(m.ListMode)

	if m.SearchAppliedQuery != "backend" {
		t.Fatalf("SearchAppliedQuery = %q, expected %q", m.SearchAppliedQuery, "backend")
	}
	if len(m.FilteredEntries) != 2 {
		t.Fatalf("FilteredEntries length = %d, expected %d", len(m.FilteredEntries), 2)
	}
	if m.SelectedIdx != 2 {
		t.Fatalf("SelectedIdx = %d, expected %d", m.SelectedIdx, 2)
	}
}
