package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"time-tracker/models"
	"time-tracker/utils"
)

// keyMap defines keybindings for the TUI
type keyMap struct {
	Help       key.Binding
	Toggle     key.Binding
	Quit       key.Binding
	Up         key.Binding
	Down       key.Binding
	JumpBottom key.Binding
}

// ShortHelp returns keybindings shown in the mini help view
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.JumpBottom, k.Toggle, k.Help, k.Quit}
}

// FullHelp returns all keybindings
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.JumpBottom},
		{k.Toggle, k.Help, k.Quit},
	}
}

// keys defines the default keybindings
var keys = keyMap{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Toggle: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "toggle start/stop"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k/↑", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j/↓", "down"),
	),
	JumpBottom: key.NewBinding(
		key.WithKeys("G"),
		key.WithHelp("G", "jump to bottom"),
	),
}

// Model represents the state of the TUI application
type Model struct {
	storage     models.Storage     // Persistent storage backend
	taskManager *utils.TaskManager // Task management operations
	entries     []models.TimeEntry // Loaded time entries
	selectedIdx int                // Index of currently selected entry
	viewportTop int                // Index of first visible row
	err         error              // Error state
	status      string             // Status message from last action
	keys        keyMap             // Keybindings
	styles      styles             // UI styling
	width       int                // Terminal width
	height      int                // Terminal height
	showHelp    bool               // Whether to show full help text
	help        help.Model         // Help component

	// Dialog state
	dialogMode  bool                  // Whether we're in dialog mode
	inputs      []textinput.Model     // Text inputs for project and title
	focusIndex  int                   // Currently focused input (0 = project, 1 = title)
}

// styles defines the visual styling for different UI elements
type styles struct {
	header          lipgloss.Style // Header row style
	footer          lipgloss.Style // Footer style
	selected        lipgloss.Style // Selected row style
	unselected      lipgloss.Style // Unselected row style
	running         lipgloss.Style // Running entry style
	gap             lipgloss.Style // Gap/blank entry style
	dialogFocused   lipgloss.Style // Dialog focused input style
	dialogBlurred   lipgloss.Style // Dialog blurred input style
}

// NewModel creates a new TUI model
func NewModel(storage models.Storage, taskManager *utils.TaskManager) *Model {
	h := help.New()
	h.ShowAll = false

	// Create textinput models for dialog
	inputs := make([]textinput.Model, 2)

	// Project input
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Project"
	inputs[0].CharLimit = 128
	inputs[0].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	inputs[0].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Title input
	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Title"
	inputs[1].CharLimit = 128
	inputs[1].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	inputs[1].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	return &Model{
		storage:     storage,
		taskManager: taskManager,
		entries:     []models.TimeEntry{},
		selectedIdx: 0,
		keys:        keys,
		showHelp:    false,
		help:        h,
		dialogMode:  false,
		inputs:      inputs,
		focusIndex:  0,
		styles: styles{
			header:        lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")),
			footer:        lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
			selected:      lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true),
			unselected:    lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
			running:       lipgloss.NewStyle().Foreground(lipgloss.Color("11")),
			gap:           lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true),
			dialogFocused: lipgloss.NewStyle().Foreground(lipgloss.Color("205")),
			dialogBlurred: lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		},
	}
}

// LoadEntries loads time entries from storage
func (m *Model) LoadEntries() error {
	entries, err := m.storage.Load()
	if err != nil {
		return err
	}

	m.entries = entries

	// Select most recent entry (last item)
	if len(m.entries) > 0 {
		m.selectedIdx = len(m.entries) - 1
	} else {
		m.selectedIdx = 0
	}

	return nil
}

// getColumnWidths calculates the width of each column based on content
func (m *Model) getColumnWidths() (int, int, int, int, int) {
	// Minimum widths for headers
	startWidth := len("Start")
	endWidth := len("End")
	projectWidth := len("Project")
	titleWidth := len("Title")
	durationWidth := len("Duration")

	// Measure content
	for _, entry := range m.entries {
		startStr := entry.Start.Format("2006-01-02 15:04")
		if len(startStr) > startWidth {
			startWidth = len(startStr)
		}

		endStr := "running"
		if entry.End != nil {
			endStr = entry.End.Format("2006-01-02 15:04")
		}
		if len(endStr) > endWidth {
			endWidth = len(endStr)
		}

		if len(entry.Project) > projectWidth {
			projectWidth = len(entry.Project)
		}

		if len(entry.Title) > titleWidth {
			titleWidth = len(entry.Title)
		}

		duration := formatDuration(entry.Duration())
		if len(duration) > durationWidth {
			durationWidth = len(duration)
		}
	}

	return startWidth, endWidth, projectWidth, titleWidth, durationWidth
}

// formatDuration converts a time.Duration to a human-readable string (e.g., "2h 15m")
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

// ensureSelectionVisible adjusts viewport so selected item is visible
func (m *Model) ensureSelectionVisible(maxVisibleRows int) {
	if m.selectedIdx < m.viewportTop {
		// Selected item is above viewport, scroll up
		m.viewportTop = m.selectedIdx
	} else if m.selectedIdx >= m.viewportTop+maxVisibleRows {
		// Selected item is below viewport, scroll down
		m.viewportTop = m.selectedIdx - maxVisibleRows + 1
	}

	// Ensure viewport doesn't go past the end
	if m.viewportTop > len(m.entries)-maxVisibleRows {
		m.viewportTop = len(m.entries) - maxVisibleRows
	}
	if m.viewportTop < 0 {
		m.viewportTop = 0
	}
}
