package utils

import (
	"sort"
	"time"
	"time-tracker/models"
)

// DailyTotal represents total time for a specific date
type DailyTotal struct {
	Date  time.Time
	Total time.Duration
}

// WeeklyTotal represents total time for a week starting on Monday
type WeeklyTotal struct {
	WeekStart time.Time
	Total     time.Duration
}

// ProjectTotal represents total time for a project
type ProjectTotal struct {
	Project string
	Total   time.Duration
}

// CalculateDailyTotals calculates total time per day for the past 7 days
func CalculateDailyTotals(entries []models.TimeEntry) []DailyTotal {
	now := time.Now()
	dailyMap := make(map[string]time.Duration)

	for _, entry := range entries {
		duration := entry.Duration()
		date := entry.Start.Format("2006-01-02")

		// Check if within past 7 days
		daysDiff := int(now.Sub(entry.Start).Hours() / 24)
		if daysDiff >= 0 && daysDiff < 7 {
			dailyMap[date] += duration
		}
	}

	var totals []DailyTotal
	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		total := dailyMap[dateStr]
		totals = append(totals, DailyTotal{Date: date, Total: total})
	}

	return totals
}

// CalculateWeeklyTotals calculates total time per week for the past 4 weeks
func CalculateWeeklyTotals(entries []models.TimeEntry) []WeeklyTotal {
	now := time.Now()
	weeklyMap := make(map[string]time.Duration)

	for _, entry := range entries {
		duration := entry.Duration()

		// Find Monday of the week
		weekStart := entry.Start.AddDate(0, 0, -int(entry.Start.Weekday()-time.Monday))
		weekStr := weekStart.Format("2006-01-02")

		// Check if within past 28 days
		daysDiff := int(now.Sub(entry.Start).Hours() / 24)
		if daysDiff >= 0 && daysDiff < 28 {
			weeklyMap[weekStr] += duration
		}
	}

	var totals []WeeklyTotal
	for i := 3; i >= 0; i-- {
		// Find Monday i weeks ago
		weekStart := now.AddDate(0, 0, -int(now.Weekday()-time.Monday)-7*i)
		weekStr := weekStart.Format("2006-01-02")
		total := weeklyMap[weekStr]
		totals = append(totals, WeeklyTotal{WeekStart: weekStart, Total: total})
	}

	return totals
}

// CalculateProjectTotals calculates total time per project for the past week
func CalculateProjectTotals(entries []models.TimeEntry) []ProjectTotal {
	now := time.Now()
	projectMap := make(map[string]time.Duration)

	for _, entry := range entries {
		duration := entry.Duration()

		// Check if within past 7 days
		daysDiff := int(now.Sub(entry.Start).Hours() / 24)
		if daysDiff >= 0 && daysDiff < 7 {
			projectMap[entry.Project] += duration
		}
	}

	var totals []ProjectTotal
	for project, total := range projectMap {
		totals = append(totals, ProjectTotal{Project: project, Total: total})
	}

	// Sort by total time descending
	sort.Slice(totals, func(i, j int) bool {
		return totals[i].Total > totals[j].Total
	})

	return totals
}
