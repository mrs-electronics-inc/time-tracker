package unit

import (
	"testing"
	"time"

	"time-tracker/models"
	"time-tracker/utils"
)

func TestMigrateToV1(t *testing.T) {

	tenAM := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	elevenAM := time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC)
	noonPM := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	onePM := time.Date(2023, 1, 1, 13, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		input     []models.V0Entry
		expected  []models.V1Entry
		expectErr bool
	}{
		{
			name:      "empty entries",
			input:     []models.V0Entry{},
			expected:  nil, // nil because empty input returns nil in current implementation
			expectErr: false,
		},
		{
			name: "inserts blank entry for gap",
			input: []models.V0Entry{
				{ID: 1, Start: tenAM, End: &elevenAM, Project: "p1", Title: "tenAM"},
				{ID: 2, Start: noonPM, End: nil, Project: "p2", Title: "elevenAM"},
			},
			expected: []models.V1Entry{
				{ID: 1, Start: tenAM, End: &elevenAM, Project: "p1", Title: "tenAM"},
				{ID: 3, Start: elevenAM, End: &noonPM, Project: "", Title: ""},
				{ID: 2, Start: noonPM, End: nil, Project: "p2", Title: "elevenAM"},
			},
			expectErr: false,
		},
		{
			name: "sorts entries by start time",
			input: []models.V0Entry{
				{ID: 2, Start: noonPM, End: nil, Project: "p2", Title: "elevenAM"},
				{ID: 1, Start: tenAM, End: &elevenAM, Project: "p1", Title: "tenAM"},
			},
			expected: []models.V1Entry{
				{ID: 1, Start: tenAM, End: &elevenAM, Project: "p1", Title: "tenAM"},
				{ID: 3, Start: elevenAM, End: &noonPM, Project: "", Title: ""},
				{ID: 2, Start: noonPM, End: nil, Project: "p2", Title: "elevenAM"},
			},
			expectErr: false,
		},
		{
			name: "no blank entry when adjacent",
			input: []models.V0Entry{
				{ID: 1, Start: tenAM, End: &noonPM, Project: "p1", Title: "tenAM"},
				{ID: 2, Start: noonPM, End: nil, Project: "p2", Title: "elevenAM"},
			},
			expected: []models.V1Entry{
				{ID: 1, Start: tenAM, End: &noonPM, Project: "p1", Title: "tenAM"},
				{ID: 2, Start: noonPM, End: nil, Project: "p2", Title: "elevenAM"},
			},
			expectErr: false,
		},
		{
			name: "handles nil end time",
			input: []models.V0Entry{
				{ID: 1, Start: tenAM, End: nil, Project: "p1", Title: "tenAM"},
			},
			expected: []models.V1Entry{
				{ID: 1, Start: tenAM, End: nil, Project: "p1", Title: "tenAM"},
			},
			expectErr: false,
		},
		{
			name: "handles duplicate IDs",
			input: []models.V0Entry{
				{ID: 1, Start: tenAM, End: &elevenAM, Project: "p1", Title: "tenAM"},
				{ID: 1, Start: noonPM, End: nil, Project: "p2", Title: "elevenAM"},
			},
			expected: []models.V1Entry{
				{ID: 1, Start: tenAM, End: &elevenAM, Project: "p1", Title: "tenAM"},
				{ID: 2, Start: elevenAM, End: &noonPM, Project: "", Title: ""},
				{ID: 1, Start: noonPM, End: nil, Project: "p2", Title: "elevenAM"},
			},
			expectErr: false,
		},
		{
			name: "handles equal start times",
			input: []models.V0Entry{
				{ID: 1, Start: tenAM, End: &elevenAM, Project: "p1", Title: "tenAM"},
				{ID: 2, Start: tenAM, End: &noonPM, Project: "p2", Title: "elevenAM"},
			},
			expected: []models.V1Entry{
				{ID: 1, Start: tenAM, End: &elevenAM, Project: "p1", Title: "tenAM"},
				{ID: 2, Start: tenAM, End: &noonPM, Project: "p2", Title: "elevenAM"},
			},
			expectErr: false,
		},
		{
			name: "handles large ID values",
			input: []models.V0Entry{
				{ID: 999999, Start: tenAM, End: &elevenAM, Project: "p1", Title: "tenAM"},
				{ID: 999998, Start: noonPM, End: nil, Project: "p2", Title: "elevenAM"},
			},
			expected: []models.V1Entry{
				{ID: 999999, Start: tenAM, End: &elevenAM, Project: "p1", Title: "tenAM"},
				{ID: 1000000, Start: elevenAM, End: &noonPM, Project: "", Title: ""},
				{ID: 999998, Start: noonPM, End: nil, Project: "p2", Title: "elevenAM"},
			},
			expectErr: false,
		},
		{
			name: "handles zero and negative IDs",
			input: []models.V0Entry{
				{ID: 0, Start: tenAM, End: &elevenAM, Project: "p1", Title: "tenAM"},
				{ID: -1, Start: noonPM, End: &onePM, Project: "p2", Title: "elevenAM"},
			},
			expected: []models.V1Entry{
				{ID: 0, Start: tenAM, End: &elevenAM, Project: "p1", Title: "tenAM"},
				{ID: 1, Start: elevenAM, End: &noonPM, Project: "", Title: ""},
				{ID: -1, Start: noonPM, End: &onePM, Project: "p2", Title: "elevenAM"},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := utils.TransformV0ToV1(tt.input)
			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error=%v, got err=%v", tt.expectErr, err)
			}
			if tt.expectErr {
				return
			}

			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d entries, got %d", len(tt.expected), len(result))
			}
			for i, exp := range tt.expected {
				got := result[i]
				if got.ID != exp.ID || got.Project != exp.Project || got.Title != exp.Title {
					t.Errorf("entry %d: got %+v, expected %+v", i, got, exp)
				}
				if !got.Start.Equal(exp.Start) {
					t.Errorf("entry %d: start time mismatch", i)
				}
				if (got.End == nil) != (exp.End == nil) || (got.End != nil && !got.End.Equal(*exp.End)) {
					t.Errorf("entry %d: end time mismatch", i)
				}
			}
		})
	}
}

// TestMigrateToV1InputIsolation verifies that MigrateToV1 does not mutate the input slice
func TestMigrateToV1InputIsolation(t *testing.T) {
	t1 := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC)
	t3 := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	input := []models.V0Entry{
		{ID: 1, Start: t1, End: &t2, Project: "p1", Title: "t1"},
		{ID: 2, Start: t3, End: nil, Project: "p2", Title: "t2"},
	}

	result, err := utils.TransformV0ToV1(input)
	if err != nil {
		t.Fatalf("TransformV0ToV1 failed: %v", err)
	}

	// Mutate the input slice after migration
	input[0].Project = "mutated"
	input[0].Title = "mutated"

	// The first entry should still have the original values (blank entry will be at index 1)
	if result[0].Project != "p1" || result[0].Title != "t1" {
		t.Errorf("mutation of input affected result: got %+v", result[0])
	}
	if result[0].ID != 1 || !result[0].Start.Equal(t1) {
		t.Errorf("first result entry corrupted: got %+v", result[0])
	}
}

func TestMigrateToV2FilterBlankEntries(t *testing.T) {
	tenAM := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	elevenAM := time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC)
	noonPM := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	elevenAMPlus2s := time.Date(2023, 1, 1, 11, 0, 2, 0, time.UTC)

	tests := []struct {
		name      string
		input     []models.V1Entry
		wantIDs   []int
		wantCount int
		expectErr bool
	}{
		{
			name:      "empty entries",
			input:     []models.V1Entry{},
			wantIDs:   []int{}, // Changed to empty slice from nil/uninitialized
			wantCount: 0,
			expectErr: false,
		},
		{
			name: "preserves all regular entries",
			input: []models.V1Entry{
				{ID: 1, Start: tenAM, End: &elevenAM, Project: "p1", Title: "tenAM"},
				{ID: 2, Start: elevenAM, End: &noonPM, Project: "p2", Title: "elevenAM"},
				{ID: 3, Start: noonPM, End: nil, Project: "p3", Title: "noonPM"},
			},
			wantIDs:   []int{1, 2, 3},
			wantCount: 3,
			expectErr: false,
		},
		{
			name: "preserves long blank entries",
			input: []models.V1Entry{
				{ID: 1, Start: tenAM, End: &elevenAM, Project: "p1", Title: "tenAM"},
				{ID: 2, Start: elevenAM, End: &noonPM, Project: "", Title: ""}, // 1 hour blank - keep
				{ID: 3, Start: noonPM, End: nil, Project: "p3", Title: "noonPM"},
			},
			wantIDs:   []int{1, 2, 3},
			wantCount: 3,
			expectErr: false,
		},
		{
			name: "filters out short blank entries",
			input: []models.V1Entry{
				{ID: 1, Start: tenAM, End: &elevenAM, Project: "p1", Title: "tenAM"},
				{ID: 2, Start: elevenAM, End: &elevenAMPlus2s, Project: "", Title: ""}, // 2 second blank - filtered
				{ID: 3, Start: noonPM, End: nil, Project: "p3", Title: "noonPM"},
			},
			wantIDs:   []int{1, 3},
			wantCount: 2,
			expectErr: false,
		},
		{
			name: "preserves non-blank entries with short duration",
			input: []models.V1Entry{
				{ID: 1, Start: tenAM, End: &elevenAM, Project: "p1", Title: "tenAM"},
				{ID: 2, Start: elevenAM, End: &elevenAMPlus2s, Project: "p2", Title: "elevenAM"}, // 2 second non-blank - keep
				{ID: 3, Start: noonPM, End: nil, Project: "p3", Title: "noonPM"},
			},
			wantIDs:   []int{1, 2, 3},
			wantCount: 3,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := utils.TransformV1ToV2(tt.input)
			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error=%v, got err=%v", tt.expectErr, err)
			}
			if tt.expectErr {
				return
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

func TestMigrateToV3(t *testing.T) {
	tenAM := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	elevenAM := time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC)
	noonPM := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		input       []models.V2Entry
		expectCount int
		expectErr   bool
	}{
		{
			name:        "empty entries",
			input:       []models.V2Entry{},
			expectCount: 0,
			expectErr:   false,
		},
		{
			name: "removes ID field from entries",
			input: []models.V2Entry{
				{ID: 1, Start: tenAM, Project: "p1", Title: "tenAM"},
				{ID: 2, Start: elevenAM, Project: "p2", Title: "elevenAM"},
				{ID: 3, Start: noonPM, Project: "p3", Title: "noonPM"},
			},
			expectCount: 3,
			expectErr:   false,
		},
		{
			name: "removes ID field from blank entries",
			input: []models.V2Entry{
				{ID: 1, Start: tenAM, Project: "p1", Title: "tenAM"},
				{ID: 2, Start: elevenAM, Project: "", Title: ""},
				{ID: 3, Start: noonPM, Project: "p3", Title: "noonPM"},
			},
			expectCount: 3,
			expectErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := utils.TransformV2ToV3(tt.input)
			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error=%v, got err=%v", tt.expectErr, err)
			}
			if tt.expectErr {
				return
			}

			if len(result) != tt.expectCount {
				t.Fatalf("expected %d entries, got %d", tt.expectCount, len(result))
			}

			for i, entry := range result {
				// Check that essential fields are preserved
				// (We cannot check for ID absence explicitly because V3Entry does not have it)
				if !entry.Start.Equal(tt.input[i].Start) {
					t.Errorf("entry %d: start time changed", i)
				}
				if entry.Project != tt.input[i].Project {
					t.Errorf("entry %d: project changed", i)
				}
				if entry.Title != tt.input[i].Title {
					t.Errorf("entry %d: title changed", i)
				}
			}
		})
	}
}
