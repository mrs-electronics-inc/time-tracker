package unit

import (
	"encoding/json"
	"testing"
	"time"

	"time-tracker/models"
	"time-tracker/utils"
)

func TestMigrateToV1(t *testing.T) {
	t1 := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC)
	t3 := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

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
				{ID: 1, Start: t1, End: &t2, Project: "p1", Title: "t1"},
				{ID: 2, Start: t3, End: nil, Project: "p2", Title: "t2"},
			},
			expected: []models.TimeEntry{
				{ID: 1, Start: t1, End: &t2, Project: "p1", Title: "t1"},
				{ID: 3, Start: t2, End: &t3, Project: "", Title: ""},
				{ID: 2, Start: t3, End: nil, Project: "p2", Title: "t2"},
			},
		},
		{
			name: "out of order",
			input: []models.TimeEntry{
				{ID: 2, Start: t3, End: nil, Project: "p2", Title: "t2"},
				{ID: 1, Start: t1, End: &t2, Project: "p1", Title: "t1"},
			},
			expected: []models.TimeEntry{
				{ID: 1, Start: t1, End: &t2, Project: "p1", Title: "t1"},
				{ID: 3, Start: t2, End: &t3, Project: "", Title: ""},
				{ID: 2, Start: t3, End: nil, Project: "p2", Title: "t2"},
			},
		},
		{
			name: "no gap",
			input: []models.TimeEntry{
				{ID: 1, Start: t1, End: &t3, Project: "p1", Title: "t1"},
				{ID: 2, Start: t3, End: nil, Project: "p2", Title: "t2"},
			},
			expected: []models.TimeEntry{
				{ID: 1, Start: t1, End: &t3, Project: "p1", Title: "t1"},
				{ID: 2, Start: t3, End: nil, Project: "p2", Title: "t2"},
			},
		},
		{
			name: "nil end",
			input: []models.TimeEntry{
				{ID: 1, Start: t1, End: nil, Project: "p1", Title: "t1"},
			},
			expected: []models.TimeEntry{
				{ID: 1, Start: t1, End: nil, Project: "p1", Title: "t1"},
			},
		},
		{
			name: "duplicate ids",
			input: []models.TimeEntry{
				{ID: 1, Start: t1, End: &t2, Project: "p1", Title: "t1"},
				{ID: 1, Start: t3, End: nil, Project: "p2", Title: "t2"},
			},
			expected: []models.TimeEntry{
				{ID: 1, Start: t1, End: &t2, Project: "p1", Title: "t1"},
				{ID: 2, Start: t2, End: &t3, Project: "", Title: ""},
				{ID: 1, Start: t3, End: nil, Project: "p2", Title: "t2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputJson, _ := json.Marshal(tt.input)
			resultJson := utils.MigrateToV1(inputJson)
			var result []models.TimeEntry
			json.Unmarshal(resultJson, &result)
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
