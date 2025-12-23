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

// getStatsRowCount calculates the number of rendered rows (data rows + separators)
// This must match the row count used in buildStatsRows/renderStatsTableContent
func getStatsRowCount(aggregated []utils.ProjectDateEntry) int {
	if len(aggregated) == 0 {
		return 0
	}

	// Get week boundaries
	separators := utils.GetWeekSeparators(aggregated)
	separatorSet := make(map[int]bool)
	for _, sep := range separators {
		separatorSet[sep] = true
	}

	// Count rows (data rows + daily separators + week separators)
	rowCount := 0
	var currentDate time.Time

	for i, entry := range aggregated {
		// Count week separator before this index if needed
		if separatorSet[i] && rowCount > 0 {
			rowCount++ // week separator
		}

		// Check if we're moving to a new day
		if currentDate.IsZero() || !currentDate.Equal(entry.Date) {
			// Count daily separator for previous day
			if !currentDate.IsZero() {
				rowCount++ // day separator
			}
			currentDate = entry.Date
		}

		rowCount++ // data row
	}

	// Count final day total if there are entries
	if !currentDate.IsZero() {
		rowCount++ // final day separator
	}

	// Count final week separator if there are entries
	if len(aggregated) > 0 {
		rowCount++ // final week separator
	}

	return rowCount
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
		// Get the actual number of rendered rows (includes separators)
		aggregated := utils.AggregateByProjectDate(m.Entries)
		rowCount := getStatsRowCount(aggregated)

		switch msg.String() {
		case "?":
			m.PreviousMode = m.CurrentMode
			m.CurrentMode = m.HelpMode
			return m, nil

		case "q", "esc":
			return m, tea.Quit

		case "tab":
			m.CurrentMode = m.ListMode
			m.SelectedIdx = 0
			m.ViewportTop = 0
			m.Status = ""
			return m, nil

		case "k", "up":
			if m.SelectedIdx > 0 {
				m.SelectedIdx--
			}
			m.Status = ""
			return m, nil

		case "j", "down":
			if rowCount > 0 {
				m.SelectedIdx++
				if m.SelectedIdx >= rowCount {
					m.SelectedIdx = rowCount - 1
				}
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

	// Clamp selection to valid rows
	if m.SelectedIdx >= len(rows) {
		m.SelectedIdx = max(len(rows)-1, 0)
	}

	headerHeight := 2
	availableRowHeight := max(availableHeight-headerHeight, 1)
	m.EnsureSelectionVisible(availableRowHeight)

	header := renderStatsTableHeader(m)
	content := renderStatsTableContent(m, rows, availableRowHeight)

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

// renderStatsTableHeader renders the stats table header
func renderStatsTableHeader(m *Model) string {
	// Column layout: Project | Date | Spacer | Duration
	// Duration is right-aligned at the end

	// Get minimum column widths
	projectCol := len("Project")
	dateCol := len("Date")
	durationCol := len("Duration")

	padding := 1
	projectCol += padding
	dateCol += padding
	durationCol += padding

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
	separatorText := strings.Repeat("─", min(m.Width, m.Width))
	output += lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")).Render(separatorText) + "\n"

	return output
}

// renderStatsTableContent renders the stats table rows with viewport scrolling
func renderStatsTableContent(m *Model, rows []StatsRow, maxHeight int) string {
	if len(rows) == 0 {
		return ""
	}

	projectCol := 15
	dateCol := 12
	durationCol := 12

	padding := 1
	projectCol += padding
	dateCol += padding
	durationCol += padding

	// Calculate available width for spacer column
	fixedWidth := projectCol + dateCol + durationCol + 3 // 3 for column separators
	spacerWidth := max(m.Width-fixedWidth, 0)

	var output strings.Builder

	// Render visible rows
	endIdx := min(m.ViewportTop+maxHeight, len(rows))

	for i := m.ViewportTop; i < endIdx; i++ {
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

			// Add tasks as separate lines
			if len(row.Tasks) > 0 {
				for j, task := range row.Tasks {
					output.WriteString("\n  ")
					if j == 0 {
						output.WriteString("• " + task)
					} else {
						output.WriteString("• " + task)
					}
				}
			}
			output.WriteString("\n")
		}
	}

	return output.String()
}
