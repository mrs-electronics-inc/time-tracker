package headless

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ParseKeyMsg converts a key string to tea.KeyMsg
// Supports formats like: "j", "enter", "ctrl+c", "shift+tab"
func ParseKeyMsg(key string) tea.KeyMsg {
	// Check for modifier combinations
	if strings.Contains(key, "+") {
		parts := strings.SplitN(key, "+", 2)
		modifier := strings.ToLower(parts[0])
		baseKey := strings.ToLower(parts[1])

		switch modifier {
		case "ctrl":
			return tea.KeyMsg{Type: tea.KeyCtrlA + tea.KeyType(baseKey[0]-'a')}
		case "shift":
			if baseKey == "tab" {
				return tea.KeyMsg{Type: tea.KeyShiftTab}
			}
			// Shift + letter = uppercase
			if len(baseKey) == 1 {
				return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{rune(strings.ToUpper(baseKey)[0])}}
			}
		case "alt":
			msg := ParseKeyMsg(baseKey)
			msg.Alt = true
			return msg
		}
	}

	// Special keys
	switch strings.ToLower(key) {
	case "enter", "return":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "esc", "escape":
		return tea.KeyMsg{Type: tea.KeyEscape}
	case "tab":
		return tea.KeyMsg{Type: tea.KeyTab}
	case "backspace":
		return tea.KeyMsg{Type: tea.KeyBackspace}
	case "delete":
		return tea.KeyMsg{Type: tea.KeyDelete}
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp}
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	case "left":
		return tea.KeyMsg{Type: tea.KeyLeft}
	case "right":
		return tea.KeyMsg{Type: tea.KeyRight}
	case "home":
		return tea.KeyMsg{Type: tea.KeyHome}
	case "end":
		return tea.KeyMsg{Type: tea.KeyEnd}
	case "pgup", "pageup":
		return tea.KeyMsg{Type: tea.KeyPgUp}
	case "pgdown", "pagedown":
		return tea.KeyMsg{Type: tea.KeyPgDown}
	case "space":
		return tea.KeyMsg{Type: tea.KeySpace}
	}

	// Single character
	if len(key) == 1 {
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
	}

	// Unknown key - treat as runes
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
}

// ParseTypeToKeyMsgs converts a text string to a sequence of tea.KeyMsg
func ParseTypeToKeyMsgs(text string) []tea.KeyMsg {
	msgs := make([]tea.KeyMsg, len(text))
	for i, r := range text {
		msgs[i] = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
	}
	return msgs
}

// NewWindowSizeMsg creates a tea.WindowSizeMsg
func NewWindowSizeMsg(rows, cols int) tea.WindowSizeMsg {
	return tea.WindowSizeMsg{Width: cols, Height: rows}
}
