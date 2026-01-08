package headless

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestParseKeyMsgSingleChar(t *testing.T) {
	tests := []struct {
		input    string
		expected tea.KeyMsg
	}{
		{"j", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}},
		{"k", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}},
		{"?", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}},
		{"s", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ParseKeyMsg(tt.input)
			if result.Type != tt.expected.Type {
				t.Errorf("expected type %v, got %v", tt.expected.Type, result.Type)
			}
			if string(result.Runes) != string(tt.expected.Runes) {
				t.Errorf("expected runes %v, got %v", tt.expected.Runes, result.Runes)
			}
		})
	}
}

func TestParseKeyMsgSpecialKeys(t *testing.T) {
	tests := []struct {
		input    string
		expected tea.KeyType
	}{
		{"enter", tea.KeyEnter},
		{"Enter", tea.KeyEnter},
		{"esc", tea.KeyEscape},
		{"escape", tea.KeyEscape},
		{"tab", tea.KeyTab},
		{"backspace", tea.KeyBackspace},
		{"up", tea.KeyUp},
		{"down", tea.KeyDown},
		{"left", tea.KeyLeft},
		{"right", tea.KeyRight},
		{"space", tea.KeySpace},
		{"home", tea.KeyHome},
		{"end", tea.KeyEnd},
		{"pgup", tea.KeyPgUp},
		{"pgdown", tea.KeyPgDown},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ParseKeyMsg(tt.input)
			if result.Type != tt.expected {
				t.Errorf("expected type %v, got %v", tt.expected, result.Type)
			}
		})
	}
}

func TestParseKeyMsgModifiers(t *testing.T) {
	// shift+tab
	result := ParseKeyMsg("shift+tab")
	if result.Type != tea.KeyShiftTab {
		t.Errorf("expected KeyShiftTab, got %v", result.Type)
	}

	// ctrl+c
	result = ParseKeyMsg("ctrl+c")
	if result.Type != tea.KeyCtrlC {
		t.Errorf("expected KeyCtrlC, got %v", result.Type)
	}

	// alt+j
	result = ParseKeyMsg("alt+j")
	if !result.Alt {
		t.Error("expected Alt to be true")
	}
	if string(result.Runes) != "j" {
		t.Errorf("expected rune 'j', got %v", result.Runes)
	}
}

func TestParseTypeToKeyMsgs(t *testing.T) {
	msgs := ParseTypeToKeyMsgs("hello")

	if len(msgs) != 5 {
		t.Fatalf("expected 5 messages, got %d", len(msgs))
	}

	expected := []rune{'h', 'e', 'l', 'l', 'o'}
	for i, msg := range msgs {
		if msg.Type != tea.KeyRunes {
			t.Errorf("message %d: expected KeyRunes, got %v", i, msg.Type)
		}
		if len(msg.Runes) != 1 || msg.Runes[0] != expected[i] {
			t.Errorf("message %d: expected rune %c, got %v", i, expected[i], msg.Runes)
		}
	}
}

func TestNewWindowSizeMsg(t *testing.T) {
	msg := NewWindowSizeMsg(24, 80)

	if msg.Height != 24 {
		t.Errorf("expected height 24, got %d", msg.Height)
	}
	if msg.Width != 80 {
		t.Errorf("expected width 80, got %d", msg.Width)
	}
}
