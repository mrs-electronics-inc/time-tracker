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
	t4 := time.Date(2023, 1, 1, 13, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		input     []models.TimeEntry
		expected  []models.TimeEntry
		expectErr bool
	}{
		{
			name:      "empty entries",
			input:     []models.TimeEntry{},
			expected:  []models.TimeEntry{},
			expectErr: false,
		},
		{
			name: "inserts blank entry for gap",
			input: []models.TimeEntry{
				{ID: 1, Start: t1, End: &t2, Project: "p1", Title: "t1"},
				{ID: 2, Start: t3, End: nil, Project: "p2", Title: "t2"},
			},
			expected: []models.TimeEntry{
				{ID: 1, Start: t1, End: &t2, Project: "p1", Title: "t1"},
				{ID: 3, Start: t2, End: &t3, Project: "", Title: ""},
				{ID: 2, Start: t3, End: nil, Project: "p2", Title: "t2"},
			},
			expectErr: false,
		},
		{
			name: "sorts entries by start time",
			input: []models.TimeEntry{
				{ID: 2, Start: t3, End: nil, Project: "p2", Title: "t2"},
				{ID: 1, Start: t1, End: &t2, Project: "p1", Title: "t1"},
			},
			expected: []models.TimeEntry{
				{ID: 1, Start: t1, End: &t2, Project: "p1", Title: "t1"},
				{ID: 3, Start: t2, End: &t3, Project: "", Title: ""},
				{ID: 2, Start: t3, End: nil, Project: "p2", Title: "t2"},
			},
			expectErr: false,
		},
		{
			name: "no blank entry when adjacent",
			input: []models.TimeEntry{
				{ID: 1, Start: t1, End: &t3, Project: "p1", Title: "t1"},
				{ID: 2, Start: t3, End: nil, Project: "p2", Title: "t2"},
			},
			expected: []models.TimeEntry{
				{ID: 1, Start: t1, End: &t3, Project: "p1", Title: "t1"},
				{ID: 2, Start: t3, End: nil, Project: "p2", Title: "t2"},
			},
			expectErr: false,
		},
		{
			name: "handles nil end time",
			input: []models.TimeEntry{
				{ID: 1, Start: t1, End: nil, Project: "p1", Title: "t1"},
			},
			expected: []models.TimeEntry{
				{ID: 1, Start: t1, End: nil, Project: "p1", Title: "t1"},
			},
			expectErr: false,
		},
		{
			name: "handles duplicate IDs",
			input: []models.TimeEntry{
				{ID: 1, Start: t1, End: &t2, Project: "p1", Title: "t1"},
				{ID: 1, Start: t3, End: nil, Project: "p2", Title: "t2"},
			},
			expected: []models.TimeEntry{
				{ID: 1, Start: t1, End: &t2, Project: "p1", Title: "t1"},
				{ID: 2, Start: t2, End: &t3, Project: "", Title: ""},
				{ID: 1, Start: t3, End: nil, Project: "p2", Title: "t2"},
			},
			expectErr: false,
		},
		{
			name: "handles equal start times",
			input: []models.TimeEntry{
				{ID: 1, Start: t1, End: &t2, Project: "p1", Title: "t1"},
				{ID: 2, Start: t1, End: &t3, Project: "p2", Title: "t2"},
			},
			expected: []models.TimeEntry{
				{ID: 1, Start: t1, End: &t2, Project: "p1", Title: "t1"},
				{ID: 2, Start: t1, End: &t3, Project: "p2", Title: "t2"},
			},
			expectErr: false,
		},
		{
			name: "handles large ID values",
			input: []models.TimeEntry{
				{ID: 999999, Start: t1, End: &t2, Project: "p1", Title: "t1"},
				{ID: 999998, Start: t3, End: nil, Project: "p2", Title: "t2"},
			},
			expected: []models.TimeEntry{
				{ID: 999999, Start: t1, End: &t2, Project: "p1", Title: "t1"},
				{ID: 1000000, Start: t2, End: &t3, Project: "", Title: ""},
				{ID: 999998, Start: t3, End: nil, Project: "p2", Title: "t2"},
			},
			expectErr: false,
		},
		{
			name: "handles zero and negative IDs",
			input: []models.TimeEntry{
				{ID: 0, Start: t1, End: &t2, Project: "p1", Title: "t1"},
				{ID: -1, Start: t3, End: &t4, Project: "p2", Title: "t2"},
			},
			expected: []models.TimeEntry{
				{ID: 0, Start: t1, End: &t2, Project: "p1", Title: "t1"},
				{ID: 1, Start: t2, End: &t3, Project: "", Title: ""},
				{ID: -1, Start: t3, End: &t4, Project: "p2", Title: "t2"},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputJson, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("failed to marshal input: %v", err)
			}
			resultJson, err := utils.MigrateToV1(inputJson)
			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error=%v, got err=%v", tt.expectErr, err)
			}
			if tt.expectErr {
				return
			}
			var result []models.TimeEntry
			if err := json.Unmarshal(resultJson, &result); err != nil {
				t.Fatalf("failed to unmarshal result: %v", err)
			}
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

// TestMigrateToV1InputIsolation verifies that MigrateToV1 does not mutate the input slice
// and returns an independent allocation. This prevents bugs where input modifications
// accidentally affect migrated results.
func TestMigrateToV1InputIsolation(t *testing.T) {
	t1 := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC)
	t3 := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	input := []models.TimeEntry{
		{ID: 1, Start: t1, End: &t2, Project: "p1", Title: "t1"},
		{ID: 2, Start: t3, End: nil, Project: "p2", Title: "t2"},
	}

	inputJson, _ := json.Marshal(input)
	resultJson, err := utils.MigrateToV1(inputJson)
	if err != nil {
		t.Fatalf("MigrateToV1 failed: %v", err)
	}

	// Mutate the input slice after migration
	input[0].Project = "mutated"
	input[0].Title = "mutated"

	// Verify the migrated result is unaffected
	var result []models.TimeEntry
	json.Unmarshal(resultJson, &result)

	// The first entry should still have the original values (blank entry will be at index 1)
	if result[0].Project != "p1" || result[0].Title != "t1" {
		t.Errorf("mutation of input affected result: got %+v", result[0])
	}
	if result[0].ID != 1 || result[0].Start != t1 {
		t.Errorf("first result entry corrupted: got %+v", result[0])
	}
}

// TestMigrateToV1InvalidJSON verifies that MigrateToV1 returns an error when given invalid JSON
func TestMigrateToV1InvalidJSON(t *testing.T) {
	invalidJson := []byte("not valid json {")
	_, err := utils.MigrateToV1(invalidJson)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestMigrateToV2(t *testing.T) {
	t1 := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC)
	t3 := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	t2_plus_2s := time.Date(2023, 1, 1, 11, 0, 2, 0, time.UTC)

	tests := []struct {
		name      string
		input     []models.TimeEntry
		wantIDs   []int
		wantCount int
		expectErr bool
	}{
		{
			name:      "empty entries",
			input:     []models.TimeEntry{},
			wantIDs:   []int{},
			wantCount: 0,
			expectErr: false,
		},
		{
			name: "preserves all regular entries",
			input: []models.TimeEntry{
				{ID: 1, Start: t1, End: &t2, Project: "p1", Title: "t1"},
				{ID: 2, Start: t2, End: &t3, Project: "p2", Title: "t2"},
				{ID: 3, Start: t3, End: nil, Project: "p3", Title: "t3"},
			},
			wantIDs:   []int{1, 2, 3},
			wantCount: 3,
			expectErr: false,
		},
		{
			name: "preserves long blank entries",
			input: []models.TimeEntry{
				{ID: 1, Start: t1, End: &t2, Project: "p1", Title: "t1"},
				{ID: 2, Start: t2, End: &t3, Project: "", Title: ""}, // 1 hour blank - keep
				{ID: 3, Start: t3, End: nil, Project: "p3", Title: "t3"},
			},
			wantIDs:   []int{1, 2, 3},
			wantCount: 3,
			expectErr: false,
		},
		{
			name: "filters out short blank entries",
			input: []models.TimeEntry{
				{ID: 1, Start: t1, End: &t2, Project: "p1", Title: "t1"},
				{ID: 2, Start: t2, End: &t2_plus_2s, Project: "", Title: ""}, // 2 second blank - filtered
				{ID: 3, Start: t3, End: nil, Project: "p3", Title: "t3"},
			},
			wantIDs:   []int{1, 3},
			wantCount: 2,
			expectErr: false,
		},
		{
			name: "preserves non-blank entries with short duration",
			input: []models.TimeEntry{
				{ID: 1, Start: t1, End: &t2, Project: "p1", Title: "t1"},
				{ID: 2, Start: t2, End: &t2_plus_2s, Project: "p2", Title: "t2"}, // 2 second non-blank - keep
				{ID: 3, Start: t3, End: nil, Project: "p3", Title: "t3"},
			},
			wantIDs:   []int{1, 2, 3},
			wantCount: 3,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputJson, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("failed to marshal input: %v", err)
			}
			resultJson, err := utils.MigrateToV2(inputJson)
			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error=%v, got err=%v", tt.expectErr, err)
			}
			if tt.expectErr {
				return
			}
			var result []models.TimeEntry
			if err := json.Unmarshal(resultJson, &result); err != nil {
				t.Fatalf("failed to unmarshal result: %v", err)
			}
			if len(result) != tt.wantCount {
				t.Fatalf("expected %d entries, got %d", tt.wantCount, len(result))
			}
			for i, want := range tt.wantIDs {
				got := result[i]
				if got.ID != want {
					t.Errorf("entry %d: expected ID %d, got %d", i, want, got.ID)
				}
			}
		})
	}
}

func TestMigrateToV2EndTimeReconstruction(t *testing.T) {
	t1 := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC)
	t3 := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		input     []models.TimeEntry
		expectEnd []time.Time // expected End times for each entry (last entry should be zero time)
		expectErr bool
	}{
		{
			name: "sets end times from next entry's start",
			input: []models.TimeEntry{
				{ID: 1, Start: t1, End: &t2, Project: "p1", Title: "t1"},
				{ID: 2, Start: t2, End: &t3, Project: "p2", Title: "t2"},
				{ID: 3, Start: t3, End: nil, Project: "p3", Title: "t3"},
			},
			expectEnd: []time.Time{t2, t3}, // last entry doesn't need End
			expectErr: false,
		},
		{
			name: "single entry has no end time",
			input: []models.TimeEntry{
				{ID: 1, Start: t1, End: nil, Project: "p1", Title: "t1"},
			},
			expectEnd: []time.Time{}, // single entry, no End
			expectErr: false,
		},
		{
			name: "two entries set first end from second start",
			input: []models.TimeEntry{
				{ID: 1, Start: t1, End: &t2, Project: "p1", Title: "t1"},
				{ID: 2, Start: t2, End: nil, Project: "p2", Title: "t2"},
			},
			expectEnd: []time.Time{t2}, // first entry's End = second entry's Start
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputJson, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("failed to marshal input: %v", err)
			}
			resultJson, err := utils.MigrateToV2(inputJson)
			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error=%v, got err=%v", tt.expectErr, err)
			}
			if tt.expectErr {
				return
			}
			var result []models.TimeEntry
			if err := json.Unmarshal(resultJson, &result); err != nil {
				t.Fatalf("failed to unmarshal result: %v", err)
			}

			// Check end times for all entries except the last
			for i := 0; i < len(tt.expectEnd); i++ {
				if result[i].End == nil {
					t.Errorf("entry %d: expected End to be set, got nil", i)
				} else if *result[i].End != tt.expectEnd[i] {
					t.Errorf("entry %d: expected End=%v, got %v", i, tt.expectEnd[i], *result[i].End)
				}
			}

			// Check that last entry has no End
			if len(result) > 0 {
				lastIdx := len(result) - 1
				if result[lastIdx].End != nil {
					t.Errorf("last entry: expected End to be nil, got %v", *result[lastIdx].End)
				}
			}
		})
	}
}
