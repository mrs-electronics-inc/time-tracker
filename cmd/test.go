package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"time-tracker/cmd/tui"
	"time-tracker/config"
	"time-tracker/utils"
)

// Command represents an input command from stdin
type Command struct {
	Cmd  string `json:"cmd"`
	Key  string `json:"key,omitempty"`
	Text string `json:"text,omitempty"`
	Rows int    `json:"rows,omitempty"`
	Cols int    `json:"cols,omitempty"`
}

// Response represents the JSON response sent to stdout
type Response struct {
	RenderPath string `json:"render_path,omitempty"`
	Error      string `json:"error,omitempty"`
}

var (
	renderDir   string
	keepRenders bool
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run TUI in test mode for automated testing",
	Long: `Run the TUI in test mode, accepting JSON commands on stdin and
rendering the screen to PNG images after each command.

This enables automated testing and AI agent interaction with the TUI.`,
	RunE: runTestMode,
}

func init() {
	testCmd.Flags().StringVar(&renderDir, "render-dir", "/tmp/time-tracker/renders", "Directory to save render images")
	testCmd.Flags().BoolVar(&keepRenders, "keep-renders", false, "Keep render images after exit")
	rootCmd.AddCommand(testCmd)
}

func runTestMode(cmd *cobra.Command, args []string) error {
	// Create render directory
	if err := os.MkdirAll(renderDir, 0755); err != nil {
		return fmt.Errorf("failed to create render directory: %w", err)
	}

	// Track created renders for cleanup
	var createdRenders []string

	// Setup cleanup on exit
	if !keepRenders {
		cleanup := func() {
			for _, path := range createdRenders {
				os.Remove(path)
			}
		}
		defer cleanup()

		// Handle signals for cleanup
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigChan
			cleanup()
			os.Exit(0)
		}()
	}

	// Initialize storage and model
	storage, err := utils.NewFileStorage(config.DataFilePath())
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	taskManager := utils.NewTaskManager(storage)

	model := tui.NewModel(storage, taskManager)
	if err := model.LoadEntries(); err != nil {
		return fmt.Errorf("failed to load entries: %w", err)
	}

	// Set default terminal size
	model.Width = 80
	model.Height = 24

	// Initialize the model
	model.Init()

	// Create renderer
	renderer, err := NewTerminalRenderer()
	if err != nil {
		return fmt.Errorf("failed to create renderer: %w", err)
	}

	// Helper to render and send response
	sendRender := func() error {
		view := model.View()
		timestamp := time.Now().Format("2006-01-02T15-04-05.000")
		filename := filepath.Join(renderDir, timestamp+".png")

		if err := renderer.RenderToFile(view, model.Width, model.Height, filename); err != nil {
			return err
		}
		createdRenders = append(createdRenders, filename)

		resp := Response{RenderPath: filename}
		return json.NewEncoder(os.Stdout).Encode(resp)
	}

	// Send initial render
	if err := sendRender(); err != nil {
		return fmt.Errorf("failed to send initial render: %w", err)
	}

	// Process commands from stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var cmdInput Command
		if err := json.Unmarshal([]byte(line), &cmdInput); err != nil {
			resp := Response{Error: fmt.Sprintf("invalid JSON: %v", err)}
			json.NewEncoder(os.Stdout).Encode(resp)
			continue
		}

		// Convert command to tea.Msg and update model
		switch cmdInput.Cmd {
		case "key":
			msg := keyToMsg(cmdInput.Key)
			if msg.Type == tea.KeyRunes && len(msg.Runes) == 0 && msg.Alt == false {
				resp := Response{Error: fmt.Sprintf("unknown key: %s", cmdInput.Key)}
				json.NewEncoder(os.Stdout).Encode(resp)
				continue
			}
			newModel, _ := model.Update(msg)
			model = newModel.(*tui.Model)

		case "type":
			for _, r := range cmdInput.Text {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
				newModel, _ := model.Update(msg)
				model = newModel.(*tui.Model)
			}

		case "resize":
			if cmdInput.Rows > 0 {
				model.Height = cmdInput.Rows
			}
			if cmdInput.Cols > 0 {
				model.Width = cmdInput.Cols
			}
			msg := tea.WindowSizeMsg{Width: model.Width, Height: model.Height}
			newModel, _ := model.Update(msg)
			model = newModel.(*tui.Model)

		default:
			resp := Response{Error: fmt.Sprintf("unknown command: %s", cmdInput.Cmd)}
			json.NewEncoder(os.Stdout).Encode(resp)
			continue
		}

		// Render and send response
		if err := sendRender(); err != nil {
			resp := Response{Error: fmt.Sprintf("render failed: %v", err)}
			json.NewEncoder(os.Stdout).Encode(resp)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stdin: %w", err)
	}

	return nil
}

// keyToMsg converts a key string to a tea.KeyMsg
func keyToMsg(key string) tea.KeyMsg {
	// Handle special keys
	switch key {
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "esc", "escape":
		return tea.KeyMsg{Type: tea.KeyEscape}
	case "tab":
		return tea.KeyMsg{Type: tea.KeyTab}
	case "shift+tab":
		return tea.KeyMsg{Type: tea.KeyShiftTab}
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp}
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	case "left":
		return tea.KeyMsg{Type: tea.KeyLeft}
	case "right":
		return tea.KeyMsg{Type: tea.KeyRight}
	case "backspace":
		return tea.KeyMsg{Type: tea.KeyBackspace}
	case "delete":
		return tea.KeyMsg{Type: tea.KeyDelete}
	case "home":
		return tea.KeyMsg{Type: tea.KeyHome}
	case "end":
		return tea.KeyMsg{Type: tea.KeyEnd}
	case "pgup", "pageup":
		return tea.KeyMsg{Type: tea.KeyPgUp}
	case "pgdown", "pagedown":
		return tea.KeyMsg{Type: tea.KeyPgDown}
	case "space":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}}
	case "ctrl+c":
		return tea.KeyMsg{Type: tea.KeyCtrlC}
	case "ctrl+d":
		return tea.KeyMsg{Type: tea.KeyCtrlD}
	case "ctrl+z":
		return tea.KeyMsg{Type: tea.KeyCtrlZ}
	}

	// Single character keys
	if len(key) == 1 {
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
	}

	// Unknown key
	return tea.KeyMsg{Type: tea.KeyRunes}
}
