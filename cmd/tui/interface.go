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
	keyBg := lipgloss.Color("5")   // Magenta for keys
	keyFg := lipgloss.Color("0")   // Black text on keys
	labelBg := lipgloss.Color("0") // Black for labels
	labelFg := lipgloss.Color("8") // Gray text for labels
	modeBg := lipgloss.Color("10") // Green for mode
	modeFg := lipgloss.Color("0")  // Black text on mode

	// Styles
	modeStyle := lipgloss.NewStyle().
		Background(modeBg).
		Foreground(modeFg).
		Bold(true).
		Padding(0, 1)

	modeSep := lipgloss.NewStyle().
		Background(labelBg).
		Foreground(modeBg).
		Render("\uE0B0")

	keyStyle := lipgloss.NewStyle().
		Background(keyBg).
		Foreground(keyFg).
		Padding(0, 1)

	keySepLeft := lipgloss.NewStyle().
		Background(keyBg).
		Foreground(labelBg).
		Render("\uE0B0")

	keySepRight := lipgloss.NewStyle().
		Background(labelBg).
		Foreground(keyBg).
		Render("\uE0B0")

	labelStyle := lipgloss.NewStyle().
		Background(labelBg).
		Foreground(labelFg).
		Padding(0, 1)

	// Helper to render a key-label pair with powerline separators
	renderPair := func(key, label string) string {
		return keySepLeft + keyStyle.Render(key) + keySepRight + labelStyle.Render(label)
	}

	var parts []string

	// Mode indicator and keybindings based on current mode
	switch m.mode {
	case ModeStart:
		parts = append(parts, modeStyle.Render("START")+modeSep)
		parts = append(parts, renderPair("Tab", "switch"))
		parts = append(parts, renderPair("↵", "submit"))
		parts = append(parts, renderPair("Esc", "cancel"))
	case ModeHelp:
		parts = append(parts, modeStyle.Render("HELP")+modeSep)
		parts = append(parts, renderPair("Esc", "close"))
	default: // ModeList
		parts = append(parts, modeStyle.Render("LIST")+modeSep)
		parts = append(parts, renderPair("j/k", "navigate"))
		parts = append(parts, renderPair("s", "start"))
		parts = append(parts, renderPair("?", "help"))
		parts = append(parts, renderPair("q", "quit"))
	}

	statusBar := strings.Join(parts, "")

	// Pad to full width with black background
	statusBarWidth := lipgloss.Width(statusBar)
	if statusBarWidth < m.width {
		padding := lipgloss.NewStyle().
			Background(labelBg).
			Render(strings.Repeat(" ", m.width-statusBarWidth))
		statusBar = statusBar + padding
	}

	return statusBar
}
