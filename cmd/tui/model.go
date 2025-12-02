package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"time-tracker/models"
	"time-tracker/utils"
)

// keyMap defines keybindings for the TUI
type keyMap struct {
	Quit key.Binding
}

// ShortHelp returns keybindings shown in the mini help view
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

// FullHelp returns keybindings for the expanded help view
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit},
	}
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// Model represents the state of the TUI application
type Model struct {
	storage      models.Storage
	taskManager  *utils.TaskManager
	entries      []models.TimeEntry
	currentEntry *models.TimeEntry
	err          error
	width        int
	height       int
	keys         keyMap
	help         help.Model
	styles       styles
}

type styles struct {
	header   lipgloss.Style
	footer   lipgloss.Style
	divider  lipgloss.Style
	active   lipgloss.Style
	inactive lipgloss.Style
}

// NewModel creates a new TUI model
func NewModel(storage models.Storage, taskManager *utils.TaskManager) *Model {
	return &Model{
		storage:     storage,
		taskManager: taskManager,
		keys:        keys,
		help:        help.New(),
		styles: styles{
			header:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")),
			footer:   lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
			divider:  lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
			active:   lipgloss.NewStyle().Foreground(lipgloss.Color("10")).SetString("●"),
			inactive: lipgloss.NewStyle().Foreground(lipgloss.Color("8")).SetString("○"),
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

	// Find currently running entry
	for i := range entries {
		if entries[i].IsRunning() {
			m.currentEntry = &entries[i]
			break
		}
	}

	return nil
}
