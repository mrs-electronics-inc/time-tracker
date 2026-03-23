package modes

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"time-tracker/models"
)

type projectSuggestionStorage struct {
	projects []models.Project
}

func (s projectSuggestionStorage) Load() ([]models.TimeEntry, error)       { return nil, nil }
func (s projectSuggestionStorage) Save([]models.TimeEntry) error           { return nil }
func (s projectSuggestionStorage) LoadProjects() ([]models.Project, error) { return s.projects, nil }
func (s projectSuggestionStorage) SaveProjects([]models.Project) error     { return nil }

func newFormDateTestModel() *Model {
	inputs := make([]textinput.Model, InputMinute+1)
	for i := range inputs {
		inputs[i] = textinput.New()
	}
	inputs[InputProject].ShowSuggestions = true

	return &Model{
		Inputs:     inputs,
		NewMode:    NewMode,
		EditMode:   EditMode,
		ResumeMode: ResumeMode,
	}
}

func formDateString(m *Model) string {
	return m.Inputs[InputYear].Value() + "-" + m.Inputs[InputMonth].Value() + "-" + m.Inputs[InputDay].Value()
}

func assertCurrentDateDefault(t *testing.T, before time.Time, m *Model) {
	t.Helper()

	after := time.Now()
	actual := formDateString(m)
	beforeDate := before.Format("2006-01-02")
	afterDate := after.Format("2006-01-02")

	if actual != beforeDate && actual != afterDate {
		t.Fatalf("date default = %q, expected %q or %q", actual, beforeDate, afterDate)
	}
}

func TestOpenNewModeSetsCurrentDateDefaults(t *testing.T) {
	m := newFormDateTestModel()

	before := time.Now()
	openNewMode(m)

	assertCurrentDateDefault(t, before, m)
}

func TestOpenEditModeSetsDateDefaultsFromEntry(t *testing.T) {
	m := newFormDateTestModel()
	entry := models.TimeEntry{Start: time.Date(2025, time.January, 7, 9, 30, 0, 0, time.UTC)}

	openEditMode(m, entry, 0)

	if got := m.Inputs[InputYear].Value(); got != "2025" {
		t.Fatalf("year = %q, expected 2025", got)
	}
	if got := m.Inputs[InputMonth].Value(); got != "01" {
		t.Fatalf("month = %q, expected 01", got)
	}
	if got := m.Inputs[InputDay].Value(); got != "07" {
		t.Fatalf("day = %q, expected 07", got)
	}
}

func TestOpenResumeModeSetsCurrentDateDefaults(t *testing.T) {
	m := newFormDateTestModel()
	entry := models.TimeEntry{Project: "backend", Title: "review"}

	before := time.Now()
	openResumeMode(m, entry)

	assertCurrentDateDefault(t, before, m)
}

func TestOpenFormModesSetProjectSuggestions(t *testing.T) {
	projects := []models.Project{{Name: " beta "}, {Name: "Alpha"}, {Name: "alpha"}, {Name: ""}, {Name: "Gamma"}}
	m := newFormDateTestModel()
	m.Storage = projectSuggestionStorage{projects: projects}

	openNewMode(m)
	if got := m.Inputs[InputProject].AvailableSuggestions(); len(got) != 3 || got[0] != "Alpha" || got[1] != "beta" || got[2] != "Gamma" {
		t.Fatalf("new mode suggestions = %+v", got)
	}

	openEditMode(m, models.TimeEntry{Project: "beta", Title: "task", Start: time.Now()}, 0)
	if got := m.Inputs[InputProject].AvailableSuggestions(); len(got) != 3 || got[0] != "Alpha" || got[1] != "beta" || got[2] != "Gamma" {
		t.Fatalf("edit mode suggestions = %+v", got)
	}

	openResumeMode(m, models.TimeEntry{Project: "beta", Title: "task"})
	if got := m.Inputs[InputProject].AvailableSuggestions(); len(got) != 3 || got[0] != "Alpha" || got[1] != "beta" || got[2] != "Gamma" {
		t.Fatalf("resume mode suggestions = %+v", got)
	}
}

func TestFormModeTabAcceptsProjectSuggestionBeforeChangingFocus(t *testing.T) {
	m := newFormDateTestModel()
	m.Storage = projectSuggestionStorage{projects: []models.Project{{Name: "Backend"}}}

	openNewMode(m)

	var cmd tea.Cmd
	m.Inputs[InputProject], cmd = m.Inputs[InputProject].Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	if cmd != nil {
		cmd()
	}

	if got := m.Inputs[InputProject].MatchedSuggestions(); len(got) != 1 || got[0] != "Backend" {
		t.Fatalf("matched suggestions = %+v", got)
	}

	updatedModel, _ := NewMode.HandleKeyMsg(m, tea.KeyMsg{Type: tea.KeyTab})

	if updatedModel.Inputs[InputProject].Value() != "backend" {
		t.Fatalf("project value = %q, expected suggestion to be accepted", updatedModel.Inputs[InputProject].Value())
	}
	if updatedModel.FocusIndex != InputProject {
		t.Fatalf("focus index = %d, expected project input to stay focused", updatedModel.FocusIndex)
	}
}

func TestFormModeTabNavigatesWhenProjectAlreadyExactlyMatchesSuggestion(t *testing.T) {
	m := newFormDateTestModel()
	m.Storage = projectSuggestionStorage{projects: []models.Project{{Name: "Backend"}}}

	openNewMode(m)
	m.Inputs[InputProject].SetValue("Backend")

	updatedModel, _ := NewMode.HandleKeyMsg(m, tea.KeyMsg{Type: tea.KeyTab})

	if updatedModel.Inputs[InputProject].Value() != "Backend" {
		t.Fatalf("project value = %q, expected project value to remain unchanged", updatedModel.Inputs[InputProject].Value())
	}
	if updatedModel.FocusIndex != InputTitle {
		t.Fatalf("focus index = %d, expected focus to move to title input", updatedModel.FocusIndex)
	}
}
