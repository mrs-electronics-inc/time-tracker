package modes

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"time-tracker/models"
	"time-tracker/utils"
)

func newProjectSuggestionStorage(t *testing.T, projects []models.Project) *utils.MemoryStorage {
	t.Helper()

	storage := utils.NewMemoryStorage()
	if err := storage.SaveProjects(projects); err != nil {
		t.Fatalf("SaveProjects() error = %v", err)
	}
	return storage
}

func TestOpenFormModesSetProjectSuggestions(t *testing.T) {
	projects := []models.Project{{Name: " beta "}, {Name: "Alpha"}, {Name: "alpha"}, {Name: ""}, {Name: "Gamma"}}
	m := newFormDateTestModel()
	m.Storage = newProjectSuggestionStorage(t, projects)

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
	m.Storage = newProjectSuggestionStorage(t, []models.Project{{Name: "Backend"}})

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

	if updatedModel.Inputs[InputProject].Value() != "Backend" {
		t.Fatalf("project value = %q, expected suggestion to be accepted", updatedModel.Inputs[InputProject].Value())
	}
	if updatedModel.FocusIndex != InputProject {
		t.Fatalf("focus index = %d, expected project input to stay focused", updatedModel.FocusIndex)
	}
}

func TestFormModeTabNavigatesWhenProjectAlreadyExactlyMatchesSuggestion(t *testing.T) {
	m := newFormDateTestModel()
	m.Storage = newProjectSuggestionStorage(t, []models.Project{{Name: "Backend"}})

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
