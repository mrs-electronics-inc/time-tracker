package cmd

import (
	"encoding/json"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestCommandParsing(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		want    Command
	}{
		{
			name:    "key command",
			jsonStr: `{"cmd": "key", "key": "j"}`,
			want:    Command{Cmd: "key", Key: "j"},
		},
		{
			name:    "type command",
			jsonStr: `{"cmd": "type", "text": "hello world"}`,
			want:    Command{Cmd: "type", Text: "hello world"},
		},
		{
			name:    "resize command",
			jsonStr: `{"cmd": "resize", "rows": 24, "cols": 80}`,
			want:    Command{Cmd: "resize", Rows: 24, Cols: 80},
		},
		{
			name:    "key enter",
			jsonStr: `{"cmd": "key", "key": "enter"}`,
			want:    Command{Cmd: "key", Key: "enter"},
		},
		{
			name:    "key escape",
			jsonStr: `{"cmd": "key", "key": "esc"}`,
			want:    Command{Cmd: "key", Key: "esc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Command
			err := json.Unmarshal([]byte(tt.jsonStr), &got)
			if err != nil {
				t.Fatalf("failed to parse JSON: %v", err)
			}

			if got.Cmd != tt.want.Cmd {
				t.Errorf("Cmd = %q, want %q", got.Cmd, tt.want.Cmd)
			}
			if got.Key != tt.want.Key {
				t.Errorf("Key = %q, want %q", got.Key, tt.want.Key)
			}
			if got.Text != tt.want.Text {
				t.Errorf("Text = %q, want %q", got.Text, tt.want.Text)
			}
			if got.Rows != tt.want.Rows {
				t.Errorf("Rows = %d, want %d", got.Rows, tt.want.Rows)
			}
			if got.Cols != tt.want.Cols {
				t.Errorf("Cols = %d, want %d", got.Cols, tt.want.Cols)
			}
		})
	}
}

func TestKeyToMsg(t *testing.T) {
	tests := []struct {
		key      string
		wantType tea.KeyType
		wantRune rune
	}{
		{"enter", tea.KeyEnter, 0},
		{"esc", tea.KeyEscape, 0},
		{"escape", tea.KeyEscape, 0},
		{"tab", tea.KeyTab, 0},
		{"shift+tab", tea.KeyShiftTab, 0},
		{"up", tea.KeyUp, 0},
		{"down", tea.KeyDown, 0},
		{"left", tea.KeyLeft, 0},
		{"right", tea.KeyRight, 0},
		{"backspace", tea.KeyBackspace, 0},
		{"delete", tea.KeyDelete, 0},
		{"home", tea.KeyHome, 0},
		{"end", tea.KeyEnd, 0},
		{"pgup", tea.KeyPgUp, 0},
		{"pageup", tea.KeyPgUp, 0},
		{"pgdown", tea.KeyPgDown, 0},
		{"pagedown", tea.KeyPgDown, 0},
		{"ctrl+c", tea.KeyCtrlC, 0},
		{"ctrl+d", tea.KeyCtrlD, 0},
		{"j", tea.KeyRunes, 'j'},
		{"k", tea.KeyRunes, 'k'},
		{"q", tea.KeyRunes, 'q'},
		{"space", tea.KeyRunes, ' '},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			msg := keyToMsg(tt.key)

			if msg.Type != tt.wantType {
				t.Errorf("Type = %v, want %v", msg.Type, tt.wantType)
			}

			if tt.wantRune != 0 {
				if len(msg.Runes) == 0 || msg.Runes[0] != tt.wantRune {
					t.Errorf("Runes = %v, want [%c]", msg.Runes, tt.wantRune)
				}
			}
		})
	}
}

func TestResponseJSON(t *testing.T) {
	resp := Response{RenderPath: "/tmp/time-tracker/renders/test.png"}
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal response: %v", err)
	}

	expected := `{"render_path":"/tmp/time-tracker/renders/test.png"}`
	if string(data) != expected {
		t.Errorf("got %s, want %s", string(data), expected)
	}

	errResp := Response{Error: "something went wrong"}
	data, err = json.Marshal(errResp)
	if err != nil {
		t.Fatalf("failed to marshal error response: %v", err)
	}

	expected = `{"error":"something went wrong"}`
	if string(data) != expected {
		t.Errorf("got %s, want %s", string(data), expected)
	}
}
