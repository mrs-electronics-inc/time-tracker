package modes

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time-tracker/models"
	"time-tracker/utils"
)

// KeyBinding represents a single key binding with its display and description
type KeyBinding struct {
	Keys        string // Key sequence (e.g., "q", "ctrl+c")
	Label       string // Shown in the status bar
	Description string // Shown in the help page
}

// Mode represents a TUI mode with its keybindings and renderer
type Mode struct {
	Name          string
	KeyBindings   []KeyBinding
	HandleKeyMsg  func(*Model, tea.KeyMsg) (*Model, tea.Cmd)
	RenderContent func(*Model, int) string
}

// Styles defines the visual styling for different UI elements
type Styles struct {
	Header        lipgloss.Style // Header row style
	Footer        lipgloss.Style // Footer style
	Selected      lipgloss.Style // Selected row style
	Unselected    lipgloss.Style // Unselected row style
	Running       lipgloss.Style // Running entry style
	Gap           lipgloss.Style // Gap/blank entry style
	InputFocused  lipgloss.Style // Focused input style
	InputBlurred  lipgloss.Style // Blurred input style
	Title         lipgloss.Style // Screen title style
	Label         lipgloss.Style // Input label style
	StatusError   lipgloss.Style // Status error message style
	StatusSuccess lipgloss.Style // Status success message style
}

// Model represents the state of the TUI application
type Model struct {
	Storage     models.Storage     // Persistent storage backend
	TaskManager *utils.TaskManager // Task management operations
	Entries     []models.TimeEntry // Loaded time entries
	SelectedIdx int                // Index of currently selected entry (list mode)
	ViewportTop int                // Index of first visible row (list mode) or viewport scroll position (stats mode)
	Err         error              // Error state
	Status      string             // Status message from last action
	Styles      Styles             // UI styling
	Width       int                // Terminal width
	Height      int                // Terminal height

	// Mode state
	CurrentMode  *Mode             // Current TUI mode
	PreviousMode *Mode             // Previous mode (used for help context)
	Inputs       []textinput.Model // Text inputs for project, title, hour, minute
	FocusIndex   int               // Currently focused input (0 = project, 1 = title, 2 = hour, 3 = minute)

	// Loading state
	Loading bool // Whether we're waiting for a data operation

	// Mode references for navigation
	ListMode    *Mode
	StartMode   *Mode
	HelpMode    *Mode
	StatsMode   *Mode
	NewMode     *Mode
	EditMode    *Mode
	ResumeMode  *Mode
	ConfirmMode *Mode

	// Form state for new/edit/resume modes
	FormState FormState

	// Confirm state for delete confirmation
	ConfirmState ConfirmState
}

// LoadEntries loads time entries from storage
func (m *Model) LoadEntries() error {
	entries, err := m.Storage.Load()
	if err != nil {
		return err
	}

	m.Entries = entries

	// Select most recent entry (last item)
	if len(m.Entries) > 0 {
		m.SelectedIdx = len(m.Entries) - 1
	} else {
		m.SelectedIdx = 0
	}

	return nil
}

// GetColumnWidths calculates the width of each column based on content
func (m *Model) GetColumnWidths() (int, int, int, int, int) {
	// Minimum widths for headers
	startWidth := len("Start")
	endWidth := len("End")
	projectWidth := len("Project")
	titleWidth := len("Title")
	durationWidth := len("Duration")

	// Measure content
	for _, entry := range m.Entries {
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

		duration := utils.FormatDuration(entry.Duration())
		if len(duration) > durationWidth {
			durationWidth = len(duration)
		}
	}

	return startWidth, endWidth, projectWidth, titleWidth, durationWidth
}

// EnsureSelectionVisible adjusts viewport so selected item is visible
func (m *Model) EnsureSelectionVisible(maxVisibleRows int) {
	if m.SelectedIdx < m.ViewportTop {
		// Selected item is above viewport, scroll up
		m.ViewportTop = m.SelectedIdx
	} else if m.SelectedIdx >= m.ViewportTop+maxVisibleRows {
		// Selected item is below viewport, scroll down
		m.ViewportTop = m.SelectedIdx - maxVisibleRows + 1
	}

	// Ensure viewport doesn't go past the end
	if m.ViewportTop > len(m.Entries)-maxVisibleRows {
		m.ViewportTop = len(m.Entries) - maxVisibleRows
	}
	if m.ViewportTop < 0 {
		m.ViewportTop = 0
	}
}

// ResetScroll resets scroll position and viewport. Use -1 for ViewportTop to signal
// that stats mode should initialize to the bottom on first render.
func (m *Model) ResetScroll() {
	m.SelectedIdx = 0
	m.ViewportTop = -1
}

// SelectMostRecentEntry sets selection to the most recent entry
func (m *Model) SelectMostRecentEntry() {
	if len(m.Entries) > 0 {
		m.SelectedIdx = len(m.Entries) - 1
	} else {
		m.SelectedIdx = 0
	}
}

// SwitchMode changes mode, resets scroll, and clears status
func (m *Model) SwitchMode(newMode *Mode) {
	m.PreviousMode = m.CurrentMode
	m.CurrentMode = newMode
	m.ResetScroll()
	m.Status = ""

	// When switching to list mode, select most recent entry
	if newMode.Name == "list" {
		m.SelectMostRecentEntry()
	}
}
