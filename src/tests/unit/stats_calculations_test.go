package unit

import (
	"testing"
	"time"
	"time-tracker/models"
	"time-tracker/utils"
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

	totals := utils.CalculateDailyTotals(entries)

	if len(totals) != 7 {
		t.Errorf("Expected 7 daily totals, got %d", len(totals))
	}

	// Check yesterday's total
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
	found := false
	for _, total := range totals {
		if total.Date.Format("2006-01-02") == yesterday {
			if total.Total != 2*time.Hour {
				t.Errorf("Expected 2 hours, got %v", total.Total)
			}
			found = true
		}
	}
	if !found {
		t.Errorf("Yesterday's total not found")
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

	totals := utils.CalculateWeeklyTotals(entries)

	if len(totals) != 4 {
		t.Errorf("Expected 4 weekly totals, got %d", len(totals))
	}
}
