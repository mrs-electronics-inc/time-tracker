package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"time-tracker/cmd/tui/modes"
	"time-tracker/models"
	"time-tracker/utils"
)

// OperationCompleteMsg is sent when an async operation completes
type OperationCompleteMsg struct {
	Error error
}

// Model wraps the modes.Model and adds TUI-specific state
type Model struct {
	*modes.Model
}

// NewModel creates a new TUI model
func NewModel(storage models.Storage, taskManager *utils.TaskManager) *Model {
	// Create textinput models for start mode
	inputs := make([]textinput.Model, 4)

	// Project input
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Project"
	inputs[0].CharLimit = 128
	inputs[0].Width = 40
	inputs[0].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	inputs[0].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Title input
	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Title"
	inputs[1].CharLimit = 128
	inputs[1].Width = 40
	inputs[1].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	inputs[1].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	// Hour input
	inputs[2] = textinput.New()
	inputs[2].Placeholder = "HH"
	inputs[2].CharLimit = 2
	inputs[2].Width = 2
	inputs[2].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	inputs[2].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	// Minute input
	inputs[3] = textinput.New()
	inputs[3].Placeholder = "MM"
	inputs[3].CharLimit = 2
	inputs[3].Width = 2
	inputs[3].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	inputs[3].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	modesModel := &modes.Model{
		Storage:     storage,
		TaskManager: taskManager,
		Entries:     []models.TimeEntry{},
		SelectedIdx: 0,
		Inputs:      inputs,
		FocusIndex:  0,
		Loading:     false,
		Styles: modes.Styles{
			Header:        lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")),
			Footer:        lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
			Selected:      lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true),
			Unselected:    lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
			Running:       lipgloss.NewStyle().Foreground(lipgloss.Color("11")),
			Gap:           lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true),
			InputFocused:  lipgloss.NewStyle().Foreground(lipgloss.Color("205")),
			InputBlurred:  lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
			Title:         lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")),
			Label:         lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
			StatusError:   lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true),
			StatusSuccess: lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true),
		},
	}

	// Initialize modes
	modesModel.ListMode = modes.ListMode
	modesModel.StartMode = modes.StartMode
	modesModel.HelpMode = modes.HelpMode
	modesModel.CurrentMode = modesModel.ListMode

	return &Model{Model: modesModel}
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case OperationCompleteMsg:
		m.Loading = false
		if msg.Error != nil {
			m.Status = "Error: " + msg.Error.Error()
		}
		return m, nil

	case tea.KeyMsg:
		// Delegate to current mode's key handler
		if m.CurrentMode.HandleKeyMsg != nil {
			model, cmd := m.CurrentMode.HandleKeyMsg(m.Model, msg)
			return &Model{Model: model}, cmd
		}
	}

	return m, nil
}

// View renders the UI
func (m *Model) View() string {
	if m.Err != nil {
		return "Error: " + m.Err.Error() + "\n"
	}

	// Calculate available height (status bar is always 1 line)
	statusBarHeight := 1
	availableHeight := max(m.Height-statusBarHeight, 1)

	// Render mode-specific content
	content := m.CurrentMode.RenderContent(m.Model, availableHeight)

	// Add status bar
	statusBar := m.renderStatusBar()
	contentLines := countNewlines(content)
	spacerLines := m.Height - contentLines - statusBarHeight
	spacerLines = max(spacerLines, 0)
	var result strings.Builder
	result.WriteString(content)
	if spacerLines > 0 {
		result.WriteString(strings.Repeat("\n", spacerLines))
	}
	result.WriteString(statusBar)
	return result.String()
}

// renderStatusBar renders a zellij-style status bar with mode and keybindings
func (m *Model) renderStatusBar() string {
	// Colors
	black := lipgloss.Color("0")
	magenta := lipgloss.Color("5")
	gray := lipgloss.Color("8")
	green := lipgloss.Color("10")

	// Styles
	modeStyle := lipgloss.NewStyle().
		Background(green).
		Foreground(black).
		Bold(true).
		Padding(0, 1)

	keyStyle := lipgloss.NewStyle().
		Background(black).
		Foreground(magenta).
		Padding(0, 1)

	labelStyle := lipgloss.NewStyle().
		Background(gray).
		Foreground(black).
		Bold(true).
		Padding(0, 1)

	// Separators
	powerlineSeparator := "\uE0B0"

	modeSep := lipgloss.NewStyle().
		Background(black).
		Foreground(green).
		Render(powerlineSeparator)

	keySep := lipgloss.NewStyle().
		Background(gray).
		Foreground(black).
		Render(powerlineSeparator)

	labelSep := lipgloss.NewStyle().
		Background(black).
		Foreground(gray).
		Render(powerlineSeparator)

	// Helper to render a key-label pair with powerline separators
	renderPair := func(key, label string) string {
		return keyStyle.Render(key) + keySep + labelStyle.Render(label) + labelSep
	}

	var parts []string

	// Mode indicator
	parts = append(parts, modeStyle.Render(strings.ToUpper(m.CurrentMode.Name))+modeSep)

	// Add keybindings from current mode
	for _, binding := range m.CurrentMode.StatusBarKeys {
		parts = append(parts, renderPair(binding.Key, binding.Label))
	}

	// Build left side of status bar
	leftSide := strings.Join(parts, "")

	// Add status message on the right side if present
	if m.Status != "" {
		statusStyle := lipgloss.NewStyle().
			Foreground(magenta).
			Padding(0, 1)
		rightSide := statusStyle.Render(m.Status)

		// Calculate padding to right-align status
		leftWidth := lipgloss.Width(leftSide)
		rightWidth := lipgloss.Width(rightSide)
		totalWidth := leftWidth + rightWidth
		paddingWidth := max(m.Width-totalWidth, 0)

		padding := strings.Repeat(" ", paddingWidth)
		return leftSide + padding + rightSide
	}

	return leftSide
}

// Helper functions
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func countNewlines(s string) int {
	return strings.Count(s, "\n")
}
