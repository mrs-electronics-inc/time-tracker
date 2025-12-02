package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
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

	switch m.currentScreen {
	case ScreenList:
		return m.viewList()
	case ScreenMenu:
		return m.viewMenu()
	default:
		return "Unknown screen\n"
	}
}

// viewList renders the entry list view
func (m *Model) viewList() string {
	output := "=== Time Tracker ===\n\n"
	
	if len(m.entries) == 0 {
		output += "No time entries yet.\n"
	} else {
		output += "Recent entries:\n"
		// Show up to 10 most recent entries
		start := 0
		if len(m.entries) > 10 {
			start = len(m.entries) - 10
		}
		for i := len(m.entries) - 1; i >= start; i-- {
			entry := m.entries[i]
			status := "●" // Running indicator
			if entry.End != nil {
				status = "○"
			}
			output += status + " " + entry.Project + " > " + entry.Title + "\n"
		}
	}

	output += "\nKeybindings:\n"
	output += "  q: quit\n"
	
	return output
}

// viewMenu renders the menu view
func (m *Model) viewMenu() string {
	return "Menu not yet implemented\n"
}
