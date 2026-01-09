package utils

import (
	"testing"
	"time"
	"time-tracker/models"
)

func TestAggregateByProjectDate(t *testing.T) {
	t.Run("empty entries", func(t *testing.T) {
		result := AggregateByProjectDate([]models.TimeEntry{})
		if len(result) != 0 {
			t.Errorf("expected 0 entries, got %d", len(result))
		}
	})

	t.Run("skip blank entries", func(t *testing.T) {
		entries := []models.TimeEntry{
			{
				Start:   time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
				End:     timePtr(time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)),
				Project: "proj",
				Title:   "task",
			},
			{
				Start:   time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
				End:     timePtr(time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC)),
				Project: "", // blank entry
				Title:   "",
			},
		}
		result := AggregateByProjectDate(entries)
		if len(result) != 1 {
			t.Errorf("expected 1 entry after filtering blanks, got %d", len(result))
		}
	})

	t.Run("single entry", func(t *testing.T) {
		start := time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC)
		end := time.Date(2025, 1, 1, 10, 30, 0, 0, time.UTC)
		entries := []models.TimeEntry{
			{
				Start:   start,
				End:     &end,
				Project: "ProjectA",
				Title:   "Task 1",
			},
		}
		result := AggregateByProjectDate(entries)
		if len(result) != 1 {
			t.Fatalf("expected 1 entry, got %d", len(result))
		}
		if result[0].Project != "ProjectA" {
			t.Errorf("expected project 'ProjectA', got '%s'", result[0].Project)
		}
		if result[0].Duration != 90*time.Minute {
			t.Errorf("expected 90m duration, got %v", result[0].Duration)
		}
		if len(result[0].Tasks) != 1 || result[0].Tasks[0] != "Task 1" {
			t.Errorf("expected tasks ['Task 1'], got %v", result[0].Tasks)
		}
	})

	t.Run("multiple entries same project same date", func(t *testing.T) {
		start1 := time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC)
		end1 := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
		start2 := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
		end2 := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
		entries := []models.TimeEntry{
			{Start: start1, End: &end1, Project: "ProjectA", Title: "Task 1"},
			{Start: start2, End: &end2, Project: "ProjectA", Title: "Task 2"},
		}
		result := AggregateByProjectDate(entries)
		if len(result) != 1 {
			t.Fatalf("expected 1 aggregated entry, got %d", len(result))
		}
		if result[0].Duration != 3*time.Hour {
			t.Errorf("expected 3h duration, got %v", result[0].Duration)
		}
		if len(result[0].Tasks) != 2 {
			t.Errorf("expected 2 tasks, got %d: %v", len(result[0].Tasks), result[0].Tasks)
		}
	})

	t.Run("task deduplication", func(t *testing.T) {
		start1 := time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC)
		end1 := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
		start2 := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
		end2 := time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC)
		entries := []models.TimeEntry{
			{Start: start1, End: &end1, Project: "ProjectA", Title: "Same Task"},
			{Start: start2, End: &end2, Project: "ProjectA", Title: "Same Task"},
		}
		result := AggregateByProjectDate(entries)
		if len(result) != 1 {
			t.Fatalf("expected 1 aggregated entry, got %d", len(result))
		}
		if len(result[0].Tasks) != 1 {
			t.Errorf("expected 1 deduplicated task, got %d: %v", len(result[0].Tasks), result[0].Tasks)
		}
	})

	t.Run("different projects same date", func(t *testing.T) {
		start1 := time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC)
		end1 := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
		start2 := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
		end2 := time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC)
		entries := []models.TimeEntry{
			{Start: start1, End: &end1, Project: "ProjectA", Title: "Task A"},
			{Start: start2, End: &end2, Project: "ProjectB", Title: "Task B"},
		}
		result := AggregateByProjectDate(entries)
		if len(result) != 2 {
			t.Fatalf("expected 2 aggregated entries, got %d", len(result))
		}
		// Check sorting: same date, so should be sorted by project name (ascending order: A, B)
		if result[0].Project != "ProjectA" || result[1].Project != "ProjectB" {
			t.Errorf("expected projects [ProjectA, ProjectB] (sorted by name ascending), got [%s, %s]", result[0].Project, result[1].Project)
		}
	})

	t.Run("multiple dates sorted descending", func(t *testing.T) {
		date1 := time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC)
		date2 := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
		date3 := time.Date(2025, 1, 2, 9, 0, 0, 0, time.UTC)

		endDate1 := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
		endDate2 := time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC)
		endDate3 := time.Date(2025, 1, 2, 10, 0, 0, 0, time.UTC)

		entries := []models.TimeEntry{
			{Start: date1, End: &endDate1, Project: "PA", Title: "T1"},
			{Start: date2, End: &endDate2, Project: "PB", Title: "T2"},
			{Start: date3, End: &endDate3, Project: "PC", Title: "T3"},
		}
		result := AggregateByProjectDate(entries)
		if len(result) != 3 {
			t.Fatalf("expected 3 entries, got %d", len(result))
		}
		// Oldest date first (2025-01-01), then 2025-01-02 entries (sorted by project name)
		if !result[0].Date.Equal(date1.Truncate(24 * time.Hour)) {
			t.Errorf("expected first entry date 2025-01-01, got %v", result[0].Date)
		}
		if result[0].Project != "PA" {
			t.Errorf("expected first entry project PA, got %s", result[0].Project)
		}
		if result[1].Project != "PB" {
			t.Errorf("expected second entry project PB, got %s", result[1].Project)
		}
		if !result[2].Date.Equal(date3.Truncate(24 * time.Hour)) {
			t.Errorf("expected third entry date 2025-01-02, got %v", result[2].Date)
		}
	})

	t.Run("skip running entries", func(t *testing.T) {
		start := time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC)
		end := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)

		entries := []models.TimeEntry{
			{Start: start, End: &end, Project: "P", Title: "Completed Task"},
			{Start: start, End: nil, Project: "P", Title: "Running Task"},
		}
		result := AggregateByProjectDate(entries)
		if len(result) != 1 {
			t.Fatalf("expected 1 entry (running excluded), got %d", len(result))
		}
		// Only the completed task should remain
		if len(result[0].Tasks) != 1 || result[0].Tasks[0] != "Completed Task" {
			t.Errorf("expected only 'Completed Task', got %v", result[0].Tasks)
		}
	})

	t.Run("tasks sorted alphabetically", func(t *testing.T) {
		start1 := time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC)
		end1 := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
		start2 := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
		end2 := time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC)
		start3 := time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC)
		end3 := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

		entries := []models.TimeEntry{
			{Start: start1, End: &end1, Project: "P", Title: "Zebra"},
			{Start: start2, End: &end2, Project: "P", Title: "Apple"},
			{Start: start3, End: &end3, Project: "P", Title: "Banana"},
		}
		result := AggregateByProjectDate(entries)
		if len(result) != 1 {
			t.Fatalf("expected 1 aggregated entry, got %d", len(result))
		}
		expected := []string{"Apple", "Banana", "Zebra"}
		if len(result[0].Tasks) != 3 {
			t.Errorf("expected 3 tasks, got %d", len(result[0].Tasks))
		}
		for i, task := range result[0].Tasks {
			if task != expected[i] {
				t.Errorf("expected task[%d]=%s, got %s", i, expected[i], task)
			}
		}
	})
}

func TestGetWeekSeparators(t *testing.T) {
	t.Run("empty aggregated", func(t *testing.T) {
		result := GetWeekSeparators([]ProjectDateEntry{})
		if len(result) != 0 {
			t.Errorf("expected 0 separators, got %d", len(result))
		}
	})

	t.Run("single week", func(t *testing.T) {
		// Mon Jan 6, Tue Jan 7, Wed Jan 8, 2025
		entries := []ProjectDateEntry{
			{Date: time.Date(2025, 1, 8, 0, 0, 0, 0, time.UTC), Project: "P"},
			{Date: time.Date(2025, 1, 7, 0, 0, 0, 0, time.UTC), Project: "P"},
			{Date: time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC), Project: "P"},
		}
		result := GetWeekSeparators(entries)
		if len(result) != 0 {
			t.Errorf("expected 0 separators within same week, got %d", len(result))
		}
	})

	t.Run("multiple weeks", func(t *testing.T) {
		// Week 1: Jan 6-12 (Mon-Sun)
		// Week 2: Dec 30, 2024 - Jan 5, 2025 (Mon-Sun)
		// Dates in descending order: Jan 8, Jan 7, Jan 1, Dec 31, Dec 30
		entries := []ProjectDateEntry{
			{Date: time.Date(2025, 1, 8, 0, 0, 0, 0, time.UTC), Project: "P"},   // Week 1
			{Date: time.Date(2025, 1, 7, 0, 0, 0, 0, time.UTC), Project: "P"},   // Week 1
			{Date: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), Project: "P"},   // Week 2
			{Date: time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC), Project: "P"}, // Week 2
			{Date: time.Date(2024, 12, 30, 0, 0, 0, 0, time.UTC), Project: "P"}, // Week 2
		}
		result := GetWeekSeparators(entries)
		// Should have separator at index 2 (where week changes from Jan 8,7 to Jan 1)
		if len(result) != 1 {
			t.Errorf("expected 1 separator, got %d: %v", len(result), result)
		}
		if result[0] != 2 {
			t.Errorf("expected separator at index 2, got %v", result)
		}
	})
}

func TestGetWeeklyTotal(t *testing.T) {
	t.Run("single week total", func(t *testing.T) {
		mon := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC) // Monday
		entries := []ProjectDateEntry{
			{Date: mon.AddDate(0, 0, 0), Project: "P", Duration: 1 * time.Hour},
			{Date: mon.AddDate(0, 0, 1), Project: "P", Duration: 2 * time.Hour},
			{Date: mon.AddDate(0, 0, 2), Project: "P", Duration: 3 * time.Hour},
		}
		total := GetWeeklyTotal(entries, mon)
		if total != 6*time.Hour {
			t.Errorf("expected 6h total, got %v", total)
		}
	})

	t.Run("empty week", func(t *testing.T) {
		entries := []ProjectDateEntry{
			{Date: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), Project: "P", Duration: 1 * time.Hour},
		}
		weekStart := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC) // Different week
		total := GetWeeklyTotal(entries, weekStart)
		if total != 0 {
			t.Errorf("expected 0 total for empty week, got %v", total)
		}
	})
}

// Helper function to create a pointer to time.Time
func timePtr(t time.Time) *time.Time {
	return &t
}
