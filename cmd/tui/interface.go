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
		// Start mode key handling
		if m.mode == ModeStart {
			return m.handleStartKeyMsg(msg)
		}

		// Help mode key handling
		if m.mode == ModeHelp {
			return m.handleHelpKeyMsg(msg)
		}

		// List mode key handling
		if key.Matches(msg, m.keys.Help) {
			m.prevMode = m.mode
			m.mode = ModeHelp
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

		// Toggle start/stop
		if key.Matches(msg, m.keys.Toggle) {
			if len(m.entries) == 0 {
				// No entries yet - open blank start mode
				m.openStartModeBlank()
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
					// Start new entry based on selected
					m.openStartMode(entry)
				} else {
					// Blank entry - open blank start mode
					m.openStartModeBlank()
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

	// If in start mode, render start screen
	if m.mode == ModeStart {
		return m.renderStartScreen()
	}

	// If in help mode, render help screen
	if m.mode == ModeHelp {
		return m.renderHelpScreen()
	}

	// Show loading indicator if operation in progress
	if m.loading {
		return m.renderLoading()
	}

	// Render footer first to know its height
	footer := m.renderFooter()
	footerHeight := strings.Count(footer, "\n") + 1

	// Header takes 2 lines (header + separator)
	headerHeight := 2

	// Available height for list rows
	availableHeight := max(m.height-headerHeight-footerHeight, 1)

	// Ensure selection is visible
	m.ensureSelectionVisible(availableHeight)

	// Render header and rows separately
	header := m.renderTableHeader()
	rows := m.renderTableRows(availableHeight)

	// Combine header and rows
	table := header + rows

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

	return strings.Join(parts, "")
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
	durationWidth += padding

	// Calculate available width for title column
	fixedWidth := startWidth + endWidth + projectWidth + durationWidth + 4 // 4 for column separators
	availableTitleWidth := max(m.width-fixedWidth, len("Title")+padding)
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
	output := m.styles.header.Render(headerText) + "\n"

	// Render separator
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
	durationWidth += padding

	// Calculate available width for title column
	fixedWidth := startWidth + endWidth + projectWidth + durationWidth + 4 // 4 for column separators
	availableTitleWidth := max(m.width-fixedWidth, len("Title")+padding)
	titleWidth = availableTitleWidth

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
		} else if entry.IsBlank() {
			endStr = "stopped"
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

// renderFooter renders the footer with status bar
func (m *Model) renderFooter() string {
	return m.renderStatusBar()
}

// renderStatusBar renders a zellij-style status bar with mode and keybindings
func (m *Model) renderStatusBar() string {
	// Colors
	black := lipgloss.Color("0")
	magenta := lipgloss.Color("5")
	gray := lipgloss.Color("8")
	green := lipgloss.Color("10")

	// Styles
	modeStyle := lipgloss.NewStyle().
		Background(green).
		Foreground(black).
		Bold(true).
		Padding(0, 1)

	keyStyle := lipgloss.NewStyle().
		Background(black).
		Foreground(magenta).
		Padding(0, 1)

	labelStyle := lipgloss.NewStyle().
		Background(gray).
		Foreground(black).
		Bold(true).
		Padding(0, 1)

	// Separators
	powerlineSeparator := "\uE0B0"

	modeSep := lipgloss.NewStyle().
		Background(black).
		Foreground(green).
		Render(powerlineSeparator)

	keySep := lipgloss.NewStyle().
		Background(gray).
		Foreground(black).
		Render(powerlineSeparator)

	labelSep := lipgloss.NewStyle().
		Background(black).
		Foreground(gray).
		Render(powerlineSeparator)

	// Helper to render a key-label pair with powerline separators
	renderPair := func(key, label string) string {
		return keyStyle.Render(key) + keySep + labelStyle.Render(label) + labelSep
	}

	var parts []string

	// Mode indicator and keybindings based on current mode
	switch m.mode {
	case ModeStart:
		parts = append(parts, modeStyle.Render("START")+modeSep)
		parts = append(parts, renderPair("Tab", "NEXT"))
		parts = append(parts, renderPair("Enter", "SUBMIT"))
		parts = append(parts, renderPair("Esc", "CANCEL"))
	case ModeHelp:
		parts = append(parts, modeStyle.Render("HELP")+modeSep)
		parts = append(parts, renderPair("Esc", "BACK"))
	default: // ModeList
		parts = append(parts, modeStyle.Render("LIST")+modeSep)
		parts = append(parts, renderPair("j/k", "NAVIGATE"))
		parts = append(parts, renderPair("G", "GO TO CURRENT"))
		parts = append(parts, renderPair("s", "START/STOP"))
		parts = append(parts, renderPair("?", "HELP"))
		parts = append(parts, renderPair("Esc", "QUIT"))
	}

	// Build left side of status bar
	leftSide := strings.Join(parts, "")

	// Add status message on the right side if present
	if m.status != "" {
		statusStyle := lipgloss.NewStyle().
			Foreground(magenta).
			Padding(0, 1)
		rightSide := statusStyle.Render(m.status)

		// Calculate padding to right-align status
		leftWidth := lipgloss.Width(leftSide)
		rightWidth := lipgloss.Width(rightSide)
		totalWidth := leftWidth + rightWidth
		paddingWidth := max(m.width-totalWidth, 0)

		padding := strings.Repeat(" ", paddingWidth)
		return leftSide + padding + rightSide
	}

	return leftSide
}
