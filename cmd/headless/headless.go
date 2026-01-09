package headless

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
	"time-tracker/cmd/tui"
	"time-tracker/config"
	"time-tracker/utils"
)

var (
	bindAddr   string
	port       int
	maxRenders int
)

const (
	defaultWidth  = 160
	defaultHeight = 40
)

// Server holds the headless server state
type Server struct {
	mu         sync.RWMutex
	renders    map[string][]byte // timestamp -> PNG data
	renderKeys []string          // ordered keys for FIFO eviction
	latest     string            // latest render timestamp
	maxRenders int
	renderer   *Renderer
	model      *tui.Model
	width      int
	height     int
	ansi       string
}

// NewServer creates a new headless server
func NewServer(maxRenders int) *Server {
	return &Server{
		renders:    make(map[string][]byte),
		renderKeys: make([]string, 0),
		maxRenders: maxRenders,
		width:      defaultWidth,
		height:     defaultHeight,
	}
}

// Initialize sets up the TUI model and renderer
func (s *Server) Initialize() error {
	// Force ANSI color output
	lipgloss.SetColorProfile(termenv.ANSI)

	// Create storage and task manager
	storage, err := utils.NewFileStorage(config.DataFilePath())
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	taskManager := utils.NewTaskManager(storage)

	// Create TUI model
	s.model = tui.NewModel(storage, taskManager)
	if err := s.model.LoadEntries(); err != nil {
		return fmt.Errorf("failed to load entries: %w", err)
	}

	// Send initial window size
	updated, _ := s.model.Update(tea.WindowSizeMsg{Width: s.width, Height: s.height})
	s.model = updated.(*tui.Model)

	// Create renderer
	s.renderer, err = NewRenderer(s.width, s.height)
	if err != nil {
		return fmt.Errorf("failed to create renderer: %w", err)
	}

	// Create initial render
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateRender()
}

// updateRender captures the current view and creates a PNG render
// Must be called with s.mu held
func (s *Server) updateRender() error {
	s.ansi = s.model.View()

	pngData, err := s.renderer.Render(s.ansi)
	if err != nil {
		return err
	}

	timestamp := time.Now().Format("2006-01-02T15-04-05-000")
	s.addRenderLocked(timestamp, pngData)
	return nil
}

// AddRender stores a render and evicts old ones if needed (FIFO)
func (s *Server) AddRender(timestamp string, data []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.addRenderLocked(timestamp, data)
}

// addRenderLocked stores a render (must be called with s.mu held)
func (s *Server) addRenderLocked(timestamp string, data []byte) {
	// Evict oldest if at capacity
	for len(s.renderKeys) >= s.maxRenders {
		oldest := s.renderKeys[0]
		delete(s.renders, oldest)
		s.renderKeys = s.renderKeys[1:]
	}

	s.renders[timestamp] = data
	s.renderKeys = append(s.renderKeys, timestamp)
	s.latest = timestamp
}

// InputRequest represents an input action
type InputRequest struct {
	Action string `json:"action"` // "key", "type", "resize"
	Key    string `json:"key"`    // for "key" action
	Text   string `json:"text"`   // for "type" action
	Rows   int    `json:"rows"`   // for "resize" action
	Cols   int    `json:"cols"`   // for "resize" action
}

// StateResponse represents the server state
type StateResponse struct {
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Mode      string `json:"mode"`
	RenderURL string `json:"render_url"`
	ANSI      string `json:"ansi"`
}

// ErrorResponse represents an error
type ErrorResponse struct {
	Error string `json:"error"`
}

func (s *Server) getState() StateResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	renderURL := ""
	if s.latest != "" {
		renderURL = "/render/" + s.latest + ".png"
	}

	return StateResponse{
		Width:     s.width,
		Height:    s.height,
		Mode:      s.model.CurrentMode.Name,
		RenderURL: renderURL,
		ANSI:      s.ansi,
	}
}

func (s *Server) handleInput(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "method not allowed"})
		return
	}

	var req InputRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid JSON: " + err.Error()})
		return
	}

	s.mu.Lock()
	var msgs []tea.Msg

	switch req.Action {
	case "key":
		msgs = []tea.Msg{ParseKeyMsg(req.Key)}
	case "type":
		for _, km := range ParseTypeToKeyMsgs(req.Text) {
			msgs = append(msgs, km)
		}
	case "resize":
		s.width = req.Cols
		s.height = req.Rows
		msgs = []tea.Msg{NewWindowSizeMsg(req.Rows, req.Cols)}
		// Recreate renderer with new size
		if renderer, err := NewRenderer(s.width, s.height); err == nil {
			s.renderer = renderer
		}
	default:
		s.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid action: " + req.Action})
		return
	}

	// Process all messages
	for _, msg := range msgs {
		updated, _ := s.model.Update(msg)
		s.model = updated.(*tui.Model)
	}

	// Update render
	s.updateRender()
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.getState())
}

func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "method not allowed"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.getState())
}

func (s *Server) handleRenderLatest(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	latest := s.latest
	s.mu.RUnlock()

	if latest == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "no renders available"})
		return
	}

	http.Redirect(w, r, "/render/"+latest+".png", http.StatusFound)
}

func (s *Server) handleRender(w http.ResponseWriter, r *http.Request) {
	// Extract timestamp from path: /render/{timestamp}.png
	path := r.URL.Path
	if len(path) < 12 || path[len(path)-4:] != ".png" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid render path"})
		return
	}

	timestamp := path[8 : len(path)-4] // strip "/render/" and ".png"

	s.mu.RLock()
	data, ok := s.renders[timestamp]
	s.mu.RUnlock()

	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "render not found"})
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(data)
}

// HeadlessCmd is the cobra command for headless mode
var HeadlessCmd = &cobra.Command{
	Use:   "headless",
	Short: "Run TUI as HTTP server for programmatic interaction",
	Long: `Run the TUI as an HTTP server, enabling AI agents and automated tests
to interact programmatically.

Endpoints:
  POST /input         Send input action, receive updated state
  GET  /state         Get current state
  GET  /render/latest Redirect to most recent render
  GET  /render/{ts}.png  Get specific render`,
	RunE: func(cmd *cobra.Command, args []string) error {
		server := NewServer(maxRenders)

		if err := server.Initialize(); err != nil {
			return err
		}

		mux := http.NewServeMux()
		mux.HandleFunc("/input", server.handleInput)
		mux.HandleFunc("/state", server.handleState)
		mux.HandleFunc("/render/latest", server.handleRenderLatest)
		mux.HandleFunc("/render/", server.handleRender)

		addr := fmt.Sprintf("%s:%d", bindAddr, port)
		fmt.Printf("Starting headless server on %s\n", addr)

		return http.ListenAndServe(addr, mux)
	},
}

func init() {
	HeadlessCmd.Flags().StringVar(&bindAddr, "bind", "127.0.0.1", "Bind address")
	HeadlessCmd.Flags().IntVar(&port, "port", 8484, "Port number")
	HeadlessCmd.Flags().IntVar(&maxRenders, "max-renders", 100, "Max renders to keep in memory")
}
