package tui

import (
	"sort"
	"strings"
	"time-tracker/models"
)

// AutocompleteResult represents a single autocomplete suggestion
type AutocompleteResult struct {
	Project string // Project name
	Title   string // Task title
	Score   int    // Relevance score (higher = more recent)
}

// AutocompleteSuggestions holds project and task suggestions
type AutocompleteSuggestions struct {
	Projects           []string           // Unique projects
	Tasks              []AutocompleteResult // (project, title) combinations
	FilteredResults    []AutocompleteResult // Currently filtered task results
	FilteredProjects   []string            // Currently filtered project results
	selectedIdx        int                // Index of currently selected suggestion
	isProjectFiltering bool               // Whether we're currently filtering projects
}

// NewAutocompleteSuggestions creates a new autocomplete suggestions object
func NewAutocompleteSuggestions() *AutocompleteSuggestions {
	return &AutocompleteSuggestions{
		Projects:           []string{},
		Tasks:              []AutocompleteResult{},
		FilteredResults:    []AutocompleteResult{},
		FilteredProjects:   []string{},
		selectedIdx:        0,
		isProjectFiltering: false,
	}
}

// ExtractFromEntries extracts unique projects and (project, title) combinations
// from time entries, weighted toward most recent ones
func (a *AutocompleteSuggestions) ExtractFromEntries(entries []models.TimeEntry) {
	// Reset slices to avoid accumulation from previous calls
	a.Projects = []string{}
	a.Tasks = []AutocompleteResult{}
	a.FilteredResults = []AutocompleteResult{}
	a.FilteredProjects = []string{}
	a.selectedIdx = 0

	projectMap := make(map[string]bool)
	taskMap := make(map[string]AutocompleteResult) // key: "project|title"

	// Iterate in reverse (most recent first) to weight scores properly
	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]

		// Skip blank entries
		if entry.IsBlank() {
			continue
		}

		// Track unique projects
		if !projectMap[entry.Project] {
			projectMap[entry.Project] = true
			a.Projects = append(a.Projects, entry.Project)
		}

		// Track unique (project, title) combinations with scores
		key := entry.Project + "|" + entry.Title
		if _, exists := taskMap[key]; !exists {
			// Score is based on reverse position (most recent = higher score)
			score := len(entries) - i
			taskMap[key] = AutocompleteResult{
				Project: entry.Project,
				Title:   entry.Title,
				Score:   score,
			}
		}
	}

	// Convert task map to slice and sort by score (descending)
	for _, task := range taskMap {
		a.Tasks = append(a.Tasks, task)
	}

	// Sort by score descending (most recent first)
	sort.Slice(a.Tasks, func(i, j int) bool {
		return a.Tasks[i].Score > a.Tasks[j].Score
	})

	// Sort projects alphabetically
	sort.Strings(a.Projects)

	// Initialize filtered results with all tasks
	a.FilteredResults = a.Tasks
	a.FilteredProjects = a.Projects
	a.selectedIdx = 0
}

// FilterProjects filters project suggestions based on input text
func (a *AutocompleteSuggestions) FilterProjects(input string) {
	if input == "" {
		a.FilteredProjects = a.Projects
		a.isProjectFiltering = true
		a.selectedIdx = 0
		return
	}

	input = strings.ToLower(input)
	var results []string
	for _, project := range a.Projects {
		if strings.Contains(strings.ToLower(project), input) {
			results = append(results, project)
		}
	}
	a.FilteredProjects = results
	a.isProjectFiltering = true
	a.selectedIdx = 0
}

// FilterTasks filters task suggestions based on input text
// If projectFilter is not empty, only tasks from that project are returned
func (a *AutocompleteSuggestions) FilterTasks(input string, projectFilter string) {
	input = strings.ToLower(input)
	projectFilter = strings.ToLower(projectFilter)

	var results []AutocompleteResult
	for _, task := range a.Tasks {
		// If project filter is set, match on project first
		if projectFilter != "" && strings.ToLower(task.Project) != projectFilter {
			continue
		}

		// Match input against title only if input is provided
		if input != "" && !strings.Contains(strings.ToLower(task.Title), input) {
			continue
		}

		results = append(results, task)
	}

	a.FilteredResults = results
	a.isProjectFiltering = false
	a.selectedIdx = 0
}

// GetSelectedSuggestion returns the currently selected suggestion
func (a *AutocompleteSuggestions) GetSelectedSuggestion() *AutocompleteResult {
	if a.selectedIdx >= 0 && a.selectedIdx < len(a.FilteredResults) {
		return &a.FilteredResults[a.selectedIdx]
	}
	return nil
}

// SelectNext moves selection to next suggestion
func (a *AutocompleteSuggestions) SelectNext() {
	if len(a.FilteredResults) > 0 {
		a.selectedIdx = (a.selectedIdx + 1) % len(a.FilteredResults)
	}
}

// SelectPrev moves selection to previous suggestion
func (a *AutocompleteSuggestions) SelectPrev() {
	if len(a.FilteredResults) > 0 {
		a.selectedIdx--
		if a.selectedIdx < 0 {
			a.selectedIdx = len(a.FilteredResults) - 1
		}
	}
}
