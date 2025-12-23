package modes

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time-tracker/utils"
)

// StatsRow represents a single row in stats mode (either a data row or a separator)
type StatsRow struct {
	Project         string
	Date            string
	DurationMinutes int
	Tasks           []string
	IsWeekSeparator bool
	IsDaySeparator  bool
	WeekStartDate   string // For weekly separator rows
	DayTotalDate    string // For daily separator rows
}

// IsWeeklySeparator returns true if this is a weekly separator row
func (sr *StatsRow) IsWeeklySeparator() bool {
	return sr.IsWeekSeparator
}

// IsDailySeparator returns true if this is a daily separator row
func (sr *StatsRow) IsDailySeparator() bool {
	return sr.IsDaySeparator
}

// StatsRowFromEntry converts a ProjectDateEntry to a StatsRow
func StatsRowFromEntry(entry utils.ProjectDateEntry) StatsRow {
	return StatsRow{
		Project:         entry.Project,
		Date:            entry.Date.Format("2006-01-02"),
		DurationMinutes: int(entry.Duration.Minutes()),
		Tasks:           entry.Tasks,
		IsWeekSeparator: false,
	}
}

// StatsWeeklySeparatorRow creates a separator row for a week
func StatsWeeklySeparatorRow(weekStart time.Time, totalMinutes int) StatsRow {
	return StatsRow{
		IsWeekSeparator: true,
		WeekStartDate:   weekStart.Format("2006-01-02"),
		DurationMinutes: totalMinutes,
	}
}

// StatsDailySeparatorRow creates a separator row for a day
func StatsDailySeparatorRow(date time.Time, totalMinutes int) StatsRow {
	return StatsRow{
		IsDaySeparator:  true,
		DayTotalDate:    date.Format("2006-01-02"),
		DurationMinutes: totalMinutes,
	}
}

// StatsMode is the stats view mode
var StatsMode = &Mode{
	Name: "stats",
	KeyBindings: []KeyBinding{
		{Keys: "Tab", Label: "LIST", Description: "Switch mode"},
		{Keys: "k / ↑", Label: "UP", Description: "Move up"},
		{Keys: "j / ↓", Label: "DOWN", Description: "Move down"},
		{Keys: "?", Label: "HELP", Description: "Toggle help"},
		{Keys: "q / Esc", Label: "QUIT", Description: "Quit"},
	},
	HandleKeyMsg: func(m *Model, msg tea.KeyMsg) (*Model, tea.Cmd) {
		// Get the rows (which include separators)
		aggregated := utils.AggregateByProjectDate(m.Entries)
		rows := buildStatsRows(aggregated)

		switch msg.String() {
		case "?":
			m.PreviousMode = m.CurrentMode
			m.CurrentMode = m.HelpMode
			return m, nil

		case "q", "esc":
			return m, tea.Quit

		case "tab":
			m.SwitchMode(m.ListMode)
			return m, nil

		case "k", "up":
			if m.ViewportTop > 0 {
				m.ViewportTop--
			}
			m.Status = ""
			return m, nil

		case "j", "down":
			if len(rows) > 0 {
				m.ViewportTop++
			}
			m.Status = ""
			return m, nil
		}
		return m, nil
	},
	RenderContent: func(m *Model, availableHeight int) string {
		if m.Loading {
			return renderLoading()
		}

		return renderStatsContent(m, availableHeight)
	},
}

// renderStatsContent renders the stats mode content
func renderStatsContent(m *Model, availableHeight int) string {
	// Aggregate entries by project and date
	aggregated := utils.AggregateByProjectDate(m.Entries)

	if len(aggregated) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)
		msg := "No time entries found.\n"
		return emptyStyle.Render(msg)
	}

	// Convert to StatsRows with weekly separators
	rows := buildStatsRows(aggregated)

	// Calculate column widths based on content
	projectCol, dateCol, durationCol := getStatsColumnWidths(rows)

	headerHeight := 2
	contentHeight := max(availableHeight-headerHeight, 1)

	// On first render (ViewportTop == -1), start at bottom to show most recent data.
	// This must be here (not in SwitchMode) because we need row count.
	if m.ViewportTop == -1 && len(rows) > 0 {
		// Calculate ViewportTop by working backward from the last row
		// counting visual lines until we reach contentHeight
		visualLines := 0
		for i := len(rows) - 1; i >= 0; i-- {
			row := rows[i]
			if row.IsWeeklySeparator() || row.IsDailySeparator() {
				visualLines += 2 // separator + line
			} else {
				visualLines += 1 + len(row.Tasks) // data row + task lines
			}
			if visualLines >= contentHeight {
				m.ViewportTop = i
				break
			}
		}
		if m.ViewportTop == -1 {
			m.ViewportTop = 0
		}
	}

	// Clamp viewport top to valid range (at least 0, at most last row)
	if m.ViewportTop > len(rows)-1 {
		m.ViewportTop = max(len(rows)-1, 0)
	}
	if m.ViewportTop < 0 {
		m.ViewportTop = 0
	}

	header := renderStatsTableHeader(m, projectCol, dateCol, durationCol)
	content := renderStatsTableContent(m, rows, contentHeight, projectCol, dateCol, durationCol)

	return header + content
}

// buildStatsRows creates StatsRow entries with daily and weekly separators inserted
func buildStatsRows(aggregated []utils.ProjectDateEntry) []StatsRow {
	var rows []StatsRow

	// Get week boundaries
	separators := utils.GetWeekSeparators(aggregated)
	separatorSet := make(map[int]bool)
	for _, sep := range separators {
		separatorSet[sep] = true
	}

	// Build rows with daily and weekly separators
	var currentDate time.Time
	var dayTotal int

	for i, entry := range aggregated {
		// Insert week separator before this index if needed
		if separatorSet[i] && len(rows) > 0 {
			// First, add daily total for the previous day if we have one
			if !currentDate.IsZero() {
				rows = append(rows, StatsDailySeparatorRow(currentDate, dayTotal))
			}

			// Then add week separator
			prevWeekStart := utils.GetMondayOfWeek(aggregated[i-1].Date)
			weekTotal := utils.GetWeeklyTotal(aggregated, prevWeekStart)
			rows = append(rows, StatsWeeklySeparatorRow(prevWeekStart, int(weekTotal.Minutes())))

			// Reset day tracking
			currentDate = time.Time{}
			dayTotal = 0
		}

		// Check if we're moving to a new day
		if currentDate.IsZero() || !currentDate.Equal(entry.Date) {
			// If we had a previous day, add its total
			if !currentDate.IsZero() {
				rows = append(rows, StatsDailySeparatorRow(currentDate, dayTotal))
			}
			// Start tracking the new day
			currentDate = entry.Date
			dayTotal = 0
		}

		rows = append(rows, StatsRowFromEntry(entry))
		dayTotal += int(entry.Duration.Minutes())
	}

	// Add final day total if there are entries
	if !currentDate.IsZero() {
		rows = append(rows, StatsDailySeparatorRow(currentDate, dayTotal))
	}

	// Add final week separator if there are entries
	if len(aggregated) > 0 {
		lastWeekStart := utils.GetMondayOfWeek(aggregated[len(aggregated)-1].Date)
		weekTotal := utils.GetWeeklyTotal(aggregated, lastWeekStart)
		rows = append(rows, StatsWeeklySeparatorRow(lastWeekStart, int(weekTotal.Minutes())))
	}

	return rows
}

// formatDurationMinutes converts minutes to hh:mm format
func formatDurationMinutes(minutes int) string {
	hours := minutes / 60
	mins := minutes % 60
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, mins)
	}
	return fmt.Sprintf("%dm", mins)
}

// getStatsColumnWidths calculates column widths based on content
func getStatsColumnWidths(rows []StatsRow) (int, int, int) {
	projectCol := len("Project")
	dateCol := len("Date")
	durationCol := len("Duration")

	for _, row := range rows {
		if len(row.Project) > projectCol {
			projectCol = len(row.Project)
		}
		if len(row.Date) > dateCol {
			dateCol = len(row.Date)
		}
		if len(row.WeekStartDate) > dateCol {
			dateCol = len(row.WeekStartDate)
		}
		if len(row.DayTotalDate) > dateCol {
			dateCol = len(row.DayTotalDate)
		}
		durationStr := formatDurationMinutes(row.DurationMinutes)
		if len(durationStr) > durationCol {
			durationCol = len(durationStr)
		}
	}

	padding := 1
	return projectCol + padding, dateCol + padding, durationCol + padding
}

// renderStatsTableHeader renders the stats table header
func renderStatsTableHeader(m *Model, projectCol, dateCol, durationCol int) string {
	// Column layout: Project | Date | Spacer | Duration
	// Duration is right-aligned at the end

	// Calculate available width for spacer column
	fixedWidth := projectCol + dateCol + durationCol + 3 // 3 for column separators
	spacerWidth := max(m.Width-fixedWidth, 0)

	headerText := fmt.Sprintf(
		"%-*s %-*s %-*s %*s",
		projectCol, "Project",
		dateCol, "Date",
		spacerWidth, "",
		durationCol, "Duration",
	)

	output := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")).Render(headerText) + "\n"

	// Separator line
	separatorText := strings.Repeat("─", m.Width)
	output += lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")).Render(separatorText) + "\n"

	return output
}

// renderStatsTableContent renders the stats table rows with viewport scrolling
func renderStatsTableContent(m *Model, rows []StatsRow, maxHeight int, projectCol, dateCol, durationCol int) string {
	if len(rows) == 0 {
		return ""
	}

	// Calculate available width for spacer column
	fixedWidth := projectCol + dateCol + durationCol + 3 // 3 for column separators
	spacerWidth := max(m.Width-fixedWidth, 0)

	var output strings.Builder
	renderedLines := 0

	// Render visible rows, stopping when we exceed maxHeight
	for i := m.ViewportTop; i < len(rows) && renderedLines < maxHeight; i++ {
		row := rows[i]

		if row.IsWeeklySeparator() {
			// Render weekly separator with purple foreground
			separatorText := fmt.Sprintf(
				"%-*s %-*s %-*s %*s",
				projectCol, "WEEK OF",
				dateCol, row.WeekStartDate,
				spacerWidth, "",
				durationCol, formatDurationMinutes(row.DurationMinutes),
			)
			styledRow := lipgloss.NewStyle().
				Foreground(lipgloss.Color("5")).
				Bold(true).
				Render(separatorText)
			output.WriteString(styledRow + "\n")

			// Add separator line below week header
			separatorLine := strings.Repeat("─", m.Width)
			styledLine := lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Render(separatorLine)
			output.WriteString(styledLine + "\n")
			renderedLines += 2
		} else if row.IsDailySeparator() {
			// Render daily separator with blue foreground
			separatorText := fmt.Sprintf(
				"%-*s %-*s %-*s %*s",
				projectCol, "TOTAL",
				dateCol, row.DayTotalDate,
				spacerWidth, "",
				durationCol, formatDurationMinutes(row.DurationMinutes),
			)
			styledRow := lipgloss.NewStyle().
				Foreground(lipgloss.Color("4")).
				Bold(true).
				Render(separatorText)
			output.WriteString(styledRow + "\n")

			// Add separator line below day total
			separatorLine := strings.Repeat("─", m.Width)
			styledLine := lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Render(separatorLine)
			output.WriteString(styledLine + "\n")
			renderedLines += 2
		} else {
			// Render regular data row
			durationFormatted := formatDurationMinutes(row.DurationMinutes)

			// First line: project, date, spacer, duration
			firstLineText := fmt.Sprintf(
				"%-*s %-*s %-*s %*s",
				projectCol, row.Project,
				dateCol, row.Date,
				spacerWidth, "",
				durationCol, durationFormatted,
			)

			styledRow := m.Styles.Unselected.Render(firstLineText)
			output.WriteString(styledRow)
			renderedLines++

			// Add tasks as separate lines
			if len(row.Tasks) > 0 {
				for _, task := range row.Tasks {
					if renderedLines >= maxHeight {
						break
					}
					output.WriteString("\n  ")
					output.WriteString("• " + task)
					renderedLines++
				}
			}
			output.WriteString("\n")
		}
	}

	return output.String()
}
