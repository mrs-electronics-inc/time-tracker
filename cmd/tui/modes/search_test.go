package modes

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"time-tracker/models"
)

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
