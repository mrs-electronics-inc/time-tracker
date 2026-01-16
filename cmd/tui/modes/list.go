package modes

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time-tracker/utils"
)

// ListMode is the list view mode
var ListMode = &Mode{
	Name: "list",
	KeyBindings: []KeyBinding{
		{Keys: "n", Label: "NEW", Description: "New entry"},
		{Keys: "s", Label: "STOP", Description: "Stop running entry"},
		{Keys: "r", Label: "RESUME", Description: "Resume entry"},
		{Keys: "e", Label: "EDIT", Description: "Edit entry"},
		{Keys: "d", Label: "DELETE", Description: "Delete entry"},
		{Keys: "Tab", Label: "STATS", Description: "Switch mode"},
		{Keys: "?", Label: "HELP", Description: "Toggle help"},
		{Keys: "q", Label: "QUIT", Description: "Quit"},
	},
	HandleKeyMsg: func(m *Model, msg tea.KeyMsg) (*Model, tea.Cmd) {
		switch msg.String() {
		case "n":
			openNewMode(m)
			return m, nil

		case "s":
			// Stop only works on running entries
			if isValidSelection(m) {
				entry := m.Entries[m.SelectedIdx]
				if entry.IsRunning() {
					if _, err := m.TaskManager.StopEntry(); err != nil {
						m.Status = "Error stopping entry: " + err.Error()
					} else {
						m.Status = "Entry stopped"
					}
					if err := m.LoadEntries(); err != nil {
						m.Err = err
					}
				}
			}
			return m, nil

		case "r":
			if isValidSelection(m) {
				entry := m.Entries[m.SelectedIdx]
				if !entry.IsBlank() {
					openResumeMode(m, entry)
				}
			}
			return m, nil

		case "e":
			if isValidSelection(m) {
				entry := m.Entries[m.SelectedIdx]
				openEditMode(m, entry, m.SelectedIdx)
			}
			return m, nil

		case "d":
			if isValidSelection(m) {
				openConfirmDelete(m, m.SelectedIdx)
			}
			return m, nil

		case "tab":
			m.SwitchMode(m.StatsMode)
			return m, nil

		case "?":
			m.PreviousMode = m.CurrentMode
			m.CurrentMode = m.HelpMode
			return m, nil

		case "q", "esc":
			return m, tea.Quit

		case "k", "up":
			if m.SelectedIdx > 0 {
				m.SelectedIdx--
			}
			m.Status = ""
			return m, nil

		case "j", "down":
			if m.SelectedIdx < len(m.Entries)-1 {
				m.SelectedIdx++
			}
			m.Status = ""
			return m, nil

		case "G":
			if len(m.Entries) > 0 {
				m.SelectedIdx = len(m.Entries) - 1
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

		headerHeight := 2
		listRowHeight := max(availableHeight-headerHeight, 1)
		m.EnsureSelectionVisible(listRowHeight)

		header := renderTableHeader(m)
		rows := renderTableRows(m, listRowHeight)

		return header + rows
	},
}

// isValidSelection checks if the selected index is valid
func isValidSelection(m *Model) bool {
	return len(m.Entries) > 0 && m.SelectedIdx >= 0 && m.SelectedIdx < len(m.Entries)
}

// renderLoading renders a loading indicator
func renderLoading() string {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	frame := frames[int(time.Now().Unix()*10)%len(frames)]

	loadingText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("11")).
		Bold(true).
		Render(frame + " Loading...")

	return "\n\n" + loadingText + "\n"
}

// getColumnWidthsWithPadding calculates column widths with padding applied
func getColumnWidthsWithPadding(m *Model) (int, int, int, int, int) {
	startWidth, endWidth, projectWidth, titleWidth, durationWidth := m.GetColumnWidths()

	// Add some padding
	padding := 1
	startWidth += padding
	endWidth += padding
	projectWidth += padding
	durationWidth += padding

	// Calculate available width for title column
	fixedWidth := startWidth + endWidth + projectWidth + durationWidth + 4 // 4 for column separators
	availableTitleWidth := max(m.Width-fixedWidth, len("Title")+padding)
	titleWidth = availableTitleWidth

	return startWidth, endWidth, projectWidth, titleWidth, durationWidth
}

// renderTableHeader renders just the table header
func renderTableHeader(m *Model) string {
	if len(m.Entries) == 0 {
		return ""
	}

	startWidth, endWidth, projectWidth, titleWidth, durationWidth := getColumnWidthsWithPadding(m)

	// Render header
	headerText := fmt.Sprintf(
		"%-*s %-*s %-*s %-*s %s",
		startWidth, "Start",
		endWidth, "End",
		projectWidth, "Project",
		titleWidth, "Title",
		"Duration",
	)
	output := m.Styles.Header.Render(headerText) + "\n"

	// Render separator
	separatorWidth := startWidth + endWidth + projectWidth + titleWidth + durationWidth + 4
	separatorText := strings.Repeat("─", separatorWidth)
	output += m.Styles.Header.Render(separatorText) + "\n"

	return output
}

// renderTableRows renders the rows with viewport scrolling
func renderTableRows(m *Model, maxHeight int) string {
	if len(m.Entries) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)
		msg := "No time entries found. Press 'n' to start tracking.\n"
		return emptyStyle.Render(msg)
	}

	startWidth, endWidth, projectWidth, titleWidth, durationWidth := getColumnWidthsWithPadding(m)

	var output strings.Builder

	// Render rows from viewport
	maxRows := maxHeight
	rowsRendered := 0
	endIdx := min(m.ViewportTop+maxRows, len(m.Entries))

	for i := m.ViewportTop; i < endIdx; i++ {
		entry := m.Entries[i]

		startStr := entry.Start.Format("2006-01-02 15:04")

		endStr := "running"
		if entry.End != nil {
			endStr = entry.End.Format("2006-01-02 15:04")
		} else if entry.IsBlank() {
			endStr = "stopped"
		}

		project := entry.Project
		title := entry.Title

		duration := utils.FormatDuration(entry.Duration())

		row := fmt.Sprintf(
			"%-*s %-*s %-*s %-*s %*s",
			startWidth, startStr,
			endWidth, endStr,
			projectWidth, project,
			titleWidth, title,
			durationWidth, duration,
		)

		// Apply styling
		var styledRow string
		if i == m.SelectedIdx {
			// Selected row - highlight with bold and inverse
			styledRow = lipgloss.NewStyle().
				Bold(true).
				Reverse(true).
				Render(row)
		} else if entry.IsRunning() {
			// Running entry - use running style
			styledRow = m.Styles.Running.Render(row)
		} else if entry.IsBlank() {
			// Gap entry - use gap style
			styledRow = m.Styles.Gap.Render(row)
		} else {
			// Regular unselected - use unselected style
			styledRow = m.Styles.Unselected.Render(row)
		}

		output.WriteString(styledRow + "\n")
		rowsRendered++
	}

	return output.String()
}
