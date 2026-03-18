package modes

import (
	"fmt"
	"strings"
	"time"

	"time-tracker/models"

	tea "github.com/charmbracelet/bubbletea"
)

// StartMode is the entry creation/editing mode
var StartMode = &Mode{
	Name: "start",
	KeyBindings: []KeyBinding{
		{Keys: "Tab / ↓", Label: "NEXT", Description: "Next field"},
		{Keys: "Shift+Tab / ↑", Label: "PREV", Description: "Previous field"},
		{Keys: "Enter", Label: "SUBMIT", Description: "Submit entry"},
		{Keys: "Esc", Label: "CANCEL", Description: "Cancel"},
	},
	HandleKeyMsg: func(m *Model, msg tea.KeyMsg) (*Model, tea.Cmd) {
		switch msg.String() {
		case "esc":
			m.CurrentMode = m.ListMode
			return m, nil

		case "tab", "down":
			m.FocusIndex = (m.FocusIndex + 1) % len(m.Inputs)
			updateInputFocus(m)
			return m, nil

		case "shift+tab", "up":
			m.FocusIndex--
			if m.FocusIndex < 0 {
				m.FocusIndex = len(m.Inputs) - 1
			}
			updateInputFocus(m)
			return m, nil

		case "enter":
			project := m.Inputs[InputProject].Value()
			title := m.Inputs[InputTitle].Value()

			if project == "" || title == "" {
				m.Status = "Project and title are required"
				return m, nil
			}

			startTime, err := parseFormTime(m)
			if err != nil {
				m.Status = "Invalid date/time: " + err.Error()
				return m, nil
			}

			if _, err := m.TaskManager.StartEntryAt(project, title, startTime); err != nil {
				m.Status = "Error starting entry: " + err.Error()
			} else {
				m.Status = "Entry started: " + project
			}

			if err := m.LoadEntries(); err != nil {
				m.Err = err
			}
			m.CurrentMode = m.ListMode
			m.SelectMostRecentEntry()
			return m, nil
		}

		cmds := make([]tea.Cmd, len(m.Inputs))
		for i := range m.Inputs {
			m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
		}
		return m, tea.Batch(cmds...)
	},
	RenderContent: func(m *Model, availableHeight int) string {
		_ = availableHeight

		title := "Start New Entry"
		projectLabel := m.Styles.Label.Render("Project:")
		projectInput := m.Inputs[InputProject].View()
		titleLabel := m.Styles.Label.Render("Title:")
		titleInput := m.Inputs[InputTitle].View()
		dateLabel := m.Styles.Label.Render("Date (YYYY-MM-DD):")
		yearInput := m.Inputs[InputYear].View()
		monthInput := m.Inputs[InputMonth].View()
		dayInput := m.Inputs[InputDay].View()
		timeLabel := m.Styles.Label.Render("Time (HH:MM):")
		hourInput := m.Inputs[InputHour].View()
		minuteInput := m.Inputs[InputMinute].View()

		var content strings.Builder
		content.WriteString(m.Styles.Title.Render(title) + "\n\n")
		content.WriteString(projectLabel + "\n")
		content.WriteString(projectInput + "\n\n")
		content.WriteString(titleLabel + "\n")
		content.WriteString(titleInput + "\n\n")
		content.WriteString(dateLabel + "\n")
		content.WriteString(yearInput + " - " + monthInput + " - " + dayInput + "\n\n")
		content.WriteString(timeLabel + "\n")
		content.WriteString(hourInput + " : " + minuteInput + "\n\n")

		if m.Status != "" {
			if strings.Contains(strings.ToLower(m.Status), "error") {
				content.WriteString(m.Styles.StatusError.Render(m.Status) + "\n\n")
			} else {
				content.WriteString(m.Styles.StatusSuccess.Render(m.Status) + "\n\n")
			}
		}

		return content.String()
	},
}

// openStartMode opens start mode with pre-filled values from an entry
func openStartMode(m *Model, entry models.TimeEntry) {
	m.CurrentMode = m.StartMode
	m.FocusIndex = InputProject

	// Pre-fill the inputs with the selected entry's values
	m.Inputs[InputProject].SetValue(entry.Project)
	m.Inputs[InputTitle].SetValue(entry.Title)

	// Set current time as default
	now := time.Now()
	m.Inputs[InputHour].SetValue(fmt.Sprintf("%02d", now.Hour()))
	m.Inputs[InputMinute].SetValue(fmt.Sprintf("%02d", now.Minute()))
	setDateDefaults(m, now)

	// Set focus to first input
	m.Inputs[InputProject].Focus()
	m.Inputs[InputProject].PromptStyle = m.Styles.InputFocused
	m.Inputs[InputProject].TextStyle = m.Styles.InputFocused

	// Blur other inputs
	for i := InputTitle; i < len(m.Inputs); i++ {
		m.Inputs[i].Blur()
		m.Inputs[i].PromptStyle = m.Styles.InputBlurred
		m.Inputs[i].TextStyle = m.Styles.InputBlurred
	}
}

// openStartModeBlank opens start mode with blank values
func openStartModeBlank(m *Model) {
	m.CurrentMode = m.StartMode
	m.FocusIndex = InputProject

	// Clear all inputs
	for i := range m.Inputs {
		m.Inputs[i].SetValue("")
	}

	// Set current time as default
	now := time.Now()
	m.Inputs[InputHour].SetValue(fmt.Sprintf("%02d", now.Hour()))
	m.Inputs[InputMinute].SetValue(fmt.Sprintf("%02d", now.Minute()))
	setDateDefaults(m, now)

	// Set focus to first input
	m.Inputs[InputProject].Focus()
	m.Inputs[InputProject].PromptStyle = m.Styles.InputFocused
	m.Inputs[InputProject].TextStyle = m.Styles.InputFocused

	// Blur other inputs
	for i := InputTitle; i < len(m.Inputs); i++ {
		m.Inputs[i].Blur()
		m.Inputs[i].PromptStyle = m.Styles.InputBlurred
		m.Inputs[i].TextStyle = m.Styles.InputBlurred
	}
}

// updateInputFocus updates the focus styling on all inputs
func updateInputFocus(m *Model) {
	for i := range m.Inputs {
		if i == m.FocusIndex {
			// Set focused state
			m.Inputs[i].Focus()
			m.Inputs[i].PromptStyle = m.Styles.InputFocused
			m.Inputs[i].TextStyle = m.Styles.InputFocused
		} else {
			// Set blurred state
			m.Inputs[i].Blur()
			m.Inputs[i].PromptStyle = m.Styles.InputBlurred
			m.Inputs[i].TextStyle = m.Styles.InputBlurred
		}
	}
}

// renderStartContent renders the start mode content
func renderStartContent(m *Model, availableHeight int) string {
	_ = availableHeight // Available for future use
	// Create title
	title := "Start New Entry"

	// Create project input section
	projectLabel := m.Styles.Label.Render("Project:")
	projectInput := m.Inputs[InputProject].View()

	// Create title input section
	titleLabel := m.Styles.Label.Render("Title:")
	titleInput := m.Inputs[InputTitle].View()

	// Create date input section
	dateLabel := m.Styles.Label.Render("Date (YYYY-MM-DD):")
	yearInput := m.Inputs[InputYear].View()
	monthInput := m.Inputs[InputMonth].View()
	dayInput := m.Inputs[InputDay].View()

	// Create time input section
	timeLabel := m.Styles.Label.Render("Time (HH:MM):")
	hourInput := m.Inputs[InputHour].View()
	minuteInput := m.Inputs[InputMinute].View()

	// Build content
	var content strings.Builder
	content.WriteString(m.Styles.Title.Render(title) + "\n\n")
	content.WriteString(projectLabel + "\n")
	content.WriteString(projectInput + "\n\n")
	content.WriteString(titleLabel + "\n")
	content.WriteString(titleInput + "\n\n")
	content.WriteString(dateLabel + "\n")
	content.WriteString(yearInput + " - " + monthInput + " - " + dayInput + "\n\n")
	content.WriteString(timeLabel + "\n")
	content.WriteString(hourInput + " : " + minuteInput + "\n\n")

	// Show status/error message if present
	if m.Status != "" {
		// Determine if it's an error or success based on message content
		if strings.Contains(strings.ToLower(m.Status), "error") {
			content.WriteString(m.Styles.StatusError.Render(m.Status) + "\n\n")
		} else {
			content.WriteString(m.Styles.StatusSuccess.Render(m.Status) + "\n\n")
		}
	}

	return content.String()
}
