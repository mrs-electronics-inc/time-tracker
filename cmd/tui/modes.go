package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time-tracker/models"
)

// openStartMode opens start mode with pre-filled values from an entry
func (m *Model) openStartMode(entry models.TimeEntry) {
	m.prevMode = ModeList
	m.mode = ModeStart
	m.focusIndex = 0

	// Pre-fill the inputs with the selected entry's values
	m.inputs[0].SetValue(entry.Project)
	m.inputs[1].SetValue(entry.Title)

	// Set current time as default
	now := time.Now()
	m.inputs[2].SetValue(fmt.Sprintf("%02d", now.Hour()))
	m.inputs[3].SetValue(fmt.Sprintf("%02d", now.Minute()))

	// Set focus to first input
	m.inputs[0].Focus()
	m.inputs[0].PromptStyle = m.styles.inputFocused
	m.inputs[0].TextStyle = m.styles.inputFocused

	// Blur other inputs
	for i := 1; i < len(m.inputs); i++ {
		m.inputs[i].Blur()
		m.inputs[i].PromptStyle = m.styles.inputBlurred
		m.inputs[i].TextStyle = m.styles.inputBlurred
	}
}

// openStartModeBlank opens start mode with blank values
func (m *Model) openStartModeBlank() {
	m.prevMode = ModeList
	m.mode = ModeStart
	m.focusIndex = 0

	// Clear all inputs
	for i := range m.inputs {
		m.inputs[i].SetValue("")
	}

	// Set current time as default
	now := time.Now()
	m.inputs[2].SetValue(fmt.Sprintf("%02d", now.Hour()))
	m.inputs[3].SetValue(fmt.Sprintf("%02d", now.Minute()))

	// Set focus to first input
	m.inputs[0].Focus()
	m.inputs[0].PromptStyle = m.styles.inputFocused
	m.inputs[0].TextStyle = m.styles.inputFocused

	// Blur other inputs
	for i := 1; i < len(m.inputs); i++ {
		m.inputs[i].Blur()
		m.inputs[i].PromptStyle = m.styles.inputBlurred
		m.inputs[i].TextStyle = m.styles.inputBlurred
	}
}

// closeStartMode closes start mode and returns to list mode
func (m *Model) closeStartMode() {
	m.mode = m.prevMode
	m.focusIndex = 0

	// Clear and blur all inputs
	for i := range m.inputs {
		m.inputs[i].SetValue("")
		m.inputs[i].Blur()
	}
}

// handleStartKeyMsg handles key messages while in start mode
func (m *Model) handleStartKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Cancel and return to list
		m.closeStartMode()
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
		// Submit entry
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

		// Reload entries and return to list
		if err := m.LoadEntries(); err != nil {
			m.err = err
		}
		m.closeStartMode()
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
			m.inputs[i].PromptStyle = m.styles.inputFocused
			m.inputs[i].TextStyle = m.styles.inputFocused
		} else {
			// Set blurred state
			m.inputs[i].Blur()
			m.inputs[i].PromptStyle = m.styles.inputBlurred
			m.inputs[i].TextStyle = m.styles.inputBlurred
		}
	}
}

// handleHelpKeyMsg handles key messages while in help mode
func (m *Model) handleHelpKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "?", "q":
		// Return to previous mode
		m.mode = m.prevMode
		return m, nil
	}
	return m, nil
}

// renderHelpScreen renders the help screen showing keybindings for the previous mode
func (m *Model) renderHelpScreen() string {
	var content strings.Builder

	title := m.styles.title.Render("Keyboard Shortcuts")
	content.WriteString(title + "\n\n")

	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))

	// Show keybindings based on the previous mode
	switch m.prevMode {
	case ModeStart:
		content.WriteString(keyStyle.Render("Tab / ↓") + descStyle.Render("  Next field") + "\n")
		content.WriteString(keyStyle.Render("Shift+Tab / ↑") + descStyle.Render("  Previous field") + "\n")
		content.WriteString(keyStyle.Render("Enter") + descStyle.Render("  Submit entry") + "\n")
		content.WriteString(keyStyle.Render("Esc") + descStyle.Render("  Cancel") + "\n")
	default: // ModeList
		content.WriteString(keyStyle.Render("j / ↓") + descStyle.Render("  Move down") + "\n")
		content.WriteString(keyStyle.Render("k / ↑") + descStyle.Render("  Move up") + "\n")
		content.WriteString(keyStyle.Render("G") + descStyle.Render("  Jump to bottom") + "\n")
		content.WriteString(keyStyle.Render("s") + descStyle.Render("  Start/stop entry") + "\n")
		content.WriteString(keyStyle.Render("?") + descStyle.Render("  Toggle help") + "\n")
		content.WriteString(keyStyle.Render("q / Esc") + descStyle.Render("  Quit") + "\n")
	}

	// Calculate content height and add spacer to push footer to bottom
	contentHeight := strings.Count(content.String(), "\n")
	footer := m.renderStatusBar()
	spacerHeight := m.height - contentHeight - 1
	if spacerHeight > 0 {
		content.WriteString(strings.Repeat("\n", spacerHeight-1))
	}
	content.WriteString(footer)

	return content.String()
}

// renderStartScreen renders the start mode screen
func (m *Model) renderStartScreen() string {
	// Create title
	title := "Start New Entry"

	// Create project input section
	projectLabel := m.styles.label.Render("Project:")
	projectInput := m.inputs[0].View()

	// Create title input section
	titleLabel := m.styles.label.Render("Title:")
	titleInput := m.inputs[1].View()

	// Create time input section
	timeLabel := m.styles.label.Render("Time (HH:MM):")
	hourInput := m.inputs[2].View()
	minuteInput := m.inputs[3].View()

	// Build content
	var content strings.Builder
	content.WriteString(m.styles.title.Render(title) + "\n\n")
	content.WriteString(projectLabel + "\n")
	content.WriteString(projectInput + "\n\n")
	content.WriteString(titleLabel + "\n")
	content.WriteString(titleInput + "\n\n")
	content.WriteString(timeLabel + "\n")
	content.WriteString(hourInput + " : " + minuteInput + "\n\n")

	// Show status/error message if present
	if m.status != "" {
		// Determine if it's an error or success based on message content
		if strings.Contains(strings.ToLower(m.status), "error") {
			content.WriteString(m.styles.statusError.Render(m.status) + "\n\n")
		} else {
			content.WriteString(m.styles.statusSuccess.Render(m.status) + "\n\n")
		}
	}

	// Calculate content height and add spacer to push footer to bottom
	contentHeight := strings.Count(content.String(), "\n")
	footer := m.renderStatusBar()
	spacerHeight := m.height - contentHeight - 1
	if spacerHeight > 0 {
		content.WriteString(strings.Repeat("\n", spacerHeight-1))
	}
	content.WriteString(footer)

	return content.String()
}
