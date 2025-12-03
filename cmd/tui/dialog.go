package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time-tracker/models"
)

// openStartDialog opens the start entry dialog with pre-filled values
func (m *Model) openStartDialog(entry models.TimeEntry) {
	m.dialogMode = true
	m.focusIndex = 0
	m.showAutocomplete = true

	// Pre-fill the inputs with the selected entry's values
	m.inputs[0].SetValue(entry.Project)
	m.inputs[1].SetValue(entry.Title)

	// Set current time as default
	now := time.Now()
	m.inputs[2].SetValue(fmt.Sprintf("%02d", now.Hour()))
	m.inputs[3].SetValue(fmt.Sprintf("%02d", now.Minute()))

	// Set focus to first input
	m.inputs[0].Focus()
	m.inputs[0].PromptStyle = m.styles.dialogFocused
	m.inputs[0].TextStyle = m.styles.dialogFocused

	// Blur other inputs
	for i := 1; i < len(m.inputs); i++ {
		m.inputs[i].Blur()
		m.inputs[i].PromptStyle = m.styles.dialogBlurred
		m.inputs[i].TextStyle = m.styles.dialogBlurred
	}

	// Initialize autocomplete filters
	m.updateAutocompleteFilter()
}

// closeDialog closes the dialog and returns to list mode
func (m *Model) closeDialog() {
	m.dialogMode = false
	m.focusIndex = 0
	m.showAutocomplete = false

	// Clear and blur all inputs
	for i := range m.inputs {
		m.inputs[i].SetValue("")
		m.inputs[i].Blur()
	}
}

// handleDialogKeyMsg handles key messages while in dialog mode
func (m *Model) handleDialogKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Cancel dialog
		m.closeDialog()
		return m, nil

	case "tab":
		// Move focus to next input
		m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
		m.updateInputFocus()
		m.updateAutocompleteFilter()
		return m, nil

	case "shift+tab":
		// Move focus to previous input
		m.focusIndex--
		if m.focusIndex < 0 {
			m.focusIndex = len(m.inputs) - 1
		}
		m.updateInputFocus()
		m.updateAutocompleteFilter()
		return m, nil

	case "down":
		// If autocomplete is showing, move to next suggestion
		// Otherwise move to next input field
		if m.showAutocomplete && m.focusIndex < 2 {
			hasResults := false
			if m.focusIndex == 0 && len(m.autocomplete.FilteredProjects) > 0 {
				hasResults = true
			} else if m.focusIndex == 1 && len(m.autocomplete.FilteredResults) > 0 {
				hasResults = true
			}
			if hasResults {
				m.autocomplete.SelectNext()
				return m, nil
			}
		}
		// Move focus to next input
		m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
		m.updateInputFocus()
		m.updateAutocompleteFilter()
		return m, nil

	case "up":
		// If autocomplete is showing, move to previous suggestion
		// Otherwise move to previous input field
		if m.showAutocomplete && m.focusIndex < 2 {
			hasResults := false
			if m.focusIndex == 0 && len(m.autocomplete.FilteredProjects) > 0 {
				hasResults = true
			} else if m.focusIndex == 1 && len(m.autocomplete.FilteredResults) > 0 {
				hasResults = true
			}
			if hasResults {
				m.autocomplete.SelectPrev()
				return m, nil
			}
		}
		// Move focus to previous input
		m.focusIndex--
		if m.focusIndex < 0 {
			m.focusIndex = len(m.inputs) - 1
		}
		m.updateInputFocus()
		m.updateAutocompleteFilter()
		return m, nil

	case "enter":
		// If autocomplete is showing and a suggestion is selected, auto-fill it
		if m.showAutocomplete && m.focusIndex < 2 {
			if m.focusIndex == 0 {
				// Fill project from filtered projects and move to title
				if m.autocomplete.selectedIdx >= 0 && m.autocomplete.selectedIdx < len(m.autocomplete.FilteredProjects) {
					project := m.autocomplete.FilteredProjects[m.autocomplete.selectedIdx]
					m.inputs[0].SetValue(project)
					m.focusIndex = 1
					m.updateInputFocus()
					m.updateAutocompleteFilter()
					return m, nil
				}
			} else if m.focusIndex == 1 {
				// Fill title from filtered tasks
				if selected := m.autocomplete.GetSelectedSuggestion(); selected != nil {
					m.inputs[1].SetValue(selected.Title)
					m.showAutocomplete = false
					return m, nil
				}
			}
		}

		// Otherwise, submit dialog
		project := m.inputs[0].Value()
		title := m.inputs[1].Value()
		hourStr := m.inputs[2].Value()
		minuteStr := m.inputs[3].Value()

		if project == "" || title == "" {
			m.status = "Project and title are required"
			return m, nil
		}

		// Validate and parse time
		if hourStr == "" {
			hourStr = "00"
		}
		if minuteStr == "" {
			minuteStr = "00"
		}

		var hour, minute int
		if n, err := fmt.Sscanf(hourStr, "%d", &hour); err != nil || n != 1 || hour < 0 || hour > 23 {
			m.status = "Invalid hour (0-23)"
			return m, nil
		}
		if n, err := fmt.Sscanf(minuteStr, "%d", &minute); err != nil || n != 1 || minute < 0 || minute > 59 {
			m.status = "Invalid minute (0-59)"
			return m, nil
		}

		// Build the start time
		now := time.Now()
		date := now
		
		// If entered time is later than current time, use yesterday's date
		if hour > now.Hour() || (hour == now.Hour() && minute > now.Minute()) {
			date = now.AddDate(0, 0, -1)
		}
		
		startTime := time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, date.Location())

		// Start the entry with specified time
		if _, err := m.taskManager.StartEntryAt(project, title, startTime); err != nil {
			m.status = "Error starting entry: " + err.Error()
		} else {
			m.status = "Entry started: " + project
		}

		// Reload entries and close dialog
		if err := m.LoadEntries(); err != nil {
			m.err = err
		}
		m.closeDialog()
		return m, nil
	}

	// Route other key messages to the focused input
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	// After input update, refresh autocomplete filter if in project/title inputs
	if m.focusIndex < 2 {
		m.updateAutocompleteFilter()
	}

	return m, tea.Batch(cmds...)
}

// updateInputFocus updates the focus styling on all inputs
func (m *Model) updateInputFocus() {
	for i := range m.inputs {
		if i == m.focusIndex {
			// Set focused state
			m.inputs[i].Focus()
			m.inputs[i].PromptStyle = m.styles.dialogFocused
			m.inputs[i].TextStyle = m.styles.dialogFocused
		} else {
			// Set blurred state
			m.inputs[i].Blur()
			m.inputs[i].PromptStyle = m.styles.dialogBlurred
			m.inputs[i].TextStyle = m.styles.dialogBlurred
		}
	}
}

// updateAutocompleteFilter updates the autocomplete suggestions based on current input
func (m *Model) updateAutocompleteFilter() {
	if m.focusIndex == 0 {
		// Filter projects based on project input
		input := m.inputs[0].Value()
		m.autocomplete.FilterProjects(input)
	} else if m.focusIndex == 1 {
		// Filter tasks based on project and title inputs
		project := m.inputs[0].Value()
		title := m.inputs[1].Value()
		m.autocomplete.FilterTasks(title, project)
	}
}

// renderAutocompleteList renders the autocomplete suggestions list
func (m *Model) renderAutocompleteList(fieldIndex int) string {
	if !m.showAutocomplete {
		return ""
	}

	// Don't show autocomplete for hour/minute fields
	if fieldIndex > 1 {
		return ""
	}

	// Only show autocomplete for the focused field
	if m.focusIndex != fieldIndex {
		return ""
	}

	var output strings.Builder
	suggestionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Bold(true)

	if fieldIndex == 0 {
		// Show filtered projects
		suggestions := m.autocomplete.FilteredProjects
		
		// If no suggestions, just return a blank line for spacing
		if len(suggestions) == 0 {
			return "\n"
		}

		// Limit to 5 visible suggestions
		maxSuggestions := 5
		if len(suggestions) > maxSuggestions {
			suggestions = suggestions[:maxSuggestions]
		}

		output.WriteString("\n")
		for i, project := range suggestions {
			if i == m.autocomplete.selectedIdx {
				output.WriteString(selectedStyle.Render("▶ " + project))
			} else {
				output.WriteString(suggestionStyle.Render("  " + project))
			}
			output.WriteString("\n")
		}
	} else if fieldIndex == 1 {
		// Show filtered tasks (from selected project only)
		suggestions := m.autocomplete.FilteredResults
		
		// If no suggestions, just return a blank line for spacing
		if len(suggestions) == 0 {
			return "\n"
		}

		// Limit to 5 visible suggestions
		maxSuggestions := 5
		if len(suggestions) > maxSuggestions {
			suggestions = suggestions[:maxSuggestions]
		}

		output.WriteString("\n")
		for i, task := range suggestions {
			if i == m.autocomplete.selectedIdx {
				output.WriteString(selectedStyle.Render("▶ " + task.Title))
			} else {
				output.WriteString(suggestionStyle.Render("  " + task.Title))
			}
			output.WriteString("\n")
		}
	}

	return output.String()
}

// renderDialog renders the start entry dialog
func (m *Model) renderDialog() string {
	// Create title
	title := "Start New Entry"
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))

	// Create project input section
	projectLabel := "Project:"
	projectInput := m.inputs[0].View()
	projectAutocomplete := m.renderAutocompleteList(0)

	// Create title input section
	titleLabel := "Title:"
	titleInput := m.inputs[1].View()
	titleAutocomplete := m.renderAutocompleteList(1)

	// Create time input section
	timeLabel := "Time (HH:MM):"
	hourInput := m.inputs[2].View()
	minuteInput := m.inputs[3].View()

	// Create help text
	helpText := "↑/↓ for suggestions • Enter to select/submit • Tab to switch fields • Esc to cancel"
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)

	// Create error text style (red)
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))

	// Build dialog content
	var dialog strings.Builder
	dialog.WriteString(titleStyle.Render(title) + "\n\n")
	dialog.WriteString(projectLabel + "\n")
	dialog.WriteString(projectInput)
	dialog.WriteString(projectAutocomplete)
	dialog.WriteString("\n")
	dialog.WriteString(titleLabel + "\n")
	dialog.WriteString(titleInput)
	dialog.WriteString(titleAutocomplete)
	dialog.WriteString("\n")
	dialog.WriteString(timeLabel + "\n")
	dialog.WriteString(hourInput + " : " + minuteInput + "\n\n")

	// Show status/error message if present
	if m.status != "" {
		dialog.WriteString(errorStyle.Render(m.status) + "\n\n")
	}

	dialog.WriteString(helpStyle.Render(helpText) + "\n")

	return dialog.String()
}
