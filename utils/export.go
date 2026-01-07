package utils

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"
	"time"
	"time-tracker/models"
)

// ExportDailyProjects exports aggregated daily project entries as TSV format.
// Excludes running entries (entries without End time).
// Returns a TSV string with columns: Project, Date, Duration (min), Description
func ExportDailyProjects(entries []ProjectDateEntry) string {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = '\t'

	// Write header
	writer.Write([]string{"Project", "Date", "Duration (min)", "Description"})

	// Write data rows
	for _, entry := range entries {
		// Skip running entries (no End time in raw data, but aggregated entries always have duration)
		// For daily-projects, we use the aggregated entries which already exclude running entries
		// (running entries are handled in AggregateByProjectDate)

		dateStr := entry.Date.Format("2006-01-02")
		durationMin := int64(entry.Duration.Minutes())
		description := strings.Join(entry.Tasks, ", ")

		writer.Write([]string{
			entry.Project,
			dateStr,
			fmt.Sprintf("%d", durationMin),
			description,
		})
	}

	writer.Flush()
	return buf.String()
}

// ExportRaw exports raw time entries as TSV format.
// Filters out blank entries (empty project and title).
// Excludes running entries (entries without End time).
// Returns a TSV string with columns: Project, Task, Start, End, Duration (min)
func ExportRaw(entries []models.TimeEntry) string {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = '\t'

	// Write header
	writer.Write([]string{"Project", "Task", "Start", "End", "Duration (min)"})

	// Write data rows
	for _, entry := range entries {
		// Skip blank entries
		if entry.IsBlank() {
			continue
		}

		// Skip running entries (no End time)
		if entry.IsRunning() {
			continue
		}

		startStr := entry.Start.Format(time.RFC3339)
		endStr := entry.End.Format(time.RFC3339)
		durationMin := int64(entry.Duration().Minutes())

		writer.Write([]string{
			entry.Project,
			entry.Title,
			startStr,
			endStr,
			fmt.Sprintf("%d", durationMin),
		})
	}

	writer.Flush()
	return buf.String()
}
