package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case OperationCompleteMsg:
		m.loading = false
		if msg.Error != nil {
			m.status = "Error: " + msg.Error.Error()
		}
		return m, nil

	case tea.KeyMsg:
		// Start mode key handling
		if m.mode == ModeStart {
			return m.handleStartKeyMsg(msg)
		}

		// Help mode key handling
		if m.mode == ModeHelp {
			return m.handleHelpKeyMsg(msg)
		}

		// List mode key handling
		if key.Matches(msg, m.keys.Help) {
			m.prevMode = m.mode
			m.mode = ModeHelp
			return m, nil
		}

		if key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		}

		// Navigation
		if key.Matches(msg, m.keys.Up) {
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
			m.status = ""
			return m, nil
		}
		if key.Matches(msg, m.keys.Down) {
			if m.selectedIdx < len(m.entries)-1 {
				m.selectedIdx++
			}
			m.status = ""
			return m, nil
		}

		// Jump to bottom
		if key.Matches(msg, m.keys.JumpBottom) {
			if len(m.entries) > 0 {
				m.selectedIdx = len(m.entries) - 1
			}
			m.status = ""
			return m, nil
		}

		// Toggle start/stop
		if key.Matches(msg, m.keys.Toggle) {
			if len(m.entries) == 0 {
				// No entries yet - open blank start mode
				m.openStartModeBlank()
			} else if m.selectedIdx >= 0 && m.selectedIdx < len(m.entries) {
				entry := m.entries[m.selectedIdx]
				if entry.IsRunning() {
					// Stop entry
					if _, err := m.taskManager.StopEntry(); err != nil {
						m.status = "Error stopping entry: " + err.Error()
					} else {
						m.status = "Entry stopped"
					}
				} else if !entry.IsBlank() {
					// Start new entry based on selected
					m.openStartMode(entry)
				} else {
					// Blank entry - open blank start mode
					m.openStartModeBlank()
				}
				// Reload entries to update display
				if err := m.LoadEntries(); err != nil {
					m.err = err
				}
			}
			return m, nil
		}
	}

	return m, nil
}

// View renders the UI
func (m *Model) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}

	// Calculate available height (status bar is always 1 line)
	statusBarHeight := 1
	availableHeight := max(m.height-statusBarHeight, 1)

	// Render mode-specific content
	content := m.renderModeContent(availableHeight)

	// Add status bar
	statusBar := m.renderStatusBar()
	contentLines := strings.Count(content, "\n")
	spacerLines := m.height - contentLines - statusBarHeight
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

	// Mode indicator and keybindings based on current mode
	switch m.mode {
	case ModeStart:
		parts = append(parts, modeStyle.Render("START")+modeSep)
		parts = append(parts, renderPair("Tab", "NEXT"))
		parts = append(parts, renderPair("Enter", "SUBMIT"))
		parts = append(parts, renderPair("Esc", "CANCEL"))
	case ModeHelp:
		parts = append(parts, modeStyle.Render("HELP")+modeSep)
		parts = append(parts, renderPair("Esc", "BACK"))
	default: // ModeList
		parts = append(parts, modeStyle.Render("LIST")+modeSep)
		parts = append(parts, renderPair("j/k", "NAVIGATE"))
		parts = append(parts, renderPair("G", "GO TO CURRENT"))
		parts = append(parts, renderPair("s", "START/STOP"))
		parts = append(parts, renderPair("?", "HELP"))
		parts = append(parts, renderPair("Esc", "QUIT"))
	}

	// Build left side of status bar
	leftSide := strings.Join(parts, "")

	// Add status message on the right side if present
	if m.status != "" {
		statusStyle := lipgloss.NewStyle().
			Foreground(magenta).
			Padding(0, 1)
		rightSide := statusStyle.Render(m.status)

		// Calculate padding to right-align status
		leftWidth := lipgloss.Width(leftSide)
		rightWidth := lipgloss.Width(rightSide)
		totalWidth := leftWidth + rightWidth
		paddingWidth := max(m.width-totalWidth, 0)

		padding := strings.Repeat(" ", paddingWidth)
		return leftSide + padding + rightSide
	}

	return leftSide
}
