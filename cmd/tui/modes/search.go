package modes

import (
	"strings"

	"time-tracker/models"
)

type VisibleEntry struct {
	Entry       models.TimeEntry
	SourceIndex int
}

func filterVisibleEntries(entries []models.TimeEntry, query string) []VisibleEntry {
	visible := make([]VisibleEntry, 0, len(entries))

	for sourceIndex, entry := range entries {
		if entryMatchesSearchQuery(entry, query) {
			visible = append(visible, VisibleEntry{
				Entry:       entry,
				SourceIndex: sourceIndex,
			})
		}
	}

	return visible
}

func entryMatchesSearchQuery(entry models.TimeEntry, query string) bool {
	if query == "" {
		return true
	}

	normalizedQuery := strings.ToLower(query)
	project := strings.ToLower(entry.Project)
	title := strings.ToLower(entry.Title)

	return strings.Contains(project, normalizedQuery) || strings.Contains(title, normalizedQuery)
}

func applySearch(m *Model) {
	previousSelection := m.SelectedIdx
	m.SearchAppliedQuery = m.SearchQueryDraft
	m.FilteredEntries = filterVisibleEntries(m.Entries, m.SearchAppliedQuery)
	ensureValidFilteredSelection(m)

	if len(m.FilteredEntries) == 0 {
		return
	}

	for _, visible := range m.FilteredEntries {
		if visible.SourceIndex == previousSelection {
			return
		}
	}

	m.SelectedIdx = m.FilteredEntries[len(m.FilteredEntries)-1].SourceIndex
}

func clearSearch(m *Model) {
	m.SearchQueryDraft = ""
	m.SearchAppliedQuery = ""
	m.SearchActive = false
	m.SearchInputFocused = false
	m.FilteredEntries = filterVisibleEntries(m.Entries, "")
}

func ensureValidFilteredSelection(m *Model) {
	if len(m.Entries) == 0 {
		m.SelectedIdx = 0
		return
	}

	if m.SearchAppliedQuery == "" {
		if m.SelectedIdx < 0 || m.SelectedIdx >= len(m.Entries) {
			m.SelectedIdx = len(m.Entries) - 1
		}
		return
	}

	if len(m.FilteredEntries) == 0 {
		if m.SelectedIdx < 0 || m.SelectedIdx >= len(m.Entries) {
			m.SelectedIdx = len(m.Entries) - 1
		}
		return
	}

	for _, visible := range m.FilteredEntries {
		if visible.SourceIndex == m.SelectedIdx {
			return
		}
	}

	m.SelectedIdx = m.FilteredEntries[len(m.FilteredEntries)-1].SourceIndex
}
