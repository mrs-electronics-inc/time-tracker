package modes

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time-tracker/utils"
)

// StatsRow represents a single row in stats mode (either a data row or a weekly separator)
type StatsRow struct {
	Project         string
	Date            string
	DurationMinutes int
	Tasks           []string
	IsWeekSeparator bool
	WeekStartDate   string // For weekly separator rows
}

// IsWeeklySeparator returns true if this is a weekly separator row
func (sr *StatsRow) IsWeeklySeparator() bool {
	return sr.IsWeekSeparator
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
			// For stats mode, we'll have aggregated rows + separators
			// We need to load the actual row count from rendered data
			// For now, assume m.Entries represents the data
			if len(m.Entries) > 0 {
				m.SelectedIdx++
				if m.SelectedIdx >= len(m.Entries) {
					m.SelectedIdx = len(m.Entries) - 1
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

	headerHeight := 2
	availableRowHeight := max(availableHeight-headerHeight, 1)
	m.EnsureSelectionVisible(availableRowHeight)

	header := renderStatsTableHeader(m.Width)
	content := renderStatsTableContent(m, rows, availableRowHeight)

	return header + content
}

// buildStatsRows creates StatsRow entries with weekly separators inserted
func buildStatsRows(aggregated []utils.ProjectDateEntry) []StatsRow {
	var rows []StatsRow

	// Get week boundaries
	separators := utils.GetWeekSeparators(aggregated)
	separatorSet := make(map[int]bool)
	for _, sep := range separators {
		separatorSet[sep] = true
	}

	// Build rows with separators
	for i, entry := range aggregated {
		// Insert separator before this index if needed
		if separatorSet[i] && len(rows) > 0 {
			// Find the week start for previous entry
			prevWeekStart := utils.GetMondayOfWeek(aggregated[i-1].Date)
			weekTotal := utils.GetWeeklyTotal(aggregated, prevWeekStart)
			rows = append(rows, StatsWeeklySeparatorRow(prevWeekStart, int(weekTotal.Minutes())))
		}

		rows = append(rows, StatsRowFromEntry(entry))
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
func renderStatsTableHeader(width int) string {
	// Column layout: Project | Date | Duration (min) | Description
	// Example: "ProjectA  2025-01-01  90  Task 1"

	projectCol := 15
	dateCol := 12
	durationCol := 12
	descCol := max(width-projectCol-dateCol-durationCol-8, 30) // -8 for separators and padding

	headerText := fmt.Sprintf(
		"%-*s %-*s %-*s %s",
		projectCol, "Project",
		dateCol, "Date",
		durationCol, "Duration",
		"Description",
	)

	output := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")).Render(headerText) + "\n"

	// Separator line
	separatorWidth := projectCol + dateCol + durationCol + descCol + 4
	separatorText := strings.Repeat("─", min(separatorWidth, width))
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

	var output strings.Builder

	// Render visible rows
	endIdx := min(m.ViewportTop+maxHeight, len(rows))

	for i := m.ViewportTop; i < endIdx; i++ {
		row := rows[i]

		if row.IsWeeklySeparator() {
			// Render weekly separator with different styling
			separatorText := fmt.Sprintf(
				"%-*s %-*s %-*s",
				projectCol, "WEEK OF",
				dateCol, row.WeekStartDate,
				durationCol, formatDurationMinutes(row.DurationMinutes),
			)
			styledRow := lipgloss.NewStyle().
				Background(lipgloss.Color("5")).
				Foreground(lipgloss.Color("0")).
				Bold(true).
				Render(separatorText)
			output.WriteString(styledRow + "\n")
		} else {
			// Render regular data row
			// Task descriptions as newline-separated bullet points
			taskDescr := ""
			if len(row.Tasks) > 0 {
				taskDescr = "• " + strings.Join(row.Tasks, "\n  • ")
			}

			durationFormatted := formatDurationMinutes(row.DurationMinutes)
			firstLineText := fmt.Sprintf(
				"%-*s %-*s %-*s %s",
				projectCol, row.Project,
				dateCol, row.Date,
				durationCol, durationFormatted,
				taskDescr,
			)

			styledRow := m.Styles.Unselected.Render(firstLineText)
			output.WriteString(styledRow + "\n")
		}
	}

	return output.String()
}
