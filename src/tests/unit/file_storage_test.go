package unit

import (
	"testing"
	"time"

	"time-tracker/models"
	"time-tracker/utils"
)

func TestMigrateToV1(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		input    []models.TimeEntry
		expected []models.TimeEntry
	}{
		{
			name:     "empty input",
			input:    []models.TimeEntry{},
			expected: []models.TimeEntry{},
		},
		{
			name: "already sorted with gap",
			input: []models.TimeEntry{
				{ID: 1, Start: now.Add(-2 * time.Hour), End: &[]time.Time{now.Add(-1 * time.Hour)}[0], Project: "p1", Title: "t1"},
				{ID: 2, Start: now, End: nil, Project: "p2", Title: "t2"},
			},
			expected: []models.TimeEntry{
				{ID: 1, Start: now.Add(-2 * time.Hour), End: &[]time.Time{now.Add(-1 * time.Hour)}[0], Project: "p1", Title: "t1"},
				{ID: 3, Start: now.Add(-1 * time.Hour), End: &[]time.Time{now}[0], Project: "", Title: ""},
				{ID: 2, Start: now, End: nil, Project: "p2", Title: "t2"},
			},
		},
		{
			name: "out of order",
			input: []models.TimeEntry{
				{ID: 2, Start: now, End: nil, Project: "p2", Title: "t2"},
				{ID: 1, Start: now.Add(-2 * time.Hour), End: &[]time.Time{now.Add(-1 * time.Hour)}[0], Project: "p1", Title: "t1"},
			},
			expected: []models.TimeEntry{
				{ID: 1, Start: now.Add(-2 * time.Hour), End: &[]time.Time{now.Add(-1 * time.Hour)}[0], Project: "p1", Title: "t1"},
				{ID: 3, Start: now.Add(-1 * time.Hour), End: &[]time.Time{now}[0], Project: "", Title: ""},
				{ID: 2, Start: now, End: nil, Project: "p2", Title: "t2"},
			},
		},
		{
			name: "no gap",
			input: []models.TimeEntry{
				{ID: 1, Start: now.Add(-1 * time.Hour), End: &[]time.Time{now}[0], Project: "p1", Title: "t1"},
				{ID: 2, Start: now, End: nil, Project: "p2", Title: "t2"},
			},
			expected: []models.TimeEntry{
				{ID: 1, Start: now.Add(-1 * time.Hour), End: &[]time.Time{now}[0], Project: "p1", Title: "t1"},
				{ID: 2, Start: now, End: nil, Project: "p2", Title: "t2"},
			},
		},
		{
			name: "nil end",
			input: []models.TimeEntry{
				{ID: 1, Start: now.Add(-1 * time.Hour), End: nil, Project: "p1", Title: "t1"},
			},
			expected: []models.TimeEntry{
				{ID: 1, Start: now.Add(-1 * time.Hour), End: nil, Project: "p1", Title: "t1"},
			},
		},
		{
			name: "duplicate ids",
			input: []models.TimeEntry{
				{ID: 1, Start: now.Add(-2 * time.Hour), End: &[]time.Time{now.Add(-1 * time.Hour)}[0], Project: "p1", Title: "t1"},
				{ID: 1, Start: now, End: nil, Project: "p2", Title: "t2"},
			},
			expected: []models.TimeEntry{
				{ID: 1, Start: now.Add(-2 * time.Hour), End: &[]time.Time{now.Add(-1 * time.Hour)}[0], Project: "p1", Title: "t1"},
				{ID: 2, Start: now.Add(-1 * time.Hour), End: &[]time.Time{now}[0], Project: "", Title: ""},
				{ID: 1, Start: now, End: nil, Project: "p2", Title: "t2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.MigrateToV1(tt.input)
			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d entries, got %d", len(tt.expected), len(result))
			}
			for i, exp := range tt.expected {
				got := result[i]
				if got.ID != exp.ID || got.Project != exp.Project || got.Title != exp.Title {
					t.Errorf("entry %d: got %+v, expected %+v", i, got, exp)
				}
				if got.Start != exp.Start {
					t.Errorf("entry %d: start time mismatch", i)
				}
				if (got.End == nil) != (exp.End == nil) || (got.End != nil && *got.End != *exp.End) {
					t.Errorf("entry %d: end time mismatch", i)
				}
			}
		})
	}
}
