package modes

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HelpMode is the help/keybinding reference mode
var HelpMode = &Mode{
	Name: "help",
	KeyBindings: []KeyBinding{
		{Keys: "Esc", Label: "BACK", Description: "Back to previous mode"},
	},
	HandleKeyMsg: func(m *Model, msg tea.KeyMsg) (*Model, tea.Cmd) {
		switch msg.String() {
		case "esc":
			m.CurrentMode = m.ListMode
			return m, nil
		}
		return m, nil
	},
	RenderContent: func(m *Model, availableHeight int) string {
		_ = availableHeight
		var content strings.Builder

		title := m.Styles.Title.Render("Keyboard Shortcuts")
		content.WriteString(title + "\n\n")

		keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
		descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))

		// Show all mode keybindings
		for _, binding := range m.CurrentMode.KeyBindings {
			content.WriteString(keyStyle.Render(binding.Keys) + " " + descStyle.Render(binding.Description) + "\n")
		}

		return content.String()
	},
}
