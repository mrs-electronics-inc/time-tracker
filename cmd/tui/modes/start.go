package modes

import (
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
		return renderStartContent(m, availableHeight)
	},
}

// openStartMode opens start mode with pre-filled values from an entry
func openStartMode(m *Model, entry models.TimeEntry) {
	m.CurrentMode = m.StartMode
	m.FocusIndex = InputProject

	// Pre-fill the inputs with the selected entry's values
	m.Inputs[InputProject].SetValue(entry.Project)
	m.Inputs[InputTitle].SetValue(entry.Title)

	// Set current date/time as default
	setCurrentDateTimeDefaults(m, time.Now())

	setupFormInputs(m)
}

// openStartModeBlank opens start mode with blank values
func openStartModeBlank(m *Model) {
	m.CurrentMode = m.StartMode
	m.FocusIndex = InputProject

	// Clear all inputs
	for i := range m.Inputs {
		m.Inputs[i].SetValue("")
	}

	// Set current date/time as default
	setCurrentDateTimeDefaults(m, time.Now())

	setupFormInputs(m)
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

	var content strings.Builder
	content.WriteString(m.Styles.Title.Render("Start New Entry") + "\n\n")
	renderEntryFormBody(m, &content)

	return content.String()
}
