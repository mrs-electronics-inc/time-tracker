package utils

import (
	"time"
	"time-tracker/models"
)

// DailyTotal represents total time for a specific date
type DailyTotal struct {
	Date     time.Time
	Total    time.Duration
	Projects map[string]time.Duration
}

// WeeklyTotal represents total time for a week starting on Monday
type WeeklyTotal struct {
	WeekStart time.Time
	Total     time.Duration
	Projects  map[string]time.Duration
}

// CalculateDailyTotals calculates total time per day for the specified number of days
func CalculateDailyTotals(entries []models.TimeEntry, numDays int) []DailyTotal {
	now := time.Now()
	dailyMap := make(map[string]time.Duration)
	dailyProjectsMap := make(map[string]map[string]time.Duration)

	for _, entry := range entries {
		duration := entry.Duration()
		date := entry.Start.Format("2006-01-02")

		daysDiff := int(now.Sub(entry.Start).Hours() / 24)
		if daysDiff >= 0 && daysDiff < numDays {
			dailyMap[date] += duration
			if dailyProjectsMap[date] == nil {
				dailyProjectsMap[date] = make(map[string]time.Duration)
			}
			dailyProjectsMap[date][entry.Project] += duration
		}
	}

	var totals []DailyTotal
	for i := numDays - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		total := dailyMap[dateStr]
		projects := dailyProjectsMap[dateStr]
		if projects == nil {
			projects = make(map[string]time.Duration)
		}
		totals = append(totals, DailyTotal{Date: date, Total: total, Projects: projects})
	}

	return totals
}

// CalculateWeeklyTotals calculates total time per week for the specified number of weeks
func CalculateWeeklyTotals(entries []models.TimeEntry, numWeeks int) []WeeklyTotal {
	now := time.Now()
	weeklyMap := make(map[string]time.Duration)
	weeklyProjectsMap := make(map[string]map[string]time.Duration)

	for _, entry := range entries {
		duration := entry.Duration()

		// Find Monday of the week
		weekStart := entry.Start.AddDate(0, 0, -int(entry.Start.Weekday()-time.Monday))
		weekStr := weekStart.Format("2006-01-02")

		// Check if within past numWeeks weeks
		daysDiff := int(now.Sub(entry.Start).Hours() / 24)
		if daysDiff >= 0 && daysDiff < numWeeks*7 {
			weeklyMap[weekStr] += duration
			if weeklyProjectsMap[weekStr] == nil {
				weeklyProjectsMap[weekStr] = make(map[string]time.Duration)
			}
			weeklyProjectsMap[weekStr][entry.Project] += duration
		}
	}

	var totals []WeeklyTotal
	for i := numWeeks - 1; i >= 0; i-- {
		// Find Monday i weeks ago
		weekStart := now.AddDate(0, 0, -int(now.Weekday()-time.Monday)-7*i)
		weekStr := weekStart.Format("2006-01-02")
		total := weeklyMap[weekStr]
		projects := weeklyProjectsMap[weekStr]
		if projects == nil {
			projects = make(map[string]time.Duration)
		}
		totals = append(totals, WeeklyTotal{WeekStart: weekStart, Total: total, Projects: projects})
	}

	return totals
}
