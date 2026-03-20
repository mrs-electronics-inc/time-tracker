package modes

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"time-tracker/models"
)

type ProjectFormMode int

const (
	ProjectFormModeNew ProjectFormMode = iota
	ProjectFormModeEdit
)

type ProjectFormState struct {
	Mode        ProjectFormMode
	EditingName string
}

var projectFormKeyBindings = []KeyBinding{
	{Keys: "Tab", Label: "NEXT", Description: "Next field"},
	{Keys: "Shift+Tab", Label: "PREV", Description: "Previous field"},
	{Keys: "Enter", Label: "SUBMIT", Description: "Save project"},
	{Keys: "Esc", Label: "CANCEL", Description: "Cancel"},
}

var ProjectNewMode = &Mode{
	Name:         "project-new",
	KeyBindings:  projectFormKeyBindings,
	HandleKeyMsg: createProjectFormKeyHandler(ProjectFormModeNew),
	RenderContent: func(m *Model, availableHeight int) string {
		return renderProjectFormContent(m, "New Project", availableHeight)
	},
}

var ProjectEditMode = &Mode{
	Name:         "project-edit",
	KeyBindings:  projectFormKeyBindings,
	HandleKeyMsg: createProjectFormKeyHandler(ProjectFormModeEdit),
	RenderContent: func(m *Model, availableHeight int) string {
		return renderProjectFormContent(m, "Edit Project", availableHeight)
	},
}

func openProjectNewMode(m *Model) {
	m.CurrentMode = m.ProjectNewMode
	m.ProjectFormState = ProjectFormState{Mode: ProjectFormModeNew}
	m.Status = ""

	for i := range m.ProjectInputs {
		m.ProjectInputs[i].SetValue("")
	}

	setupProjectFormInputs(m)
}

func openProjectEditMode(m *Model, project models.Project) {
	m.CurrentMode = m.ProjectEditMode
	m.ProjectFormState = ProjectFormState{Mode: ProjectFormModeEdit, EditingName: project.Name}
	m.Status = ""

	m.ProjectInputs[0].SetValue(project.Name)
	m.ProjectInputs[1].SetValue(project.Code)
	m.ProjectInputs[2].SetValue(project.Category)

	setupProjectFormInputs(m)
}

func createProjectFormKeyHandler(formMode ProjectFormMode) func(*Model, tea.KeyMsg) (*Model, tea.Cmd) {
	return func(m *Model, msg tea.KeyMsg) (*Model, tea.Cmd) {
		if model, cmd, handled := handleProjectFormKeyMsg(m, msg); handled {
			return model, cmd
		}

		if msg.String() == "enter" {
			return handleProjectFormSubmit(m, formMode)
		}

		cmds := make([]tea.Cmd, len(m.ProjectInputs))
		for i := range m.ProjectInputs {
			m.ProjectInputs[i], cmds[i] = m.ProjectInputs[i].Update(msg)
		}

		return m, tea.Batch(cmds...)
	}
}

func handleProjectFormKeyMsg(m *Model, msg tea.KeyMsg) (*Model, tea.Cmd, bool) {
	switch msg.String() {
	case "esc":
		m.CurrentMode = m.ProjectsMode
		m.Status = ""
		return m, nil, true

	case "tab":
		m.ProjectFocusIndex = (m.ProjectFocusIndex + 1) % len(m.ProjectInputs)
		updateProjectInputFocus(m)
		return m, nil, true

	case "shift+tab":
		m.ProjectFocusIndex--
		if m.ProjectFocusIndex < 0 {
			m.ProjectFocusIndex = len(m.ProjectInputs) - 1
		}
		updateProjectInputFocus(m)
		return m, nil, true
	}

	return m, nil, false
}

func handleProjectFormSubmit(m *Model, formMode ProjectFormMode) (*Model, tea.Cmd) {
	name := m.ProjectInputs[0].Value()
	code := m.ProjectInputs[1].Value()
	category := m.ProjectInputs[2].Value()
	trimmedName := strings.TrimSpace(name)

	if trimmedName == "" {
		if formMode == ProjectFormModeNew {
			m.Status = "Error adding project: project name cannot be empty"
		} else {
			m.Status = "Error editing project: project name cannot be empty"
		}
		return m, nil
	}

	switch formMode {
	case ProjectFormModeNew:
		if _, err := m.TaskManager.AddProject(name, code, category); err != nil {
			m.Status = "Error adding project: " + err.Error()
			return m, nil
		}
		m.Status = "Project added"

	case ProjectFormModeEdit:
		if _, err := m.TaskManager.EditProject(m.ProjectFormState.EditingName, name, code, category); err != nil {
			m.Status = "Error editing project: " + err.Error()
			return m, nil
		}
		m.Status = "Project updated"
	}

	if err := m.LoadEntries(); err != nil {
		m.Err = err
	}

	m.CurrentMode = m.ProjectsMode
	setSelectedProjectByName(m, trimmedName)
	return m, nil
}

func setupProjectFormInputs(m *Model) {
	m.ProjectFocusIndex = 0
	m.ProjectInputs[0].Focus()
	m.ProjectInputs[0].PromptStyle = m.Styles.InputFocused
	m.ProjectInputs[0].TextStyle = m.Styles.InputFocused

	for i := 1; i < len(m.ProjectInputs); i++ {
		m.ProjectInputs[i].Blur()
		m.ProjectInputs[i].PromptStyle = m.Styles.InputBlurred
		m.ProjectInputs[i].TextStyle = m.Styles.InputBlurred
	}
}

func updateProjectInputFocus(m *Model) {
	for i := range m.ProjectInputs {
		if i == m.ProjectFocusIndex {
			m.ProjectInputs[i].Focus()
			m.ProjectInputs[i].PromptStyle = m.Styles.InputFocused
			m.ProjectInputs[i].TextStyle = m.Styles.InputFocused
		} else {
			m.ProjectInputs[i].Blur()
			m.ProjectInputs[i].PromptStyle = m.Styles.InputBlurred
			m.ProjectInputs[i].TextStyle = m.Styles.InputBlurred
		}
	}
}

func renderProjectFormContent(m *Model, title string, availableHeight int) string {
	_ = availableHeight

	nameLabel := m.Styles.Label.Render("Name:")
	nameInput := m.ProjectInputs[0].View()
	codeLabel := m.Styles.Label.Render("Code:")
	codeInput := m.ProjectInputs[1].View()
	categoryLabel := m.Styles.Label.Render("Category:")
	categoryInput := m.ProjectInputs[2].View()

	var content strings.Builder
	content.WriteString(m.Styles.Title.Render(title) + "\n\n")
	content.WriteString(nameLabel + "\n")
	content.WriteString(nameInput + "\n\n")
	content.WriteString(codeLabel + "\n")
	content.WriteString(codeInput + "\n\n")
	content.WriteString(categoryLabel + "\n")
	content.WriteString(categoryInput + "\n\n")

	if m.Status != "" {
		if strings.Contains(strings.ToLower(m.Status), "error") {
			content.WriteString(m.Styles.StatusError.Render(m.Status) + "\n\n")
		} else {
			content.WriteString(m.Styles.StatusSuccess.Render(m.Status) + "\n\n")
		}
	}

	return content.String()
}
