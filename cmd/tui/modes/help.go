package modes

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HelpMode is the help/keybinding reference mode
var HelpMode = &Mode{
	Name:          "help",
	HandleKeyMsg:  handleHelpKeyMsg,
	RenderContent: renderHelpContent,
	StatusBarKeys: []StatusBarKeyBinding{
		{Key: "Esc", Label: "BACK"},
	},
	Help: []HelpEntry{
		{Keys: "Esc", Desc: "Back to previous mode"},
	},
}

// handleHelpKeyMsg handles key messages while in help mode
func handleHelpKeyMsg(m *Model, msg tea.KeyMsg) (*Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "?", "q":
		m.CurrentMode = m.ListMode
		return m, nil
	}
	return m, nil
}

// renderHelpContent renders the help mode content
func renderHelpContent(m *Model, availableHeight int) string {
	_ = availableHeight // Available for future use
	var content strings.Builder

	title := m.Styles.Title.Render("Keyboard Shortcuts")
	content.WriteString(title + "\n\n")

	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))

	// Render help entries from current mode
	for _, entry := range m.CurrentMode.Help {
		content.WriteString(keyStyle.Render(entry.Keys) + " " + descStyle.Render(entry.Desc) + "\n")
	}

	return content.String()
}
