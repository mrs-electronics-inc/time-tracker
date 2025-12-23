package modes

import (
	"testing"
	"time"
	"time-tracker/models"
	"time-tracker/utils"

	"github.com/charmbracelet/lipgloss"
)

// Helper to create test entries
func createTestEntry(start time.Time, end *time.Time, project, title string) models.TimeEntry {
	return models.TimeEntry{
		Start:   start,
		End:     end,
		Project: project,
		Title:   title,
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func TestStatsRowData(t *testing.T) {
	t.Run("converts ProjectDateEntry to StatsRow", func(t *testing.T) {
		entry := utils.ProjectDateEntry{
			Project:  "TestProj",
			Date:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			Duration: 90 * time.Minute,
			Tasks:    []string{"Task 1", "Task 2"},
		}
		row := StatsRowFromEntry(entry)

		if row.Project != "TestProj" {
			t.Errorf("expected project 'TestProj', got '%s'", row.Project)
		}
		if row.Date != "2025-01-01" {
			t.Errorf("expected date '2025-01-01', got '%s'", row.Date)
		}
		if row.DurationMinutes != 90 {
			t.Errorf("expected 90 minutes, got %d", row.DurationMinutes)
		}
		if len(row.Tasks) != 2 {
			t.Errorf("expected 2 tasks, got %d", len(row.Tasks))
		}
	})

	t.Run("handles no tasks", func(t *testing.T) {
		entry := utils.ProjectDateEntry{
			Project:  "TestProj",
			Date:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			Duration: 60 * time.Minute,
			Tasks:    []string{},
		}
		row := StatsRowFromEntry(entry)

		if len(row.Tasks) != 0 {
			t.Errorf("expected 0 tasks, got %d", len(row.Tasks))
		}
	})
}

func TestStatsRowIsWeeklySeparator(t *testing.T) {
	t.Run("regular row is not separator", func(t *testing.T) {
		row := StatsRow{
			Project:         "P",
			Date:            "2025-01-01",
			DurationMinutes: 60,
		}
		if row.IsWeeklySeparator() {
			t.Errorf("expected regular row to not be separator")
		}
	})

	t.Run("weekly separator row", func(t *testing.T) {
		row := StatsRow{
			IsWeekSeparator: true,
			WeekStartDate:   "2025-01-01",
			DurationMinutes: 600,
		}
		if !row.IsWeeklySeparator() {
			t.Errorf("expected separator row to be separator")
		}
	})
}

func TestStatsRenderTableHeader(t *testing.T) {
	t.Run("renders header with correct columns", func(t *testing.T) {
		m := &Model{
			Width:  80,
			Height: 20,
		}
		header := renderStatsTableHeader(m)
		if header == "" {
			t.Errorf("expected non-empty header")
		}
		// Check for expected column names
		if !contains(header, "Project") {
			t.Errorf("header should contain 'Project'")
		}
		if !contains(header, "Date") {
			t.Errorf("header should contain 'Date'")
		}
		if !contains(header, "Duration") {
			t.Errorf("header should contain 'Duration'")
		}
	})
}

func TestStatsKeyBindings(t *testing.T) {
	t.Run("StatsMode has expected keybindings", func(t *testing.T) {
		keybindings := StatsMode.KeyBindings
		if len(keybindings) == 0 {
			t.Errorf("expected non-empty keybindings")
		}

		keyLabels := make(map[string]bool)
		for _, kb := range keybindings {
			keyLabels[kb.Label] = true
		}

		expectedLabels := []string{"UP", "DOWN", "LIST", "HELP", "QUIT"}
		for _, label := range expectedLabels {
			if !keyLabels[label] {
				t.Errorf("expected keybinding '%s'", label)
			}
		}
	})
}

func TestStatsRenderContent(t *testing.T) {
	t.Run("renders empty state", func(t *testing.T) {
		m := &Model{
			Entries:     []models.TimeEntry{},
			CurrentMode: StatsMode,
			Width:       80,
			Height:      20,
			Styles: Styles{
				Header: lipgloss.NewStyle(),
			},
		}
		content := StatsMode.RenderContent(m, 20)
		if content == "" {
			t.Errorf("expected non-empty content for empty state")
		}
	})
}

// Helper function
func contains(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
