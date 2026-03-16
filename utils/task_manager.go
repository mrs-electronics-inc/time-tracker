package utils

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"time-tracker/models"
)

// Storage interface for abstracting data persistence
type Storage interface {
	Load() ([]models.TimeEntry, error)
	Save([]models.TimeEntry) error
	LoadProjects() ([]models.Project, error)
	SaveProjects([]models.Project) error
}

type TaskManager struct {
	storage Storage
}

type ProjectMutationResult struct {
	RewrittenEntries int
	Merged           bool
}

type ProjectInUseError struct {
	ProjectName    string
	ReferenceCount int
}

func (e *ProjectInUseError) Error() string {
	return fmt.Sprintf("cannot remove project %q: referenced by %d time entries", e.ProjectName, e.ReferenceCount)
}

func NewTaskManager(storage Storage) *TaskManager {
	return &TaskManager{storage: storage}
}

func (tm *TaskManager) StartEntry(project, title string) (*models.TimeEntry, error) {
	return tm.StartEntryAt(project, title, time.Now())
}

func (tm *TaskManager) StartEntryAt(project, title string, startTime time.Time) (*models.TimeEntry, error) {
	entries, err := tm.storage.Load()
	if err != nil {
		return nil, err
	}

	// Check if we're already tracking this exact project/title
	if len(entries) > 0 {
		lastEntry := &entries[len(entries)-1]
		if lastEntry.IsRunning() && lastEntry.Project == project && lastEntry.Title == title {
			return nil, fmt.Errorf("already tracking: %s", project)
		}
	}

	// Stop last entry
	if len(entries) > 0 {
		lastEntry := &entries[len(entries)-1]
		if lastEntry.IsRunning() {
			entries[len(entries)-1].End = &startTime
		}
	}

	// Create new entry
	newEntry := models.TimeEntry{
		Start:   startTime,
		End:     nil,
		Project: project,
		Title:   title,
	}

	entries = append(entries, newEntry)

	if err := tm.storage.Save(entries); err != nil {
		return nil, err
	}

	return &newEntry, nil
}

func (tm *TaskManager) StopEntry() (*models.TimeEntry, error) {
	entries, err := tm.storage.Load()
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("no active time entry to stop")
	}

	lastEntry := &entries[len(entries)-1]
	if !lastEntry.IsRunning() || lastEntry.IsBlank() {
		return nil, fmt.Errorf("no active time entry to stop")
	}

	now := time.Now()
	entries[len(entries)-1].End = &now

	// Add a blank entry to represent the gap after the stopped entry
	blankEntry := models.TimeEntry{
		Start:   now,
		End:     nil,
		Project: "",
		Title:   "",
	}
	entries = append(entries, blankEntry)

	if err := tm.storage.Save(entries); err != nil {
		return nil, err
	}

	return lastEntry, nil
}

func (tm *TaskManager) ListEntries() ([]models.TimeEntry, error) {
	entries, err := tm.storage.Load()
	if err != nil {
		return nil, err
	}

	// Sort by start time ascending (oldest first)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Start.Before(entries[j].Start)
	})

	return entries, nil
}

// UpdateEntry updates an existing entry's project, title, and start time
func (tm *TaskManager) UpdateEntry(idx int, project, title string, startTime time.Time) error {
	entries, err := tm.storage.Load()
	if err != nil {
		return err
	}

	if idx < 0 || idx >= len(entries) {
		return fmt.Errorf("invalid entry index: %d", idx)
	}

	entries[idx].Project = project
	entries[idx].Title = title
	entries[idx].Start = startTime

	return tm.storage.Save(entries)
}

// DeleteEntry removes blank entries or converts non-blank entries to blank
func (tm *TaskManager) DeleteEntry(idx int) error {
	entries, err := tm.storage.Load()
	if err != nil {
		return err
	}

	if idx < 0 || idx >= len(entries) {
		return fmt.Errorf("invalid entry index: %d", idx)
	}

	if entries[idx].IsBlank() {
		// Remove blank entries from the slice
		entries = append(entries[:idx], entries[idx+1:]...)
	} else {
		// Convert non-blank entries to blank
		entries[idx].Project = ""
		entries[idx].Title = ""
	}

	return tm.storage.Save(entries)
}

func (tm *TaskManager) AddProject(name, code, category string) (*models.Project, error) {
	name = strings.TrimSpace(name)
	code = strings.TrimSpace(code)
	category = strings.TrimSpace(category)

	if name == "" {
		return nil, fmt.Errorf("project name cannot be empty")
	}

	projects, err := tm.storage.LoadProjects()
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		if strings.EqualFold(project.Name, name) {
			return nil, fmt.Errorf("project %q already exists", name)
		}
	}

	newProject := models.Project{Name: name, Code: code, Category: category}
	projects = append(projects, newProject)

	if err := tm.storage.SaveProjects(projects); err != nil {
		return nil, err
	}

	return &newProject, nil
}

func (tm *TaskManager) EditProject(name, newName, code, category string) (*ProjectMutationResult, error) {
	name = strings.TrimSpace(name)
	newName = strings.TrimSpace(newName)
	code = strings.TrimSpace(code)
	category = strings.TrimSpace(category)

	if name == "" {
		return nil, fmt.Errorf("project name cannot be empty")
	}

	projects, err := tm.storage.LoadProjects()
	if err != nil {
		return nil, err
	}

	sourceIndex := -1
	for i, project := range projects {
		if strings.EqualFold(project.Name, name) {
			sourceIndex = i
			break
		}
	}
	if sourceIndex < 0 {
		return nil, fmt.Errorf("project %q not found", name)
	}

	source := projects[sourceIndex]
	if newName == "" {
		newName = source.Name
	}

	if newName == "" {
		return nil, fmt.Errorf("project name cannot be empty")
	}

	targetIndex := -1
	for i, project := range projects {
		if i != sourceIndex && strings.EqualFold(project.Name, newName) {
			targetIndex = i
			break
		}
	}

	entries, err := tm.storage.Load()
	if err != nil {
		return nil, err
	}

	rewrittenEntries := 0
	if targetIndex >= 0 {
		targetName := projects[targetIndex].Name
		for i := range entries {
			if entries[i].Project == source.Name {
				entries[i].Project = targetName
				rewrittenEntries++
			}
		}

		projects = append(projects[:sourceIndex], projects[sourceIndex+1:]...)
		if err := tm.storage.Save(entries); err != nil {
			return nil, err
		}
		if err := tm.storage.SaveProjects(projects); err != nil {
			return nil, err
		}

		return &ProjectMutationResult{RewrittenEntries: rewrittenEntries, Merged: true}, nil
	}

	renamed := source.Name != newName
	if renamed {
		for i := range entries {
			if entries[i].Project == source.Name {
				entries[i].Project = newName
				rewrittenEntries++
			}
		}
	}

	source.Name = newName
	source.Code = code
	source.Category = category
	projects[sourceIndex] = source

	if renamed {
		if err := tm.storage.Save(entries); err != nil {
			return nil, err
		}
		if err := tm.storage.SaveProjects(projects); err != nil {
			return nil, err
		}
	} else {
		if err := tm.storage.SaveProjects(projects); err != nil {
			return nil, err
		}
	}

	return &ProjectMutationResult{RewrittenEntries: rewrittenEntries, Merged: false}, nil
}

func (tm *TaskManager) RemoveProject(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	projects, err := tm.storage.LoadProjects()
	if err != nil {
		return err
	}

	projectIndex := -1
	projectName := ""
	for i, project := range projects {
		if strings.EqualFold(project.Name, name) {
			projectIndex = i
			projectName = project.Name
			break
		}
	}
	if projectIndex < 0 {
		return fmt.Errorf("project %q not found", name)
	}

	entries, err := tm.storage.Load()
	if err != nil {
		return err
	}

	referenceCount := 0
	for _, entry := range entries {
		if entry.Project == projectName {
			referenceCount++
		}
	}

	if referenceCount > 0 {
		return &ProjectInUseError{ProjectName: projectName, ReferenceCount: referenceCount}
	}

	projects = append(projects[:projectIndex], projects[projectIndex+1:]...)
	return tm.storage.SaveProjects(projects)
}
