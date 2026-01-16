package modes

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time-tracker/models"
	"time-tracker/utils"
)

type ConfirmState struct {
	DeletingIdx int
}

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("1"))
	labelStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	valueStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)
	boldStyle    = lipgloss.NewStyle().Bold(true)
)

var ConfirmMode = &Mode{
	Name: "confirm",
	KeyBindings: []KeyBinding{
		{Keys: "y", Label: "YES", Description: "Confirm delete"},
		{Keys: "n / Esc", Label: "NO", Description: "Cancel"},
	},
	HandleKeyMsg: func(m *Model, msg tea.KeyMsg) (*Model, tea.Cmd) {
		switch msg.String() {
		case "y", "Y":
			if err := m.TaskManager.DeleteEntry(m.ConfirmState.DeletingIdx); err != nil {
				m.Status = "Error deleting entry: " + err.Error()
			} else {
				m.Status = "Entry deleted"
			}
			if err := m.LoadEntries(); err != nil {
				m.Err = err
			}
			m.CurrentMode = m.ListMode
			return m, nil
		case "n", "N", "esc":
			m.CurrentMode = m.ListMode
			m.Status = ""
			return m, nil
		}
		return m, nil
	},
	RenderContent: func(m *Model, _ int) string {
		idx := m.ConfirmState.DeletingIdx
		if idx < 0 || idx >= len(m.Entries) {
			return "Invalid entry"
		}
		entry := &m.Entries[idx]
		if entry.IsBlank() {
			return renderBlankConfirmDialog(entry)
		}
		return renderNonBlankConfirmDialog(entry)
	},
}

func renderBlankConfirmDialog(entry *models.TimeEntry) string {
	var out strings.Builder

	out.WriteString(titleStyle.Render("Remove Gap?") + "\n\n")
	out.WriteString(labelStyle.Render("Start: ") + valueStyle.Render(entry.Start.Format("2006-01-02 15:04")) + "\n")

	if entry.End != nil {
		out.WriteString(labelStyle.Render("End: ") + valueStyle.Render(entry.End.Format("2006-01-02 15:04")) + "\n")
	}

	writeConfirmFooter(&out, "This will remove the gap.")
	return out.String()
}

func renderNonBlankConfirmDialog(entry *models.TimeEntry) string {
	var out strings.Builder

	out.WriteString(titleStyle.Render("Delete Entry?") + "\n\n")
	out.WriteString(labelStyle.Render("Project: ") + valueStyle.Render(entry.Project) + "\n")
	out.WriteString(labelStyle.Render("Task: ") + valueStyle.Render(entry.Title) + "\n")
	out.WriteString(labelStyle.Render("Start: ") + valueStyle.Render(entry.Start.Format("2006-01-02 15:04")) + "\n")

	if entry.End != nil {
		out.WriteString(labelStyle.Render("End: ") + valueStyle.Render(entry.End.Format("2006-01-02 15:04")) + "\n")
	} else {
		out.WriteString(labelStyle.Render("End: ") + valueStyle.Render("running") + "\n")
	}

	out.WriteString(labelStyle.Render("Duration: ") + valueStyle.Render(utils.FormatDuration(entry.Duration())) + "\n")

	writeConfirmFooter(&out, "This will convert the entry to a gap.")
	return out.String()
}

func writeConfirmFooter(out *strings.Builder, warning string) {
	out.WriteString("\n" + warningStyle.Render(warning) + "\n")
	out.WriteString("\nPress " + boldStyle.Render("y") + " to confirm, " + boldStyle.Render("n") + " to cancel\n")
}

func openConfirmDelete(m *Model, idx int) {
	m.CurrentMode = m.ConfirmMode
	m.ConfirmState = ConfirmState{DeletingIdx: idx}
	m.Status = ""
}
