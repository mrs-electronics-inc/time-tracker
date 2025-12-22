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
		{Keys: "k / ↑", Label: "UP", Description: "Move up"},
		{Keys: "j / ↓", Label: "DOWN", Description: "Move down"},
		{Keys: "G", Label: "GO TO CURRENT", Description: "Go to current"},
		{Keys: "s", Label: "START/STOP", Description: "Start/stop entry"},
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

		case "s":
			if len(m.Entries) == 0 {
				openStartModeBlank(m)
			} else if m.SelectedIdx >= 0 && m.SelectedIdx < len(m.Entries) {
				entry := m.Entries[m.SelectedIdx]
				if entry.IsRunning() {
					if _, err := m.TaskManager.StopEntry(); err != nil {
						m.Status = "Error stopping entry: " + err.Error()
					} else {
						m.Status = "Entry stopped"
					}
				} else if !entry.IsBlank() {
					openStartMode(m, entry)
				} else {
					openStartModeBlank(m)
				}
				if err := m.LoadEntries(); err != nil {
					m.Err = err
				}
			}
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

// renderTableHeader renders just the table header
func renderTableHeader(m *Model) string {
	if len(m.Entries) == 0 {
		return ""
	}

	// Get column widths
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
		msg := "No time entries found. Press 's' to start tracking.\n"
		return emptyStyle.Render(msg)
	}

	// Get column widths
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
