package tui

import (
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
		content = "Recent entries:\n"
		// Show entries that fit in the available space
		displayCount := contentHeight - 1 // -1 for the "Recent entries:" line
		start := 0
		if len(m.entries) > displayCount {
			start = len(m.entries) - displayCount
		}
		for i := len(m.entries) - 1; i >= start; i-- {
			entry := m.entries[i]
			status := "●"
			if entry.End != nil {
				status = "○"
			}
			content += status + " " + entry.Project + " > " + entry.Title + "\n"
		}
	}

	// Build final view
	divider := m.styles.divider.Render(strings.Repeat("─", m.width))
	footer := m.help.View(m.keys)

	return header + "\n\n" + content + "\n" + divider + "\n" + footer
}
