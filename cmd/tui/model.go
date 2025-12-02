package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"time-tracker/models"
	"time-tracker/utils"
)

// keyMap defines keybindings for the TUI
type keyMap struct {
	Help   key.Binding
	Toggle key.Binding
	Quit   key.Binding
	Up     key.Binding
	Down   key.Binding
}

// ShortHelp returns keybindings shown in the mini help view
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Toggle, k.Quit}
}

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
}

// Model represents the state of the TUI application
type Model struct {
	storage      models.Storage
	taskManager  *utils.TaskManager
	entries      []models.TimeEntry
	selectedIdx  int
	err          error
	status       string         // Status message from last action
	keys         keyMap
	styles       styles
	width        int
	height       int
	showHelp     bool
}

type styles struct {
	header      lipgloss.Style
	footer      lipgloss.Style
	selected    lipgloss.Style
	unselected  lipgloss.Style
	running     lipgloss.Style
	gap         lipgloss.Style
}

// NewModel creates a new TUI model
func NewModel(storage models.Storage, taskManager *utils.TaskManager) *Model {
	return &Model{
		storage:     storage,
		taskManager: taskManager,
		entries:     []models.TimeEntry{},
		selectedIdx: 0,
		keys:        keys,
		showHelp:    false,
		styles: styles{
			header:     lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")),
			footer:     lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
			selected:   lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true),
			unselected: lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
			running:    lipgloss.NewStyle().Foreground(lipgloss.Color("11")),
			gap:        lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true),
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

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

// getRunningEntry returns the currently running entry, if any
func (m *Model) getRunningEntry() *models.TimeEntry {
	for i, entry := range m.entries {
		if entry.IsRunning() {
			return &m.entries[i]
		}
	}
	return nil
}
