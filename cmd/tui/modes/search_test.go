package modes

import (
	"testing"

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
