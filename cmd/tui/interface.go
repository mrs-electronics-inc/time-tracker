package tui

import (
	"fmt"
	"strings"

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

	case tea.KeyMsg:
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

		// Toggle start/stop
		if key.Matches(msg, m.keys.Toggle) {
			if m.selectedIdx >= 0 && m.selectedIdx < len(m.entries) {
				entry := m.entries[m.selectedIdx]
				if entry.IsRunning() {
					// Stop entry
					if _, err := m.taskManager.StopEntry(); err != nil {
						m.status = "Error stopping entry: " + err.Error()
					} else {
						m.status = "Entry stopped"
					}
				} else if !entry.IsBlank() {
					// Start new entry from selected entry's project/title
					if _, err := m.taskManager.StartEntry(entry.Project, entry.Title); err != nil {
						m.status = "Error starting entry: " + err.Error()
					} else {
						m.status = "Entry started: " + entry.Project
					}
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

	// Render table
	table := m.renderTable()
	
	// Add status message if present
	if m.status != "" {
		table = table + m.status + "\n"
	}
	
	// Render footer
	footer := m.renderFooter()

	// Calculate spacer height to push footer to bottom
	contentLines := strings.Count(table, "\n")
	footerLines := strings.Count(footer, "\n") + 1
	spacerHeight := m.height - contentLines - footerLines
	
	// Build layout with spacer
	if spacerHeight > 0 {
		spacer := lipgloss.NewStyle().Height(spacerHeight).Render("")
		return lipgloss.JoinVertical(
			lipgloss.Top,
			table,
			spacer,
			footer,
		)
	}
	
	// If not enough space, just render table and footer
	return table + footer
}

// renderTable renders the table with entries
func (m *Model) renderTable() string {
	if len(m.entries) == 0 {
		return "No time entries found\n"
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

	// Render header
	headerText := fmt.Sprintf(
		"%-*s %-*s %-*s %-*s %s",
		startWidth, "Start",
		endWidth, "End",
		projectWidth, "Project",
		titleWidth, "Title",
		"Duration",
	)
	output.WriteString(headerText + "\n")

	// Render separator (4 spaces for column separators between 5 columns)
	separatorWidth := startWidth + endWidth + projectWidth + titleWidth + durationWidth + 4
	separatorText := strings.Repeat("─", separatorWidth)
	output.WriteString(separatorText + "\n")

	// Render rows
	for i, entry := range m.entries {
		startStr := entry.Start.Format("2006-01-02 15:04")

		endStr := "running"
		if entry.End != nil {
			endStr = entry.End.Format("2006-01-02 15:04")
		}

		project := entry.Project
		title := entry.Title
		if entry.IsBlank() {
			project = ""
			title = ""
		}

		duration := formatDuration(entry.Duration())

		row := fmt.Sprintf(
			"%-*s %-*s %-*s %-*s %-*s",
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
	}

	return output.String()
}

// renderFooter renders the footer with help text
func (m *Model) renderFooter() string {
	// Determine if selected entry is running
	toggleAction := "start"
	if m.selectedIdx >= 0 && m.selectedIdx < len(m.entries) {
		if m.entries[m.selectedIdx].IsRunning() {
			toggleAction = "stop"
		}
	}

	if m.showHelp {
		helpText := "Keybindings:\n"
		helpText += "  j/↓ - down          s - " + toggleAction + "\n"
		helpText += "  k/↑ - up            ? - toggle help\n"
		helpText += "  q/esc - quit\n"
		return m.styles.footer.Render(helpText)
	}
	// Show compact help footer
	return m.styles.footer.Render("j/k ↑↓: navigate | s: " + toggleAction + " | ?: help | q: quit")
}
