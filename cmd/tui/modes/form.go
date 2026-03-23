package modes

import (
	"fmt"
	"sort"
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

var formKeyBindings = []KeyBinding{
	{Keys: "Tab", Label: "NEXT", Description: "Next field"},
	{Keys: "Shift+Tab", Label: "PREV", Description: "Previous field"},
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

// openNewMode opens the form in new entry mode with empty fields
func openNewMode(m *Model) {
	m.CurrentMode = m.NewMode
	m.FormState = FormState{Mode: FormModeNew}
	m.Status = ""
	setProjectSuggestions(m)

	// Clear all inputs
	for i := range m.Inputs {
		m.Inputs[i].SetValue("")
	}

	// Set current date/time as default
	setCurrentDateTimeDefaults(m, time.Now())

	setupFormInputs(m)
}

// openEditMode opens the form in edit mode with all fields pre-filled
func openEditMode(m *Model, entry models.TimeEntry, idx int) {
	m.CurrentMode = m.EditMode
	m.FormState = FormState{Mode: FormModeEdit, EditingIdx: idx}
	m.Status = ""
	setProjectSuggestions(m)

	// Pre-fill all inputs from entry
	m.Inputs[InputProject].SetValue(entry.Project)
	m.Inputs[InputTitle].SetValue(entry.Title)
	m.Inputs[InputHour].SetValue(fmt.Sprintf("%02d", entry.Start.Hour()))
	m.Inputs[InputMinute].SetValue(fmt.Sprintf("%02d", entry.Start.Minute()))
	setDateDefaults(m, entry.Start)

	setupFormInputs(m)
}

// openResumeMode opens the form in resume mode with project/task pre-filled
func openResumeMode(m *Model, entry models.TimeEntry) {
	m.CurrentMode = m.ResumeMode
	m.FormState = FormState{Mode: FormModeResume}
	m.Status = ""
	setProjectSuggestions(m)

	// Pre-fill project and title from entry
	m.Inputs[InputProject].SetValue(entry.Project)
	m.Inputs[InputTitle].SetValue(entry.Title)

	// Set current date/time as default
	setCurrentDateTimeDefaults(m, time.Now())

	setupFormInputs(m)
}

func setProjectSuggestions(m *Model) {
	if m.Storage == nil || len(m.Inputs) <= InputProject {
		return
	}

	projects, err := m.Storage.LoadProjects()
	if err != nil {
		return
	}

	projects = normalizeProjectsForSuggestions(projects)

	suggestions := make([]string, 0, len(projects))
	for _, project := range projects {
		name := strings.TrimSpace(project.Name)
		if name == "" {
			continue
		}
		suggestions = append(suggestions, name)
	}

	m.Inputs[InputProject].SetSuggestions(suggestions)
}

func normalizeProjectsForSuggestions(projects []models.Project) []models.Project {
	type projectGroup struct {
		name string
		key  string
	}

	seen := make(map[string]models.Project, len(projects))
	groups := make([]projectGroup, 0, len(projects))

	for _, project := range projects {
		name := strings.TrimSpace(project.Name)
		if name == "" {
			continue
		}
		key := strings.ToLower(name)
		if _, ok := seen[key]; ok {
			continue
		}
		project.Name = name
		seen[key] = project
		groups = append(groups, projectGroup{name: name, key: key})
	}

	sort.SliceStable(groups, func(i, j int) bool {
		if groups[i].key != groups[j].key {
			return groups[i].key < groups[j].key
		}
		return groups[i].name < groups[j].name
	})

	normalized := make([]models.Project, 0, len(groups))
	for _, group := range groups {
		normalized = append(normalized, seen[group.key])
	}
	return normalized
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

// handleFormSubmit handles form submission for new/edit/resume modes
func handleFormSubmit(m *Model, formMode FormMode) (*Model, tea.Cmd) {
	project := m.Inputs[InputProject].Value()
	title := m.Inputs[InputTitle].Value()

	// For new/resume modes, require project and title
	// For edit mode, allow empty values (to create blank entries/gaps)
	if (formMode == FormModeNew || formMode == FormModeResume) && (project == "" || title == "") {
		m.Status = "Project and title are required"
		return m, nil
	}

	startTime, err := parseFormTime(m)
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

// handleFormKeyMsg handles common form key messages (navigation)
// Returns true if the key was handled, false otherwise
func handleFormKeyMsg(m *Model, msg tea.KeyMsg) (*Model, tea.Cmd, bool) {
	switch msg.String() {
	case "esc":
		m.CurrentMode = m.ListMode
		m.Status = ""
		return m, nil, true

	case "tab":
		if shouldPassTabToProjectSuggestions(m) {
			return m, nil, false
		}
		m.FocusIndex = (m.FocusIndex + 1) % len(m.Inputs)
		updateInputFocus(m)
		return m, nil, true

	case "shift+tab":
		m.FocusIndex--
		if m.FocusIndex < 0 {
			m.FocusIndex = len(m.Inputs) - 1
		}
		updateInputFocus(m)
		return m, nil, true
	}

	return m, nil, false
}

func shouldPassTabToProjectSuggestions(m *Model) bool {
	return m.FocusIndex == InputProject && len(m.Inputs[InputProject].MatchedSuggestions()) > 0
}

// parseFormTime parses date and time from form inputs.
func parseFormTime(m *Model) (time.Time, error) {
	yearStr := m.Inputs[InputYear].Value()
	monthStr := m.Inputs[InputMonth].Value()
	dayStr := m.Inputs[InputDay].Value()
	hourStr := m.Inputs[InputHour].Value()
	minuteStr := m.Inputs[InputMinute].Value()

	if hourStr == "" {
		hourStr = "00"
	}
	if minuteStr == "" {
		minuteStr = "00"
	}

	year, err := parseRangedInt(yearStr, "year", 1, 9999)
	if err != nil {
		return time.Time{}, err
	}
	month, err := parseRangedInt(monthStr, "month", 1, 12)
	if err != nil {
		return time.Time{}, err
	}
	day, err := parseRangedInt(dayStr, "day", 1, 31)
	if err != nil {
		return time.Time{}, err
	}
	hour, err := parseRangedInt(hourStr, "hour", 0, 23)
	if err != nil {
		return time.Time{}, err
	}
	minute, err := parseRangedInt(minuteStr, "minute", 0, 59)
	if err != nil {
		return time.Time{}, err
	}

	loc := time.Now().Location()
	startTime := time.Date(year, time.Month(month), day, hour, minute, 0, 0, loc)
	if startTime.Year() != year || int(startTime.Month()) != month || startTime.Day() != day {
		return time.Time{}, fmt.Errorf("invalid day for month/year")
	}

	return startTime, nil
}

func parseRangedInt(value, field string, min, max int) (int, error) {
	var parsed int
	if n, err := fmt.Sscanf(value, "%d", &parsed); err != nil || n != 1 || parsed < min || parsed > max {
		return 0, fmt.Errorf("invalid %s (%d-%d)", field, min, max)
	}

	return parsed, nil
}

// updateInputFocus updates the focus styling on all inputs
func updateInputFocus(m *Model) {
	for i := range m.Inputs {
		if i == m.FocusIndex {
			m.Inputs[i].Focus()
			m.Inputs[i].PromptStyle = m.Styles.InputFocused
			m.Inputs[i].TextStyle = m.Styles.InputFocused
		} else {
			m.Inputs[i].Blur()
			m.Inputs[i].PromptStyle = m.Styles.InputBlurred
			m.Inputs[i].TextStyle = m.Styles.InputBlurred
		}
	}
}

// setupFormInputs sets up form inputs with focus on first field
func setupFormInputs(m *Model) {
	m.FocusIndex = InputProject
	m.Inputs[InputProject].Focus()
	m.Inputs[InputProject].PromptStyle = m.Styles.InputFocused
	m.Inputs[InputProject].TextStyle = m.Styles.InputFocused

	for i := InputTitle; i < len(m.Inputs); i++ {
		m.Inputs[i].Blur()
		m.Inputs[i].PromptStyle = m.Styles.InputBlurred
		m.Inputs[i].TextStyle = m.Styles.InputBlurred
	}
}

func setDateDefaults(m *Model, date time.Time) {
	m.Inputs[InputYear].SetValue(fmt.Sprintf("%04d", date.Year()))
	m.Inputs[InputMonth].SetValue(fmt.Sprintf("%02d", int(date.Month())))
	m.Inputs[InputDay].SetValue(fmt.Sprintf("%02d", date.Day()))
}

func setCurrentDateTimeDefaults(m *Model, now time.Time) {
	m.Inputs[InputHour].SetValue(fmt.Sprintf("%02d", now.Hour()))
	m.Inputs[InputMinute].SetValue(fmt.Sprintf("%02d", now.Minute()))
	setDateDefaults(m, now)
}

func renderEntryFormBody(m *Model, content *strings.Builder) {
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
}

// renderFormContent renders the form with a given title
func renderFormContent(m *Model, title string, availableHeight int) string {
	_ = availableHeight // Available for future use

	var content strings.Builder
	content.WriteString(m.Styles.Title.Render(title) + "\n\n")
	renderEntryFormBody(m, &content)

	return content.String()
}
