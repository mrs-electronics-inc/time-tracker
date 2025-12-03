package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time-tracker/models"
)

// openStartDialog opens the start entry dialog with pre-filled values
func (m *Model) openStartDialog(entry models.TimeEntry) {
	m.dialogMode = true
	m.focusIndex = 0

	// Pre-fill the inputs with the selected entry's values
	m.inputs[0].SetValue(entry.Project)
	m.inputs[1].SetValue(entry.Title)

	// Set focus to first input
	m.inputs[0].Focus()
	m.inputs[0].PromptStyle = m.styles.dialogFocused
	m.inputs[0].TextStyle = m.styles.dialogFocused

	// Blur second input
	m.inputs[1].Blur()
	m.inputs[1].PromptStyle = m.styles.dialogBlurred
	m.inputs[1].TextStyle = m.styles.dialogBlurred
}

// closeDialog closes the dialog and returns to list mode
func (m *Model) closeDialog() {
	m.dialogMode = false
	m.focusIndex = 0

	// Clear inputs
	m.inputs[0].SetValue("")
	m.inputs[1].SetValue("")

	// Blur all inputs
	m.inputs[0].Blur()
	m.inputs[1].Blur()
}

// handleDialogKeyMsg handles key messages while in dialog mode
func (m *Model) handleDialogKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Cancel dialog
		m.closeDialog()
		return m, nil

	case "tab", "down":
		// Move focus to next input
		m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
		m.updateInputFocus()
		return m, nil

	case "shift+tab", "up":
		// Move focus to previous input
		m.focusIndex--
		if m.focusIndex < 0 {
			m.focusIndex = len(m.inputs) - 1
		}
		m.updateInputFocus()
		return m, nil

	case "enter":
		// Submit dialog
		project := m.inputs[0].Value()
		title := m.inputs[1].Value()

		if project == "" && title == "" {
			m.status = "Project and title cannot both be empty"
			return m, nil
		}

		// Start the entry
		if _, err := m.taskManager.StartEntry(project, title); err != nil {
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

// renderDialog renders the start entry dialog
func (m *Model) renderDialog() string {
	// Create title
	title := "Start New Entry"
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))

	// Create project input section
	projectLabel := "Project:"
	projectInput := m.inputs[0].View()

	// Create title input section
	titleLabel := "Title:"
	titleInput := m.inputs[1].View()

	// Create help text
	helpText := "Tab/↓/↑ to switch fields • Enter to submit • Esc to cancel"
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)

	// Build dialog content
	var dialog strings.Builder
	dialog.WriteString(titleStyle.Render(title) + "\n\n")
	dialog.WriteString(projectLabel + "\n")
	dialog.WriteString(projectInput + "\n\n")
	dialog.WriteString(titleLabel + "\n")
	dialog.WriteString(titleInput + "\n\n")
	dialog.WriteString(helpStyle.Render(helpText) + "\n")

	return dialog.String()
}
