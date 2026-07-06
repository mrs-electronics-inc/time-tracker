package utils

import (
	"testing"
	"time"
	"time-tracker/models"
)

func TestCalculateDailyTotals(t *testing.T) {
	// Create test entries
	now := time.Now()
	entries := []models.TimeEntry{
		{
			Start:   now.AddDate(0, 0, -1),
			End:     &[]time.Time{now.AddDate(0, 0, -1).Add(2 * time.Hour)}[0],
			Project: "test",
			Title:   "task",
		},
	}

	totals := CalculateDailyTotals(entries, 14)

	if len(totals) != 14 {
		t.Errorf("Expected 14 daily totals, got %d", len(totals))
	}

	// Check yesterday's total
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
	found := false
	for _, total := range totals {
		if total.Date.Format("2006-01-02") == yesterday {
			if total.Total != 2*time.Hour {
				t.Errorf("Expected 2 hours, got %v", total.Total)
			}
			if total.Projects["test"] != 2*time.Hour {
				t.Errorf("Expected 2 hours for project 'test', got %v", total.Projects["test"])
			}
			found = true
		}
	}
	if !found {
		t.Errorf("Yesterday's total not found")
	}
}

func TestCalculateMonthlyTotals(t *testing.T) {
	now := time.Now()
	// Create an entry about 2 months ago (but same month to avoid edge cases)
	entryTime := now.AddDate(0, -2, -5)
	entries := []models.TimeEntry{
		{
			Start:   entryTime,
			End:     &[]time.Time{entryTime.Add(3 * time.Hour)}[0],
			Project: "test-project",
			Title:   "some task",
		},
	}

	totals := CalculateMonthlyTotals(entries, 6)

	if len(totals) != 6 {
		t.Errorf("Expected 6 monthly totals, got %d", len(totals))
	}

	found := false
	for _, total := range totals {
		if total.Projects["test-project"] == 3*time.Hour {
			if total.Total != 3*time.Hour {
				t.Errorf("Expected 3 hours total, got %v", total.Total)
			}
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Project 'test-project' with 3 hours not found in any month")
	}
}

func TestCalculateWeeklyTotals(t *testing.T) {
	now := time.Now()
	entries := []models.TimeEntry{
		{
			Start:   now.AddDate(0, 0, -7),
			End:     &[]time.Time{now.AddDate(0, 0, -7).Add(2 * time.Hour)}[0],
			Project: "test",
			Title:   "task",
		},
	}

	totals := CalculateWeeklyTotals(entries, 4)

	if len(totals) != 4 {
		t.Errorf("Expected 4 weekly totals, got %d", len(totals))
	}

	// Check that the project is recorded
	found := false
	for _, total := range totals {
		if total.Projects["test"] == 2*time.Hour {
			if total.Total != 2*time.Hour {
				t.Errorf("Expected 2 hours total, got %v", total.Total)
			}
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Project 'test' with 2 hours not found in any week")
	}
}
