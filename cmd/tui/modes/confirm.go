package modes

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time-tracker/utils"
)

// ConfirmState holds the state for confirmation dialogs
type ConfirmState struct {
	DeletingIdx int  // Index of entry being deleted
	Confirmed   bool // Whether user selected yes
}

// ConfirmMode is the delete confirmation modal
var ConfirmMode = &Mode{
	Name: "confirm",
	KeyBindings: []KeyBinding{
		{Keys: "y", Label: "YES", Description: "Confirm delete"},
		{Keys: "n / Esc", Label: "NO", Description: "Cancel"},
	},
	HandleKeyMsg: func(m *Model, msg tea.KeyMsg) (*Model, tea.Cmd) {
		switch msg.String() {
		case "y", "Y":
			// Delete the entry (convert to blank)
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
	RenderContent: func(m *Model, availableHeight int) string {
		_ = availableHeight

		// Get the entry being deleted
		if m.ConfirmState.DeletingIdx < 0 || m.ConfirmState.DeletingIdx >= len(m.Entries) {
			return "Invalid entry"
		}
		entry := m.Entries[m.ConfirmState.DeletingIdx]

		// Styles
		titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("1")) // Red
		labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
		valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
		warningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)

		var content strings.Builder
		// Determine if this is a blank entry
		isBlank := entry.IsBlank()

		if isBlank {
			content.WriteString(titleStyle.Render("Remove Gap?") + "\n\n")
		} else {
			content.WriteString(titleStyle.Render("Delete Entry?") + "\n\n")
		}

		// Entry details (skip for blank entries)
		if !isBlank {
			content.WriteString(labelStyle.Render("Project: ") + valueStyle.Render(entry.Project) + "\n")
			content.WriteString(labelStyle.Render("Task: ") + valueStyle.Render(entry.Title) + "\n")
		}
		content.WriteString(labelStyle.Render("Start: ") + valueStyle.Render(entry.Start.Format("2006-01-02 15:04")) + "\n")
		if entry.End != nil {
			content.WriteString(labelStyle.Render("End: ") + valueStyle.Render(entry.End.Format("2006-01-02 15:04")) + "\n")
		} else if !isBlank {
			content.WriteString(labelStyle.Render("End: ") + valueStyle.Render("running") + "\n")
		}
		if !isBlank {
			content.WriteString(labelStyle.Render("Duration: ") + valueStyle.Render(utils.FormatDuration(entry.Duration())) + "\n\n")
		} else {
			content.WriteString("\n")
		}

		if isBlank {
			content.WriteString(warningStyle.Render("This will remove the gap.") + "\n\n")
		} else {
			content.WriteString(warningStyle.Render("This will convert the entry to a gap.") + "\n\n")
		}
		content.WriteString(fmt.Sprintf("Press %s to confirm, %s to cancel\n",
			lipgloss.NewStyle().Bold(true).Render("y"),
			lipgloss.NewStyle().Bold(true).Render("n")))

		return content.String()
	},
}

// openConfirmDelete opens the delete confirmation modal
func openConfirmDelete(m *Model, idx int) {
	m.CurrentMode = m.ConfirmMode
	m.ConfirmState = ConfirmState{DeletingIdx: idx}
	m.Status = ""
}
