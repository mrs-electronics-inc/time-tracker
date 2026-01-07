package utils

import (
	"encoding/csv"
	"strings"
	"testing"
	"time"
	"time-tracker/models"
)

func parseRawTSV(t *testing.T, tsvString string) [][]string {
	reader := csv.NewReader(strings.NewReader(tsvString))
	reader.Comma = '\t'
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to parse TSV: %v", err)
	}
	return records
}

func TestExportDailyProjectsHeader(t *testing.T) {
	entries := []ProjectDateEntry{}
	result := ExportDailyProjects(entries)
	records := parseRawTSV(t, result)

	if len(records) < 1 {
		t.Fatalf("Expected at least header row, got %d rows", len(records))
	}

	expectedHeader := []string{"Project", "Date", "Duration (min)", "Description"}
	if !sliceEqual(records[0], expectedHeader) {
		t.Errorf("Expected header %v, got %v", expectedHeader, records[0])
	}
}

func TestExportDailyProjectsBasic(t *testing.T) {
	date := time.Date(2025, 12, 23, 0, 0, 0, 0, time.UTC)
	entries := []ProjectDateEntry{
		{
			Project:  "ProjectA",
			Date:     date,
			Duration: 90 * time.Minute,
			Tasks:    []string{"Task1", "Task2"},
		},
		{
			Project:  "ProjectB",
			Date:     date,
			Duration: 30 * time.Minute,
			Tasks:    []string{"Task3"},
		},
	}

	result := ExportDailyProjects(entries)
	records := parseRawTSV(t, result)

	if len(records) != 3 {
		t.Fatalf("Expected 3 rows (header + 2 data), got %d", len(records))
	}

	// Check first data row
	expectedRow0 := []string{"ProjectA", "2025-12-23", "90", "Task1, Task2"}
	if !sliceEqual(records[1], expectedRow0) {
		t.Errorf("Expected row %v, got %v", expectedRow0, records[1])
	}

	// Check second data row
	expectedRow1 := []string{"ProjectB", "2025-12-23", "30", "Task3"}
	if !sliceEqual(records[2], expectedRow1) {
		t.Errorf("Expected row %v, got %v", expectedRow1, records[2])
	}
}

func TestExportDailyProjectsWithSpecialCharacters(t *testing.T) {
	date := time.Date(2025, 12, 23, 0, 0, 0, 0, time.UTC)
	entries := []ProjectDateEntry{
		{
			Project:  "Project\tA", // Tab in project name
			Date:     date,
			Duration: 60 * time.Minute,
			Tasks:    []string{"Task\nWith\nNewlines", "Task\"With\"Quotes"},
		},
	}

	result := ExportDailyProjects(entries)
	records := parseRawTSV(t, result)

	if len(records) != 2 {
		t.Fatalf("Expected 2 rows (header + 1 data), got %d", len(records))
	}

	// Check that special characters are properly escaped
	if records[1][0] != "Project\tA" {
		t.Errorf("Expected project to contain tab, got %q", records[1][0])
	}

	if !strings.Contains(records[1][3], "Task\nWith\nNewlines") {
		t.Errorf("Expected description to contain newlines, got %q", records[1][3])
	}
}

func TestExportDailyProjectsEmptyTasks(t *testing.T) {
	date := time.Date(2025, 12, 23, 0, 0, 0, 0, time.UTC)
	entries := []ProjectDateEntry{
		{
			Project:  "ProjectA",
			Date:     date,
			Duration: 60 * time.Minute,
			Tasks:    []string{},
		},
	}

	result := ExportDailyProjects(entries)
	records := parseRawTSV(t, result)

	if len(records) != 2 {
		t.Fatalf("Expected 2 rows, got %d", len(records))
	}

	// Description should be empty
	if records[1][3] != "" {
		t.Errorf("Expected empty description, got %q", records[1][3])
	}
}

func TestExportRawHeader(t *testing.T) {
	entries := []models.TimeEntry{}
	result := ExportRaw(entries)
	records := parseRawTSV(t, result)

	if len(records) < 1 {
		t.Fatalf("Expected at least header row, got %d rows", len(records))
	}

	expectedHeader := []string{"Project", "Task", "Start", "End", "Duration (min)"}
	if !sliceEqual(records[0], expectedHeader) {
		t.Errorf("Expected header %v, got %v", expectedHeader, records[0])
	}
}

func TestExportRawBasic(t *testing.T) {
	start := time.Date(2025, 12, 23, 9, 0, 0, 0, time.UTC)
	end := time.Date(2025, 12, 23, 10, 30, 0, 0, time.UTC)
	entries := []models.TimeEntry{
		{
			Start:   start,
			End:     &end,
			Project: "ProjectA",
			Title:   "Task1",
		},
	}

	result := ExportRaw(entries)
	records := parseRawTSV(t, result)

	if len(records) != 2 {
		t.Fatalf("Expected 2 rows (header + 1 data), got %d", len(records))
	}

	expectedRow := []string{
		"ProjectA",
		"Task1",
		start.Format(time.RFC3339),
		end.Format(time.RFC3339),
		"90",
	}
	if !sliceEqual(records[1], expectedRow) {
		t.Errorf("Expected row %v, got %v", expectedRow, records[1])
	}
}

func TestExportRawFiltersBlankEntries(t *testing.T) {
	start := time.Date(2025, 12, 23, 9, 0, 0, 0, time.UTC)
	end := time.Date(2025, 12, 23, 10, 0, 0, 0, time.UTC)

	entries := []models.TimeEntry{
		{
			Start:   start,
			End:     &end,
			Project: "ProjectA",
			Title:   "Task1",
		},
		// Blank entry (empty project and title)
		{
			Start:   start,
			End:     &end,
			Project: "",
			Title:   "",
		},
		{
			Start:   start,
			End:     &end,
			Project: "ProjectB",
			Title:   "Task2",
		},
	}

	result := ExportRaw(entries)
	records := parseRawTSV(t, result)

	if len(records) != 3 {
		t.Fatalf("Expected 3 rows (header + 2 data, blank filtered), got %d", len(records))
	}

	if records[1][0] != "ProjectA" {
		t.Errorf("Expected first data row to be ProjectA, got %s", records[1][0])
	}
	if records[2][0] != "ProjectB" {
		t.Errorf("Expected second data row to be ProjectB, got %s", records[2][0])
	}
}

func TestExportRawExcludesRunningEntries(t *testing.T) {
	start := time.Date(2025, 12, 23, 9, 0, 0, 0, time.UTC)
	end := time.Date(2025, 12, 23, 10, 0, 0, 0, time.UTC)

	entries := []models.TimeEntry{
		{
			Start:   start,
			End:     &end,
			Project: "ProjectA",
			Title:   "Task1",
		},
		// Running entry (End == nil)
		{
			Start:   start,
			End:     nil,
			Project: "ProjectB",
			Title:   "Task2",
		},
	}

	result := ExportRaw(entries)
	records := parseRawTSV(t, result)

	if len(records) != 2 {
		t.Fatalf("Expected 2 rows (header + 1 data, running excluded), got %d", len(records))
	}

	if records[1][0] != "ProjectA" {
		t.Errorf("Expected only ProjectA, got %s", records[1][0])
	}
}

func TestExportRawWithSpecialCharacters(t *testing.T) {
	start := time.Date(2025, 12, 23, 9, 0, 0, 0, time.UTC)
	end := time.Date(2025, 12, 23, 10, 0, 0, 0, time.UTC)

	entries := []models.TimeEntry{
		{
			Start:   start,
			End:     &end,
			Project: "Project\tWith\tTabs",
			Title:   "Task\nWith\nNewlines",
		},
	}

	result := ExportRaw(entries)
	records := parseRawTSV(t, result)

	if len(records) != 2 {
		t.Fatalf("Expected 2 rows, got %d", len(records))
	}

	// Check that special characters are properly preserved after round-trip
	if records[1][0] != "Project\tWith\tTabs" {
		t.Errorf("Expected project with tabs, got %q", records[1][0])
	}

	if records[1][1] != "Task\nWith\nNewlines" {
		t.Errorf("Expected task with newlines, got %q", records[1][1])
	}
}

func TestExportRawDurationCalculation(t *testing.T) {
	start := time.Date(2025, 12, 23, 9, 0, 0, 0, time.UTC)
	end := time.Date(2025, 12, 23, 9, 45, 30, 0, time.UTC) // 45.5 minutes

	entries := []models.TimeEntry{
		{
			Start:   start,
			End:     &end,
			Project: "ProjectA",
			Title:   "Task1",
		},
	}

	result := ExportRaw(entries)
	records := parseRawTSV(t, result)

	// Duration should be 45 minutes (rounded down)
	if records[1][4] != "45" {
		t.Errorf("Expected duration 45 min, got %s", records[1][4])
	}
}

func TestExportDailyProjectsMultipleDays(t *testing.T) {
	date1 := time.Date(2025, 12, 23, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2025, 12, 24, 0, 0, 0, 0, time.UTC)

	entries := []ProjectDateEntry{
		{
			Project:  "ProjectA",
			Date:     date1,
			Duration: 60 * time.Minute,
			Tasks:    []string{"Task1"},
		},
		{
			Project:  "ProjectA",
			Date:     date2,
			Duration: 120 * time.Minute,
			Tasks:    []string{"Task2"},
		},
	}

	result := ExportDailyProjects(entries)
	records := parseRawTSV(t, result)

	if len(records) != 3 {
		t.Fatalf("Expected 3 rows, got %d", len(records))
	}

	if records[1][1] != "2025-12-23" || records[2][1] != "2025-12-24" {
		t.Errorf("Expected input order preserved: 2025-12-23 then 2025-12-24, got %s and %s", records[1][1], records[2][1])
	}
}

func TestExportRawMultipleEntries(t *testing.T) {
	start1 := time.Date(2025, 12, 23, 9, 0, 0, 0, time.UTC)
	end1 := time.Date(2025, 12, 23, 10, 0, 0, 0, time.UTC)
	start2 := time.Date(2025, 12, 23, 11, 0, 0, 0, time.UTC)
	end2 := time.Date(2025, 12, 23, 12, 0, 0, 0, time.UTC)

	entries := []models.TimeEntry{
		{
			Start:   start1,
			End:     &end1,
			Project: "ProjectA",
			Title:   "Task1",
		},
		{
			Start:   start2,
			End:     &end2,
			Project: "ProjectB",
			Title:   "Task2",
		},
	}

	result := ExportRaw(entries)
	records := parseRawTSV(t, result)

	if len(records) != 3 {
		t.Fatalf("Expected 3 rows (header + 2 data), got %d", len(records))
	}
}

// Helper function to compare string slices
func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
