package utils

import (
	"sort"
	"time"
	"time-tracker/models"
)

// Now is a function variable that returns the current time.
// It can be overridden in tests to provide deterministic times.
var Now = time.Now

// ProjectDateEntry represents an aggregated group of time entries for a (project, date) combination
type ProjectDateEntry struct {
	Project     string
	Date        time.Time
	Duration    time.Duration
	Tasks       []string // Deduplicated task titles
	RawDuration time.Duration
}

// AggregateByProjectDate groups time entries by (project, date) and collects task descriptions.
// It handles blank entries, running entries, and task deduplication.
// Returns a slice sorted by date (descending), then project name.
func AggregateByProjectDate(entries []models.TimeEntry) []ProjectDateEntry {
	// Map key: "YYYY-MM-DD:project"
	aggregated := make(map[string]*ProjectDateEntry)

	for _, entry := range entries {
		// Skip blank entries
		if entry.IsBlank() {
			continue
		}

		// For running entries, use current time as end time for duration calculation
		endTime := entry.End
		if entry.IsRunning() {
			now := Now()
			endTime = &now
		}

		duration := time.Duration(0)
		if endTime != nil {
			duration = endTime.Sub(entry.Start)
		}

		// Create key and date by parsing the date string to ensure consistency
		dateStr := entry.Start.Format("2006-01-02")
		key := dateStr + ":" + entry.Project

		// Parse the date string back to a time.Time to ensure the Date field
		// uses the same date as the formatted string (avoiding timezone truncation issues)
		parsedDate, _ := time.Parse("2006-01-02", dateStr)

		// Add to aggregation
		if aggregated[key] == nil {
			aggregated[key] = &ProjectDateEntry{
				Project: entry.Project,
				Date:    parsedDate,
				Tasks:   []string{},
			}
		}

		aggregated[key].RawDuration += duration

		// Add task if not blank and not duplicate
		if entry.Title != "" {
			found := false
			for _, task := range aggregated[key].Tasks {
				if task == entry.Title {
					found = true
					break
				}
			}
			if !found {
				aggregated[key].Tasks = append(aggregated[key].Tasks, entry.Title)
			}
		}
	}

	// Convert to slice and sort
	result := make([]ProjectDateEntry, 0, len(aggregated))
	for _, entry := range aggregated {
		// Set duration from raw duration
		entry.Duration = entry.RawDuration
		// Sort tasks alphabetically
		sort.Strings(entry.Tasks)
		result = append(result, *entry)
	}

	// Sort by date (ascending), then project name (ascending)
	sort.Slice(result, func(i, j int) bool {
		dateI := result[i].Date.Truncate(24 * time.Hour)
		dateJ := result[j].Date.Truncate(24 * time.Hour)
		if !dateI.Equal(dateJ) {
			return dateI.Before(dateJ)
		}
		return result[i].Project < result[j].Project
	})

	return result
}

// GetWeekSeparators returns indices in the aggregated slice where week boundaries occur.
// Returns indices where a new week starts (where a weekly separator row should be inserted before).
// Week is Monday to Sunday (Monday = start).
func GetWeekSeparators(aggregated []ProjectDateEntry) []int {
	if len(aggregated) == 0 {
		return []int{}
	}

	separators := []int{}
	currentWeekStart := getMonday(aggregated[0].Date)

	for i := 1; i < len(aggregated); i++ {
		weekStart := getMonday(aggregated[i].Date)
		if !weekStart.Equal(currentWeekStart) {
			separators = append(separators, i)
			currentWeekStart = weekStart
		}
	}

	return separators
}

// GetWeeklyTotal returns the total duration for entries in a given week
func GetWeeklyTotal(aggregated []ProjectDateEntry, weekStart time.Time) time.Duration {
	total := time.Duration(0)
	weekEnd := weekStart.AddDate(0, 0, 6).Add(24*time.Hour - time.Nanosecond)

	for _, entry := range aggregated {
		if entry.Date.After(weekStart.Add(-time.Nanosecond)) && entry.Date.Before(weekEnd.Add(time.Nanosecond)) {
			total += entry.Duration
		}
	}

	return total
}

// getMonday returns the Monday of the week containing the given date
func getMonday(date time.Time) time.Time {
	offset := (int(date.Weekday()) + 6) % 7 // Monday -> 0, Sunday -> 6
	return date.AddDate(0, 0, -offset).Truncate(24 * time.Hour)
}

// GetMondayOfWeek is the exported version of getMonday
func GetMondayOfWeek(date time.Time) time.Time {
	return getMonday(date)
}
