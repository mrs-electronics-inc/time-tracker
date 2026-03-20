package modes

import (
	"regexp"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
)

func newDateRenderTestModel() *Model {
	inputs := make([]textinput.Model, InputMinute+1)
	for i := range inputs {
		inputs[i] = textinput.New()
		inputs[i].Prompt = ""
	}

	inputs[InputProject].SetValue("backend")
	inputs[InputTitle].SetValue("review")
	inputs[InputYear].SetValue("2025")
	inputs[InputMonth].SetValue("01")
	inputs[InputDay].SetValue("07")
	inputs[InputHour].SetValue("09")
	inputs[InputMinute].SetValue("30")

	return &Model{Inputs: inputs}
}

func assertDateRendersAboveTime(t *testing.T, content string) {
	t.Helper()

	if !strings.Contains(content, "Date (YYYY-MM-DD):") {
		t.Fatalf("expected date label to be rendered, got:\n%s", content)
	}

	datePattern := regexp.MustCompile(`2025\s*-\s*01\s*-\s*07`)
	if !datePattern.MatchString(content) {
		t.Fatalf("expected date row to be rendered as YYYY - MM - DD, got:\n%s", content)
	}

	dateIndex := strings.Index(content, "Date (YYYY-MM-DD):")
	timeIndex := strings.Index(content, "Time (HH:MM):")
	if dateIndex == -1 || timeIndex == -1 || dateIndex > timeIndex {
		t.Fatalf("expected date row to render above time row, got:\n%s", content)
	}
}

func TestRenderFormContentRendersDateAboveTime(t *testing.T) {
	m := newDateRenderTestModel()

	content := renderFormContent(m, "New Entry", 10)

	assertDateRendersAboveTime(t, content)
}
