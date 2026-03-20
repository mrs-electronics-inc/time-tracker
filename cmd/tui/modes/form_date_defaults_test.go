package modes

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
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
