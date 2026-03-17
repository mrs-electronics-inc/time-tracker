package cmd

import (
	"encoding/csv"
	"strings"
	"testing"
	"time"
	"time-tracker/models"
	"time-tracker/utils"
)

func parseExportTSV(t *testing.T, tsv string) [][]string {
	t.Helper()

	reader := csv.NewReader(strings.NewReader(tsv))
	reader.Comma = '\t'
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to parse TSV: %v", err)
	}

	return records
}

func TestBuildExportData_DailyProjectsFiltersByCategory(t *testing.T) {
	storage := utils.NewMemoryStorage()
	now := time.Date(2026, 3, 17, 12, 0, 0, 0, time.UTC)

	start1 := time.Date(2026, 3, 16, 9, 0, 0, 0, time.UTC)
	end1 := time.Date(2026, 3, 16, 10, 0, 0, 0, time.UTC)
	start2 := time.Date(2026, 3, 16, 10, 0, 0, 0, time.UTC)
	end2 := time.Date(2026, 3, 16, 11, 0, 0, 0, time.UTC)
	start3 := time.Date(2026, 3, 16, 11, 0, 0, 0, time.UTC)
	end3 := time.Date(2026, 3, 16, 12, 0, 0, 0, time.UTC)

	err := storage.Save([]models.TimeEntry{
		{Start: start1, End: &end1, Project: "Alpha", Title: "Build"},
		{Start: start2, End: &end2, Project: "Beta", Title: "Docs"},
		{Start: start3, End: &end3, Project: "Unknown", Title: "Research"},
	})
	if err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	err = storage.SaveProjects([]models.Project{
		{Name: "Alpha", Category: "Client"},
		{Name: "Beta", Category: "Internal"},
	})
	if err != nil {
		t.Fatalf("SaveProjects returned error: %v", err)
	}

	exported, err := buildExportData(storage, "daily-projects", "  client ", true, 7, now)
	if err != nil {
		t.Fatalf("buildExportData returned error: %v", err)
	}

	records := parseExportTSV(t, exported)
	if len(records) != 2 {
		t.Fatalf("expected 2 rows (header + 1 data), got %d", len(records))
	}

	if records[1][0] != "Alpha" {
		t.Fatalf("expected only Alpha row, got %q", records[1][0])
	}

	if records[1][2] != "Client" {
		t.Fatalf("expected Client category, got %q", records[1][2])
	}
}

func TestBuildExportData_RejectsWhitespaceCategory(t *testing.T) {
	storage := utils.NewMemoryStorage()
	now := time.Date(2026, 3, 17, 12, 0, 0, 0, time.UTC)

	_, err := buildExportData(storage, "daily-projects", "   ", true, 7, now)
	if err == nil {
		t.Fatal("expected error for whitespace-only category")
	}

	if !strings.Contains(err.Error(), "cannot be empty or whitespace") {
		t.Fatalf("expected whitespace validation error, got: %v", err)
	}
}

func TestBuildExportData_FiltersToPastDays(t *testing.T) {
	storage := utils.NewMemoryStorage()
	now := time.Date(2026, 3, 17, 12, 0, 0, 0, time.UTC)

	recentStart := now.AddDate(0, 0, -2)
	recentEnd := recentStart.Add(30 * time.Minute)
	oldStart := now.AddDate(0, 0, -8)
	oldEnd := oldStart.Add(45 * time.Minute)

	err := storage.Save([]models.TimeEntry{
		{Start: recentStart, End: &recentEnd, Project: "Alpha", Title: "Recent"},
		{Start: oldStart, End: &oldEnd, Project: "Alpha", Title: "Old"},
	})
	if err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	exported, err := buildExportData(storage, "raw", "", false, 7, now)
	if err != nil {
		t.Fatalf("buildExportData returned error: %v", err)
	}

	records := parseExportTSV(t, exported)
	if len(records) != 2 {
		t.Fatalf("expected 2 rows (header + 1 data), got %d", len(records))
	}

	if records[1][1] != "Recent" {
		t.Fatalf("expected only recent entry, got task %q", records[1][1])
	}
}

func TestBuildExportData_UsesCustomDaysWindow(t *testing.T) {
	storage := utils.NewMemoryStorage()
	now := time.Date(2026, 3, 17, 12, 0, 0, 0, time.UTC)

	recentStart := now.AddDate(0, 0, -2)
	recentEnd := recentStart.Add(30 * time.Minute)
	oldStart := now.AddDate(0, 0, -8)
	oldEnd := oldStart.Add(45 * time.Minute)

	err := storage.Save([]models.TimeEntry{
		{Start: recentStart, End: &recentEnd, Project: "Alpha", Title: "Recent"},
		{Start: oldStart, End: &oldEnd, Project: "Alpha", Title: "Old"},
	})
	if err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	exported, err := buildExportData(storage, "raw", "", false, 30, now)
	if err != nil {
		t.Fatalf("buildExportData returned error: %v", err)
	}

	records := parseExportTSV(t, exported)
	if len(records) != 3 {
		t.Fatalf("expected 3 rows (header + 2 data), got %d", len(records))
	}
}

func TestBuildExportData_RejectsNonPositiveDays(t *testing.T) {
	storage := utils.NewMemoryStorage()
	now := time.Date(2026, 3, 17, 12, 0, 0, 0, time.UTC)

	_, err := buildExportData(storage, "raw", "", false, 0, now)
	if err == nil {
		t.Fatal("expected error for non-positive days")
	}

	if !strings.Contains(err.Error(), "days must be a positive integer") {
		t.Fatalf("expected days validation error, got: %v", err)
	}
}
