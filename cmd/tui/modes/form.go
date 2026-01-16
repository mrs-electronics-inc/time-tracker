package modes

import (
	"fmt"
	"strings"
	"time"

	"time-tracker/models"

	tea "github.com/charmbracelet/bubbletea"
)

// FormMode represents the type of form operation
type FormMode int

const (
	FormModeNew FormMode = iota
	FormModeEdit
	FormModeResume
)

// FormState holds the state for form operations
type FormState struct {
	Mode       FormMode
	EditingIdx int // Index of entry being edited (for EditMode)
}

// renderFormContent renders the form with a given title
func renderFormContent(m *Model, title string, availableHeight int) string {
	_ = availableHeight // Available for future use

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
}

// handleFormKeyMsg handles common form key messages (navigation)
// Returns true if the key was handled, false otherwise
func handleFormKeyMsg(m *Model, msg tea.KeyMsg) (*Model, tea.Cmd, bool) {
	switch msg.String() {
	case "esc":
		m.CurrentMode = m.ListMode
		m.Status = ""
		return m, nil, true

	case "tab", "down":
		m.FocusIndex = (m.FocusIndex + 1) % len(m.Inputs)
		updateInputFocus(m)
		return m, nil, true

	case "shift+tab", "up":
		m.FocusIndex--
		if m.FocusIndex < 0 {
			m.FocusIndex = len(m.Inputs) - 1
		}
		updateInputFocus(m)
		return m, nil, true
	}

	return m, nil, false
}

// parseFormTime parses hour and minute from form inputs
// Returns the parsed time on today (or yesterday if time is in the future)
// If baseDate is provided, uses that as the date instead of today
func parseFormTime(m *Model, baseDate *time.Time) (time.Time, error) {
	hourStr := m.Inputs[2].Value()
	minuteStr := m.Inputs[3].Value()

	if hourStr == "" {
		hourStr = "00"
	}
	if minuteStr == "" {
		minuteStr = "00"
	}

	var hour, minute int
	if n, err := fmt.Sscanf(hourStr, "%d", &hour); err != nil || n != 1 || hour < 0 || hour > 23 {
		return time.Time{}, fmt.Errorf("invalid hour (0-23)")
	}
	if n, err := fmt.Sscanf(minuteStr, "%d", &minute); err != nil || n != 1 || minute < 0 || minute > 59 {
		return time.Time{}, fmt.Errorf("invalid minute (0-59)")
	}

	now := time.Now()
	date := now

	// If baseDate provided (edit mode), use its date
	if baseDate != nil {
		date = *baseDate
	} else {
		// If time is in the future, assume yesterday
		if hour > now.Hour() || (hour == now.Hour() && minute > now.Minute()) {
			date = now.AddDate(0, 0, -1)
		}
	}

	return time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, date.Location()), nil
}

// handleFormSubmit handles form submission for new/edit/resume modes
func handleFormSubmit(m *Model, formMode FormMode) (*Model, tea.Cmd) {
	project := m.Inputs[0].Value()
	title := m.Inputs[1].Value()

	if project == "" || title == "" {
		m.Status = "Project and title are required"
		return m, nil
	}

	// When editing, preserve the original entry's date
	var baseDate *time.Time
	if formMode == FormModeEdit && m.FormState.EditingIdx < len(m.Entries) {
		originalDate := m.Entries[m.FormState.EditingIdx].Start
		baseDate = &originalDate
	}

	startTime, err := parseFormTime(m, baseDate)
	if err != nil {
		m.Status = "Error: " + err.Error()
		return m, nil
	}

	switch formMode {
	case FormModeNew, FormModeResume:
		if _, err := m.TaskManager.StartEntryAt(project, title, startTime); err != nil {
			m.Status = "Error starting entry: " + err.Error()
		} else {
			m.Status = "Entry started: " + project
		}

	case FormModeEdit:
		if err := m.TaskManager.UpdateEntry(m.FormState.EditingIdx, project, title, startTime); err != nil {
			m.Status = "Error updating entry: " + err.Error()
		} else {
			m.Status = "Entry updated: " + project
		}
	}

	if err := m.LoadEntries(); err != nil {
		m.Err = err
	}
	m.CurrentMode = m.ListMode
	m.SelectMostRecentEntry()
	return m, nil
}

// setupFormInputs sets up form inputs with focus on first field
func setupFormInputs(m *Model) {
	m.FocusIndex = 0
	m.Inputs[0].Focus()
	m.Inputs[0].PromptStyle = m.Styles.InputFocused
	m.Inputs[0].TextStyle = m.Styles.InputFocused

	for i := 1; i < len(m.Inputs); i++ {
		m.Inputs[i].Blur()
		m.Inputs[i].PromptStyle = m.Styles.InputBlurred
		m.Inputs[i].TextStyle = m.Styles.InputBlurred
	}
}

// openNewMode opens the form in new entry mode with empty fields
func openNewMode(m *Model) {
	m.CurrentMode = m.NewMode
	m.FormState = FormState{Mode: FormModeNew}
	m.Status = ""

	// Clear all inputs
	for i := range m.Inputs {
		m.Inputs[i].SetValue("")
	}

	// Set current time as default
	now := time.Now()
	m.Inputs[2].SetValue(fmt.Sprintf("%02d", now.Hour()))
	m.Inputs[3].SetValue(fmt.Sprintf("%02d", now.Minute()))

	setupFormInputs(m)
}

// openEditMode opens the form in edit mode with all fields pre-filled
func openEditMode(m *Model, entry models.TimeEntry, idx int) {
	m.CurrentMode = m.EditMode
	m.FormState = FormState{Mode: FormModeEdit, EditingIdx: idx}
	m.Status = ""

	// Pre-fill all inputs from entry
	m.Inputs[0].SetValue(entry.Project)
	m.Inputs[1].SetValue(entry.Title)
	m.Inputs[2].SetValue(fmt.Sprintf("%02d", entry.Start.Hour()))
	m.Inputs[3].SetValue(fmt.Sprintf("%02d", entry.Start.Minute()))

	setupFormInputs(m)
}

// openResumeMode opens the form in resume mode with project/task pre-filled
func openResumeMode(m *Model, entry models.TimeEntry) {
	m.CurrentMode = m.ResumeMode
	m.FormState = FormState{Mode: FormModeResume}
	m.Status = ""

	// Pre-fill project and title from entry
	m.Inputs[0].SetValue(entry.Project)
	m.Inputs[1].SetValue(entry.Title)

	// Set current time as default
	now := time.Now()
	m.Inputs[2].SetValue(fmt.Sprintf("%02d", now.Hour()))
	m.Inputs[3].SetValue(fmt.Sprintf("%02d", now.Minute()))

	setupFormInputs(m)
}

// createFormKeyHandler creates a key handler for a form mode
func createFormKeyHandler(formMode FormMode) func(*Model, tea.KeyMsg) (*Model, tea.Cmd) {
	return func(m *Model, msg tea.KeyMsg) (*Model, tea.Cmd) {
		// Handle common form keys
		if model, cmd, handled := handleFormKeyMsg(m, msg); handled {
			return model, cmd
		}

		// Handle enter for submit
		if msg.String() == "enter" {
			return handleFormSubmit(m, formMode)
		}

		// Pass through to text inputs
		cmds := make([]tea.Cmd, len(m.Inputs))
		for i := range m.Inputs {
			m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
		}
		return m, tea.Batch(cmds...)
	}
}

var formKeyBindings = []KeyBinding{
	{Keys: "Tab / ↓", Label: "NEXT", Description: "Next field"},
	{Keys: "Shift+Tab / ↑", Label: "PREV", Description: "Previous field"},
	{Keys: "Enter", Label: "SUBMIT", Description: "Submit entry"},
	{Keys: "Esc", Label: "CANCEL", Description: "Cancel"},
}

// NewMode is the new entry form mode
var NewMode = &Mode{
	Name:         "new",
	KeyBindings:  formKeyBindings,
	HandleKeyMsg: createFormKeyHandler(FormModeNew),
	RenderContent: func(m *Model, availableHeight int) string {
		return renderFormContent(m, "New Entry", availableHeight)
	},
}

// EditMode is the edit entry form mode
var EditMode = &Mode{
	Name:         "edit",
	KeyBindings:  formKeyBindings,
	HandleKeyMsg: createFormKeyHandler(FormModeEdit),
	RenderContent: func(m *Model, availableHeight int) string {
		return renderFormContent(m, "Edit Entry", availableHeight)
	},
}

// ResumeMode is the resume entry form mode
var ResumeMode = &Mode{
	Name:         "resume",
	KeyBindings:  formKeyBindings,
	HandleKeyMsg: createFormKeyHandler(FormModeResume),
	RenderContent: func(m *Model, availableHeight int) string {
		return renderFormContent(m, "Resume Entry", availableHeight)
	},
}
