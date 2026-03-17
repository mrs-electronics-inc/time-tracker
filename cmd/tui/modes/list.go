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
		{Keys: "Tab", Label: "STATS", Description: "Switch mode"},
		{Keys: "n", Label: "NEW", Description: "New entry"},
		{Keys: "s", Label: "STOP", Description: "Stop running entry"},
		{Keys: "r", Label: "RESUME", Description: "Resume entry"},
		{Keys: "e", Label: "EDIT", Description: "Edit entry"},
		{Keys: "d", Label: "DELETE", Description: "Delete entry"},
		{Keys: "?", Label: "HELP", Description: "Toggle help"},
		{Keys: "q", Label: "QUIT", Description: "Quit"},
	},
	HandleKeyMsg: func(m *Model, msg tea.KeyMsg) (*Model, tea.Cmd) {
		switch msg.String() {
		case "enter":
			if m.SearchActive {
				applySearch(m)
			}
			return m, nil

		case "/":
			m.SearchActive = true
			return m, nil

		case "tab":
			m.SwitchMode(m.StatsMode)
			return m, nil

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

		case "?":
			m.PreviousMode = m.CurrentMode
			m.CurrentMode = m.HelpMode
			return m, nil

		case "esc":
			if m.SearchActive {
				clearSearch(m)
				return m, nil
			}
			return m, tea.Quit

		case "q":
			return m, tea.Quit

		case "k", "up":
			if !moveSelectionInFilteredEntries(m, -1) && m.SelectedIdx > 0 {
				m.SelectedIdx--
			}
			m.Status = ""
			return m, nil

		case "j", "down":
			if !moveSelectionInFilteredEntries(m, 1) && m.SelectedIdx < len(m.Entries)-1 {
				m.SelectedIdx++
			}
			m.Status = ""
			return m, nil

		case "G":
			if !moveSelectionToFilteredEnd(m) && len(m.Entries) > 0 {
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
		searchBarHeight := 0
		if m.SearchActive {
			searchBarHeight = 1
		}

		listRowHeight := max(availableHeight-headerHeight-searchBarHeight, 1)
		visibleRows := getVisibleRows(m)
		ensureSelectionVisibleInRows(m, visibleRows, listRowHeight)

		header := renderTableHeader(m)
		rows := renderTableRows(m, listRowHeight)
		searchInputBar := ""
		if m.SearchActive {
			searchInputBar = renderSearchInputBar(m)
		}

		return header + rows + searchInputBar
	},
}

func getVisibleRows(m *Model) []VisibleEntry {
	if m.SearchAppliedQuery != "" {
		return m.FilteredEntries
	}

	visibleRows := make([]VisibleEntry, 0, len(m.Entries))
	for sourceIndex, entry := range m.Entries {
		visibleRows = append(visibleRows, VisibleEntry{
			Entry:       entry,
			SourceIndex: sourceIndex,
		})
	}

	return visibleRows
}

func ensureSelectionVisibleInRows(m *Model, rows []VisibleEntry, maxVisibleRows int) {
	if len(rows) == 0 {
		m.ViewportTop = 0
		return
	}

	selectedRowIdx := -1
	for i, row := range rows {
		if row.SourceIndex == m.SelectedIdx {
			selectedRowIdx = i
			break
		}
	}

	if selectedRowIdx >= 0 {
		if selectedRowIdx < m.ViewportTop {
			m.ViewportTop = selectedRowIdx
		} else if selectedRowIdx >= m.ViewportTop+maxVisibleRows {
			m.ViewportTop = selectedRowIdx - maxVisibleRows + 1
		}
	}

	if m.ViewportTop > len(rows)-maxVisibleRows {
		m.ViewportTop = len(rows) - maxVisibleRows
	}
	if m.ViewportTop < 0 {
		m.ViewportTop = 0
	}
}

func renderSearchInputBar(m *Model) string {
	return m.Styles.Footer.Render("Search: "+m.SearchQueryDraft) + "\n"
}

// isValidSelection checks if the selected index is valid
func isValidSelection(m *Model) bool {
	if m.SearchAppliedQuery != "" && len(m.FilteredEntries) == 0 {
		return false
	}

	return len(m.Entries) > 0 && m.SelectedIdx >= 0 && m.SelectedIdx < len(m.Entries)
}

func moveSelectionInFilteredEntries(m *Model, direction int) bool {
	if m.SearchAppliedQuery == "" || len(m.FilteredEntries) == 0 {
		return false
	}

	currentFilteredIndex := -1
	for i, visible := range m.FilteredEntries {
		if visible.SourceIndex == m.SelectedIdx {
			currentFilteredIndex = i
			break
		}
	}

	if currentFilteredIndex == -1 {
		return false
	}

	nextIndex := currentFilteredIndex + direction
	if nextIndex < 0 || nextIndex >= len(m.FilteredEntries) {
		return true
	}

	m.SelectedIdx = m.FilteredEntries[nextIndex].SourceIndex
	return true
}

func moveSelectionToFilteredEnd(m *Model) bool {
	if m.SearchAppliedQuery == "" || len(m.FilteredEntries) == 0 {
		return false
	}

	m.SelectedIdx = m.FilteredEntries[len(m.FilteredEntries)-1].SourceIndex
	return true
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
	if m.SearchAppliedQuery != "" && len(m.FilteredEntries) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)
		msg := "No matching time entries found. Press 'esc' to clear search.\n"
		return emptyStyle.Render(msg)
	}

	if len(m.Entries) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)
		msg := "No time entries found. Press 'n' to start tracking.\n"
		return emptyStyle.Render(msg)
	}

	startWidth, endWidth, projectWidth, titleWidth, durationWidth := getColumnWidthsWithPadding(m)
	visibleRows := getVisibleRows(m)

	var output strings.Builder

	// Render rows from viewport
	maxRows := maxHeight
	rowsRendered := 0
	endIdx := min(m.ViewportTop+maxRows, len(visibleRows))

	for i := m.ViewportTop; i < endIdx; i++ {
		visible := visibleRows[i]
		entry := visible.Entry

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
		if visible.SourceIndex == m.SelectedIdx {
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
