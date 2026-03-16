package modes

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ProjectsMode is the project metadata list view mode.
var ProjectsMode = &Mode{
	Name: "projects",
	KeyBindings: []KeyBinding{
		{Keys: "Tab", Label: "LIST", Description: "Switch mode"},
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
		}

		return m, nil
	},
	RenderContent: func(m *Model, availableHeight int) string {
		_ = availableHeight
		return renderProjectsContent(m)
	},
}

func renderProjectsContent(m *Model) string {
	if len(m.Projects) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)
		return emptyStyle.Render("No projects found.\n")
	}

	nameWidth := len("Name")
	codeWidth := len("Code")
	categoryWidth := len("Category")

	for _, project := range m.Projects {
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

	for _, project := range m.Projects {
		row := fmt.Sprintf("%-*s %-*s %s", nameWidth, project.Name, codeWidth, project.Code, project.Category)
		output.WriteString(m.Styles.Unselected.Render(row))
		output.WriteString("\n")
	}

	return output.String()
}
