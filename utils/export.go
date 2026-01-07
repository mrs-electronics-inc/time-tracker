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
// Assumes entries are already aggregated and filtered (e.g., via AggregateByProjectDate).
// Running and blank entries are already excluded in the aggregation step.
// Returns a TSV string with columns: Project, Date, Duration (min), Description
func ExportDailyProjects(entries []ProjectDateEntry) string {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = '\t'

	// Write header
	_ = writer.Write([]string{"Project", "Date", "Duration", "Description"})

	// Write data rows
	for _, entry := range entries {
		dateStr := entry.Date.Format("2006-01-02")
		durationMin := int64(entry.Duration.Minutes())
		description := strings.Join(entry.Tasks, ", ")

		_ = writer.Write([]string{
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
// Filters out blank entries (empty project and title) and running entries (no End time).
// Returns a TSV string with columns: Project, Task, Start, End, Duration (min)
func ExportRaw(entries []models.TimeEntry) string {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = '\t'

	// Write header
	_ = writer.Write([]string{"Project", "Task", "Start", "End", "Duration"})

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

		_ = writer.Write([]string{
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
