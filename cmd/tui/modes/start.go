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
		{Keys: "Tab", Label: "NEXT", Description: "Next field"},
		{Keys: "Shift+Tab", Label: "PREV", Description: "Previous field"},
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
			project := m.Inputs[0].Value()
			title := m.Inputs[1].Value()
			hourStr := m.Inputs[2].Value()
			minuteStr := m.Inputs[3].Value()

			if project == "" || title == "" {
				m.Status = "Project and title are required"
				return m, nil
			}

			if hourStr == "" {
				hourStr = "00"
			}
			if minuteStr == "" {
				minuteStr = "00"
			}

			var hour, minute int
			if n, err := fmt.Sscanf(hourStr, "%d", &hour); err != nil || n != 1 || hour < 0 || hour > 23 {
				m.Status = "Invalid hour (0-23)"
				return m, nil
			}
			if n, err := fmt.Sscanf(minuteStr, "%d", &minute); err != nil || n != 1 || minute < 0 || minute > 59 {
				m.Status = "Invalid minute (0-59)"
				return m, nil
			}

			now := time.Now()
			date := now

			if hour > now.Hour() || (hour == now.Hour() && minute > now.Minute()) {
				date = now.AddDate(0, 0, -1)
			}

			startTime := time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, date.Location())

			if _, err := m.TaskManager.StartEntryAt(project, title, startTime); err != nil {
				m.Status = "Error starting entry: " + err.Error()
			} else {
				m.Status = "Entry started: " + project
			}

			if err := m.LoadEntries(); err != nil {
				m.Err = err
			}
			m.CurrentMode = m.ListMode
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
		projectInput := m.Inputs[0].View()
		titleLabel := m.Styles.Label.Render("Title:")
		titleInput := m.Inputs[1].View()
		timeLabel := m.Styles.Label.Render("Time (HH:MM):")
		hourInput := m.Inputs[2].View()
		minuteInput := m.Inputs[3].View()

		var content strings.Builder
		content.WriteString(m.Styles.Title.Render(title) + "\n\n")
		content.WriteString(projectLabel + "\n")
		content.WriteString(projectInput + "\n\n")
		content.WriteString(titleLabel + "\n")
		content.WriteString(titleInput + "\n\n")
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
	m.FocusIndex = 0

	// Pre-fill the inputs with the selected entry's values
	m.Inputs[0].SetValue(entry.Project)
	m.Inputs[1].SetValue(entry.Title)

	// Set current time as default
	now := time.Now()
	m.Inputs[2].SetValue(fmt.Sprintf("%02d", now.Hour()))
	m.Inputs[3].SetValue(fmt.Sprintf("%02d", now.Minute()))

	// Set focus to first input
	m.Inputs[0].Focus()
	m.Inputs[0].PromptStyle = m.Styles.InputFocused
	m.Inputs[0].TextStyle = m.Styles.InputFocused

	// Blur other inputs
	for i := 1; i < len(m.Inputs); i++ {
		m.Inputs[i].Blur()
		m.Inputs[i].PromptStyle = m.Styles.InputBlurred
		m.Inputs[i].TextStyle = m.Styles.InputBlurred
	}
}

// openStartModeBlank opens start mode with blank values
func openStartModeBlank(m *Model) {
	m.CurrentMode = m.StartMode
	m.FocusIndex = 0

	// Clear all inputs
	for i := range m.Inputs {
		m.Inputs[i].SetValue("")
	}

	// Set current time as default
	now := time.Now()
	m.Inputs[2].SetValue(fmt.Sprintf("%02d", now.Hour()))
	m.Inputs[3].SetValue(fmt.Sprintf("%02d", now.Minute()))

	// Set focus to first input
	m.Inputs[0].Focus()
	m.Inputs[0].PromptStyle = m.Styles.InputFocused
	m.Inputs[0].TextStyle = m.Styles.InputFocused

	// Blur other inputs
	for i := 1; i < len(m.Inputs); i++ {
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
	projectInput := m.Inputs[0].View()

	// Create title input section
	titleLabel := m.Styles.Label.Render("Title:")
	titleInput := m.Inputs[1].View()

	// Create time input section
	timeLabel := m.Styles.Label.Render("Time (HH:MM):")
	hourInput := m.Inputs[2].View()
	minuteInput := m.Inputs[3].View()

	// Build content
	var content strings.Builder
	content.WriteString(m.Styles.Title.Render(title) + "\n\n")
	content.WriteString(projectLabel + "\n")
	content.WriteString(projectInput + "\n\n")
	content.WriteString(titleLabel + "\n")
	content.WriteString(titleInput + "\n\n")
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
