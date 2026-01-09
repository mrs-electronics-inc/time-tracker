package headless

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var specialKeys = map[string]tea.KeyType{
	"enter": tea.KeyEnter, "return": tea.KeyEnter,
	"esc": tea.KeyEscape, "escape": tea.KeyEscape,
	"tab": tea.KeyTab, "backspace": tea.KeyBackspace, "delete": tea.KeyDelete,
	"up": tea.KeyUp, "down": tea.KeyDown, "left": tea.KeyLeft, "right": tea.KeyRight,
	"home": tea.KeyHome, "end": tea.KeyEnd,
	"pgup": tea.KeyPgUp, "pageup": tea.KeyPgUp,
	"pgdown": tea.KeyPgDown, "pagedown": tea.KeyPgDown,
	"space":     tea.KeySpace,
	"shift+tab": tea.KeyShiftTab,
}

var ctrlKeys = map[byte]tea.KeyType{
	'a': tea.KeyCtrlA, 'b': tea.KeyCtrlB, 'c': tea.KeyCtrlC, 'd': tea.KeyCtrlD,
	'e': tea.KeyCtrlE, 'f': tea.KeyCtrlF, 'g': tea.KeyCtrlG, 'h': tea.KeyCtrlH,
	'i': tea.KeyCtrlI, 'j': tea.KeyCtrlJ, 'k': tea.KeyCtrlK, 'l': tea.KeyCtrlL,
	'm': tea.KeyCtrlM, 'n': tea.KeyCtrlN, 'o': tea.KeyCtrlO, 'p': tea.KeyCtrlP,
	'q': tea.KeyCtrlQ, 'r': tea.KeyCtrlR, 's': tea.KeyCtrlS, 't': tea.KeyCtrlT,
	'u': tea.KeyCtrlU, 'v': tea.KeyCtrlV, 'w': tea.KeyCtrlW, 'x': tea.KeyCtrlX,
	'y': tea.KeyCtrlY, 'z': tea.KeyCtrlZ,
}

// ParseKeyMsg converts a key string to tea.KeyMsg
// Supports formats like: "j", "enter", "ctrl+c", "shift+tab", "alt+x"
func ParseKeyMsg(key string) tea.KeyMsg {
	normalized := strings.ToLower(key)

	if keyType, ok := specialKeys[normalized]; ok {
		return tea.KeyMsg{Type: keyType}
	}

	if strings.HasPrefix(normalized, "ctrl+") {
		baseKey := normalized[5:]
		if len(baseKey) == 1 {
			if keyType, ok := ctrlKeys[baseKey[0]]; ok {
				return tea.KeyMsg{Type: keyType}
			}
		}
	}

	if strings.HasPrefix(normalized, "shift+") {
		baseKey := normalized[6:]
		if len(baseKey) == 1 {
			return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{rune(strings.ToUpper(baseKey)[0])}}
		}
	}

	if strings.HasPrefix(normalized, "alt+") {
		msg := ParseKeyMsg(normalized[4:])
		msg.Alt = true
		return msg
	}

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
