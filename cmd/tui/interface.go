package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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
		m.help.Width = msg.Width
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the UI
func (m *Model) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}

	// Header
	header := m.styles.header.Render("Time Tracker")
	
	// Content area - calculate available space
	helpView := m.help.View(m.keys)
	helpLines := strings.Count(helpView, "\n")
	contentHeight := m.height - 3 - helpLines // -3 for header and spacing
	if contentHeight < 3 {
		contentHeight = 3
	}

	// Entries list
	var content string
	if len(m.entries) == 0 {
		content = "No time entries yet.\n"
	} else {
		// Column headers
		content = fmt.Sprintf("%-20s %-20s %-15s %-20s %-10s\n", "START", "END", "PROJECT", "TITLE", "DURATION")
		content += strings.Repeat("─", m.width) + "\n"

		// Show entries that fit in the available space
		displayCount := contentHeight - 2 // -2 for headers and divider
		start := 0
		if len(m.entries) > displayCount {
			start = len(m.entries) - displayCount
		}
		
		// Display oldest to newest (start to end of available entries)
		for i := start; i < len(m.entries); i++ {
			entry := m.entries[i]
			
			// Format start time
			startStr := entry.Start.Format("2006-01-02 15:04")
			
			// Format end time or "running"
			endStr := "running"
			if entry.End != nil {
				endStr = entry.End.Format("2006-01-02 15:04")
			}
			
			// Format duration
			durationStr := ""
			duration := entry.Duration()
			hours := int(duration.Hours())
			minutes := int(duration.Minutes()) % 60
			if hours > 0 {
				durationStr = fmt.Sprintf("%dh %dm", hours, minutes)
			} else {
				durationStr = fmt.Sprintf("%dm", minutes)
			}
			
			// For blank entries, only show start/end/duration
			if entry.IsBlank() {
				content += fmt.Sprintf("%-20s %-20s %-15s %-20s %-10s\n", startStr, endStr, "", "", durationStr)
			} else {
				// Truncate long strings to fit columns
				project := entry.Project
				if len(project) > 15 {
					project = project[:12] + "..."
				}
				title := entry.Title
				if len(title) > 20 {
					title = title[:17] + "..."
				}
				
				content += fmt.Sprintf("%-20s %-20s %-15s %-20s %-10s\n", startStr, endStr, project, title, durationStr)
			}
		}
	}

	// Build final view
	divider := m.styles.divider.Render(strings.Repeat("─", m.width))
	footer := m.help.View(m.keys)

	return header + "\n\n" + content + "\n" + divider + "\n" + footer
}
