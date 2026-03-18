package modes

import (
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"time-tracker/utils"
)

func newFormTimeTestModel() *Model {
	inputs := make([]textinput.Model, InputDay+1)
	for i := range inputs {
		inputs[i] = textinput.New()
	}

	storage := utils.NewMemoryStorage()

	return &Model{
		Storage:     storage,
		TaskManager: utils.NewTaskManager(storage),
		Inputs:      inputs,
		StartMode:   StartMode,
		ListMode:    ListMode,
	}
}

func TestParseFormTimeUsesDateFields(t *testing.T) {
	m := newFormTimeTestModel()
	m.Inputs[InputYear].SetValue("2024")
	m.Inputs[InputMonth].SetValue("02")
	m.Inputs[InputDay].SetValue("29")
	m.Inputs[InputHour].SetValue("23")
	m.Inputs[InputMinute].SetValue("59")

	parsed, err := parseFormTime(m)
	if err != nil {
		t.Fatalf("parseFormTime returned error: %v", err)
	}

	if parsed.Year() != 2024 || parsed.Month() != time.February || parsed.Day() != 29 {
		t.Fatalf("expected parsed date 2024-02-29, got %s", parsed.Format("2006-01-02 15:04"))
	}
	if parsed.Hour() != 23 || parsed.Minute() != 59 {
		t.Fatalf("expected parsed time 23:59, got %02d:%02d", parsed.Hour(), parsed.Minute())
	}
}

func TestParseFormTimeRejectsInvalidCalendarDate(t *testing.T) {
	m := newFormTimeTestModel()
	m.Inputs[InputYear].SetValue("2025")
	m.Inputs[InputMonth].SetValue("02")
	m.Inputs[InputDay].SetValue("29")
	m.Inputs[InputHour].SetValue("09")
	m.Inputs[InputMinute].SetValue("30")

	_, err := parseFormTime(m)
	if err == nil {
		t.Fatal("expected parseFormTime to reject invalid date 2025-02-29")
	}

	if !strings.Contains(strings.ToLower(err.Error()), "day") {
		t.Fatalf("expected day-related validation error, got %q", err.Error())
	}
}

func TestStartModeHandleEnterRejectsInvalidCalendarDate(t *testing.T) {
	m := newFormTimeTestModel()
	m.CurrentMode = m.StartMode
	m.Inputs[InputProject].SetValue("backend")
	m.Inputs[InputTitle].SetValue("review")
	m.Inputs[InputYear].SetValue("2025")
	m.Inputs[InputMonth].SetValue("02")
	m.Inputs[InputDay].SetValue("29")
	m.Inputs[InputHour].SetValue("09")
	m.Inputs[InputMinute].SetValue("30")

	updatedModel, _ := StartMode.HandleKeyMsg(m, tea.KeyMsg{Type: tea.KeyEnter})
	m = updatedModel

	if !strings.Contains(strings.ToLower(m.Status), "invalid") {
		t.Fatalf("expected invalid date error status, got %q", m.Status)
	}

	entries, err := m.Storage.Load()
	if err != nil {
		t.Fatalf("failed to load entries: %v", err)
	}

	if len(entries) != 0 {
		t.Fatalf("expected no entries saved on invalid date, got %d", len(entries))
	}
}
