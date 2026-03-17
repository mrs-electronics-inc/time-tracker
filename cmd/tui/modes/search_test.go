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
