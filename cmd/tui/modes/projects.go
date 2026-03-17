package modes

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time-tracker/models"
)

// ProjectsMode is the project metadata list view mode.
var ProjectsMode = &Mode{
	Name: "projects",
	KeyBindings: []KeyBinding{
		{Keys: "Tab", Label: "LIST", Description: "Switch mode"},
		{Keys: "n", Label: "NEW", Description: "Add project"},
		{Keys: "e", Label: "EDIT", Description: "Edit project"},
		{Keys: "d", Label: "DELETE", Description: "Delete project"},
		{Keys: "k / ↑", Label: "UP", Description: "Scroll up"},
		{Keys: "j / ↓", Label: "DOWN", Description: "Scroll down"},
		{Keys: "?", Label: "HELP", Description: "Toggle help"},
		{Keys: "q / Esc", Label: "QUIT", Description: "Quit"},
	},
	HandleKeyMsg: func(m *Model, msg tea.KeyMsg) (*Model, tea.Cmd) {
		switch msg.String() {
		case "?":
			m.PreviousMode = m.CurrentMode
			m.CurrentMode = m.HelpMode
			return m, nil

		case "q", "esc":
			return m, tea.Quit

		case "tab":
			m.SwitchMode(m.ListMode)
			return m, nil

		case "n":
			openProjectNewMode(m)
			return m, nil

		case "e":
			if len(m.Projects) > 0 {
				projects := sortedProjectsByName(m.Projects)
				selected := clampSelectedProjectIndex(m, len(projects))
				openProjectEditMode(m, projects[selected])
			}
			return m, nil

		case "d":
			if len(m.Projects) == 0 {
				return m, nil
			}

			projects := sortedProjectsByName(m.Projects)
			selected := clampSelectedProjectIndex(m, len(projects))
			name := projects[selected].Name

			if err := m.TaskManager.RemoveProject(name); err != nil {
				m.Status = "Error removing project: " + err.Error()
				return m, nil
			}

			if err := m.LoadEntries(); err != nil {
				m.Err = err
				return m, nil
			}

			m.Status = "Project removed"
			setSelectedProjectByName(m, "")
			return m, nil

		case "k", "up":
			if m.SelectedIdx > 0 {
				m.SelectedIdx--
			}
			if m.ViewportTop > 0 {
				m.ViewportTop--
			}
			m.Status = ""
			return m, nil

		case "j", "down":
			if m.SelectedIdx < len(m.Projects)-1 {
				m.SelectedIdx++
			}
			maxRows := projectVisibleRows(m)
			maxTop := max(len(m.Projects)-maxRows, 0)
			if m.ViewportTop < maxTop {
				m.ViewportTop++
			}
			m.Status = ""
			return m, nil
		}

		return m, nil
	},
	RenderContent: func(m *Model, availableHeight int) string {
		return renderProjectsContent(m, availableHeight)
	},
}

func renderProjectsContent(m *Model, availableHeight int) string {
	if len(m.Projects) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)
		return emptyStyle.Render("No projects found. Press 'n' to add one.\n")
	}

	projects := sortedProjectsByName(m.Projects)
	selected := clampSelectedProjectIndex(m, len(projects))
	headerHeight := 2
	maxRows := max(availableHeight-headerHeight, 1)
	maxTop := max(len(projects)-maxRows, 0)
	if m.ViewportTop < 0 {
		m.ViewportTop = 0
	}
	if m.ViewportTop > maxTop {
		m.ViewportTop = maxTop
	}

	nameWidth := len("Name")
	codeWidth := len("Code")
	categoryWidth := len("Category")

	for _, project := range projects {
		if len(project.Name) > nameWidth {
			nameWidth = len(project.Name)
		}
		if len(project.Code) > codeWidth {
			codeWidth = len(project.Code)
		}
		if len(project.Category) > categoryWidth {
			categoryWidth = len(project.Category)
		}
	}

	nameWidth++
	codeWidth++

	headerText := fmt.Sprintf("%-*s %-*s %s", nameWidth, "Name", codeWidth, "Code", "Category")
	separatorText := strings.Repeat("-", max(lipgloss.Width(headerText), len("Name Code Category")))

	var output strings.Builder
	output.WriteString(m.Styles.Header.Render(headerText))
	output.WriteString("\n")
	output.WriteString(m.Styles.Header.Render(separatorText))
	output.WriteString("\n")

	end := min(m.ViewportTop+maxRows, len(projects))
	for i := m.ViewportTop; i < end; i++ {
		project := projects[i]
		row := fmt.Sprintf("%-*s %-*s %s", nameWidth, project.Name, codeWidth, project.Code, project.Category)
		if i == selected {
			output.WriteString(lipgloss.NewStyle().Bold(true).Reverse(true).Render(row))
		} else {
			output.WriteString(m.Styles.Unselected.Render(row))
		}
		output.WriteString("\n")
	}

	return output.String()
}

func sortedProjectsByName(projects []models.Project) []models.Project {
	sorted := append([]models.Project(nil), projects...)
	sort.SliceStable(sorted, func(i, j int) bool {
		left := strings.ToLower(sorted[i].Name)
		right := strings.ToLower(sorted[j].Name)
		if left == right {
			return sorted[i].Name < sorted[j].Name
		}
		return left < right
	})
	return sorted
}

func clampSelectedProjectIndex(m *Model, projectCount int) int {
	if projectCount <= 0 {
		m.SelectedIdx = 0
		return 0
	}
	if m.SelectedIdx < 0 {
		m.SelectedIdx = 0
	}
	if m.SelectedIdx >= projectCount {
		m.SelectedIdx = projectCount - 1
	}
	return m.SelectedIdx
}

func setSelectedProjectByName(m *Model, name string) {
	if len(m.Projects) == 0 {
		m.SelectedIdx = 0
		return
	}

	projects := sortedProjectsByName(m.Projects)
	if name == "" {
		m.SelectedIdx = clampSelectedProjectIndex(m, len(projects))
		return
	}

	for i, project := range projects {
		if strings.EqualFold(project.Name, name) {
			m.SelectedIdx = i
			return
		}
	}

	m.SelectedIdx = clampSelectedProjectIndex(m, len(projects))
}

func projectVisibleRows(m *Model) int {
	statusBarHeight := 1
	availableHeight := max(m.Height-statusBarHeight, 1)
	headerHeight := 2
	return max(availableHeight-headerHeight, 1)
}
