package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case OperationCompleteMsg:
		m.loading = false
		if msg.Error != nil {
			m.status = "Error: " + msg.Error.Error()
		}
		return m, nil

	case tea.KeyMsg:
		// Dialog mode key handling
		if m.dialogMode {
			return m.handleDialogKeyMsg(msg)
		}

		// List mode key handling
		if key.Matches(msg, m.keys.Help) {
			m.showHelp = !m.showHelp
			return m, nil
		}

		if key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		}

		// Navigation
		if key.Matches(msg, m.keys.Up) {
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
			m.status = ""
			return m, nil
		}
		if key.Matches(msg, m.keys.Down) {
			if m.selectedIdx < len(m.entries)-1 {
				m.selectedIdx++
			}
			m.status = ""
			return m, nil
		}

		// Jump to bottom
		if key.Matches(msg, m.keys.JumpBottom) {
			if len(m.entries) > 0 {
				m.selectedIdx = len(m.entries) - 1
			}
			m.status = ""
			return m, nil
		}

		// Toggle start/stop or open dialog to start
		if key.Matches(msg, m.keys.Toggle) {
			if len(m.entries) == 0 {
				// No entries yet - open blank start dialog
				m.openStartDialogBlank()
			} else if m.selectedIdx >= 0 && m.selectedIdx < len(m.entries) {
				entry := m.entries[m.selectedIdx]
				if entry.IsRunning() {
					// Stop entry
					if _, err := m.taskManager.StopEntry(); err != nil {
						m.status = "Error stopping entry: " + err.Error()
					} else {
						m.status = "Entry stopped"
					}
				} else if !entry.IsBlank() {
					// Open dialog to start new entry
					m.openStartDialog(entry)
				} else {
					// Blank entry - open blank start dialog
					m.openStartDialogBlank()
				}
				// Reload entries to update display
				if err := m.LoadEntries(); err != nil {
					m.err = err
				}
			}
			return m, nil
		}
	}

	return m, nil
}

// View renders the UI
func (m *Model) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}

	// Show loading indicator if operation in progress
	if m.loading {
		return m.renderLoading()
	}

	// Render footer first to know its height (skip if in dialog mode)
	var footer string
	footerHeight := 0
	if !m.dialogMode {
		footer = m.renderFooter()
		footerHeight = strings.Count(footer, "\n") + 1
	}

	// Header takes 2 lines (header + separator)
	headerHeight := 2

	// Status message takes 1 line if present
	statusHeight := 0
	if m.status != "" {
		statusHeight = 1
	}

	// Available height for list rows
	availableHeight := max(m.height-headerHeight-footerHeight-statusHeight, 1)

	// Ensure selection is visible
	m.ensureSelectionVisible(availableHeight)

	// Render header and rows separately
	header := m.renderTableHeader()
	rows := m.renderTableRows(availableHeight)

	// Combine header and rows
	table := header + rows

	// Add status message if present
	if m.status != "" {
		table = table + m.status + "\n"
	}

	// Calculate spacer to push footer to bottom
	tableLines := strings.Count(table, "\n")
	usedLines := tableLines + footerHeight
	// Ensure spacer height is not negative
	spacerHeight := max(m.height-usedLines, 0)

	// Build layout with spacer
	var parts []string
	parts = append(parts, table)

	if spacerHeight > 0 {
		spacer := strings.Repeat("\n", spacerHeight)
		parts = append(parts, spacer)
	}

	parts = append(parts, footer)

	output := strings.Join(parts, "")

	// If in dialog mode, render as centered overlay
	if m.dialogMode {
		dialogContent := m.renderDialog()
		output = m.compositeOverlay(output, dialogContent)
	}

	return output
}

// renderTableHeader renders just the table header
func (m *Model) renderTableHeader() string {
	if len(m.entries) == 0 {
		return ""
	}

	// Get column widths
	startWidth, endWidth, projectWidth, titleWidth, durationWidth := m.getColumnWidths()

	// Add some padding
	padding := 1
	startWidth += padding
	endWidth += padding
	projectWidth += padding
	titleWidth += padding
	durationWidth += padding

	// Render header
	headerText := fmt.Sprintf(
		"%-*s %-*s %-*s %-*s %s",
		startWidth, "Start",
		endWidth, "End",
		projectWidth, "Project",
		titleWidth, "Title",
		"Duration",
	)
	output := m.styles.header.Render(headerText) + "\n"

	// Render separator (4 spaces for column separators between 5 columns)
	separatorWidth := startWidth + endWidth + projectWidth + titleWidth + durationWidth + 4
	separatorText := strings.Repeat("─", separatorWidth)
	output += m.styles.header.Render(separatorText) + "\n"

	return output
}

// renderTableRows renders the rows with viewport scrolling
func (m *Model) renderTableRows(maxHeight int) string {
	if len(m.entries) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)
		msg := "No time entries found. Press 's' to start tracking.\n"
		return emptyStyle.Render(msg)
	}

	// Get column widths
	startWidth, endWidth, projectWidth, titleWidth, durationWidth := m.getColumnWidths()

	// Add some padding
	padding := 1
	startWidth += padding
	endWidth += padding
	projectWidth += padding
	titleWidth += padding
	durationWidth += padding

	var output strings.Builder

	// Render rows from viewport
	maxRows := maxHeight
	rowsRendered := 0
	endIdx := min(m.viewportTop+maxRows, len(m.entries))

	for i := m.viewportTop; i < endIdx; i++ {
		entry := m.entries[i]

		startStr := entry.Start.Format("2006-01-02 15:04")

		endStr := "running"
		if entry.End != nil {
			endStr = entry.End.Format("2006-01-02 15:04")
		}

		project := entry.Project
		title := entry.Title

		duration := formatDuration(entry.Duration())

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
		if i == m.selectedIdx {
			// Selected row - highlight with bold and inverse
			styledRow = lipgloss.NewStyle().
				Bold(true).
				Reverse(true).
				Render(row)
		} else if entry.IsRunning() {
			// Running entry - use running style
			styledRow = m.styles.running.Render(row)
		} else if entry.IsBlank() {
			// Gap entry - use gap style
			styledRow = m.styles.gap.Render(row)
		} else {
			// Regular unselected - use unselected style
			styledRow = m.styles.unselected.Render(row)
		}

		output.WriteString(styledRow + "\n")
		rowsRendered++
	}

	return output.String()
}

// compositeOverlay overlays a dialog centered on the background using simple string compositing
func (m *Model) compositeOverlay(background, foreground string) string {
	bgLines := strings.Split(background, "\n")
	fgLines := strings.Split(foreground, "\n")

	// Calculate dimensions (use lipgloss.Width for proper terminal width with ANSI codes)
	fgHeight := len(fgLines)
	fgWidth := 0
	for _, line := range fgLines {
		width := lipgloss.Width(line)
		if width > fgWidth {
			fgWidth = width
		}
	}

	// Calculate starting positions for centering
	startRow := (m.height - fgHeight) / 2
	startCol := (m.width - fgWidth) / 2

	// Ensure non-negative offsets
	if startRow < 0 {
		startRow = 0
	}
	if startCol < 0 {
		startCol = 0
	}

	// Pad background to ensure it has enough lines
	for len(bgLines) < m.height {
		bgLines = append(bgLines, "")
	}

	// Overlay foreground onto background
	for i := 0; i < fgHeight && startRow+i < len(bgLines); i++ {
		bgLine := bgLines[startRow+i]
		fgLine := fgLines[i]

		// Pad background line to at least startCol characters
		bgLineWidth := lipgloss.Width(bgLine)
		if bgLineWidth < startCol {
			bgLine = bgLine + strings.Repeat(" ", startCol-bgLineWidth)
		}

		// Replace the section of bgLine with fgLine, using terminal widths not string lengths
		fgLineWidth := lipgloss.Width(fgLine)
		if bgLineWidth < startCol+fgLineWidth {
			// Dialog extends beyond background line, just append
			bgLines[startRow+i] = bgLine + fgLine
		} else {
			// Need to handle ANSI codes - just use simple replacement for now
			// This is a limitation: we're replacing based on display width but the string length doesn't match
			bgLines[startRow+i] = bgLine[:startCol] + fgLine
		}
	}

	// Trim trailing empty lines to original height
	if len(bgLines) > m.height {
		bgLines = bgLines[:m.height]
	}

	return strings.Join(bgLines, "\n")
}

// renderLoading renders a loading indicator
func (m *Model) renderLoading() string {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	frame := frames[int(time.Now().Unix()*10)%len(frames)]

	loadingText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("11")).
		Bold(true).
		Render(frame + " Loading...")

	return "\n\n" + loadingText + "\n"
}

// renderFooter renders the footer with help text
func (m *Model) renderFooter() string {
	m.help.Width = m.width
	m.help.ShowAll = m.showHelp
	return m.styles.footer.Render(m.help.View(m.keys))
}
