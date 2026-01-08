package headless

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/spf13/cobra"
)

var (
	bindAddr   string
	port       int
	maxRenders int
)

// Server holds the headless server state
type Server struct {
	mu      sync.RWMutex
	renders map[string][]byte // timestamp -> PNG data
	latest  string            // latest render timestamp
}

// NewServer creates a new headless server
func NewServer(maxRenders int) *Server {
	return &Server{
		renders: make(map[string][]byte),
	}
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

	// Validate action
	switch req.Action {
	case "key", "type", "resize":
		// valid
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid action: " + req.Action})
		return
	}

	// TODO: Actually process input and update TUI

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(StateResponse{
		Width:     160,
		Height:    40,
		Mode:      "list",
		RenderURL: "/render/latest",
		ANSI:      "",
	})
}

func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "method not allowed"})
		return
	}

	// TODO: Return actual state
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(StateResponse{
		Width:     160,
		Height:    40,
		Mode:      "list",
		RenderURL: "/render/latest",
		ANSI:      "",
	})
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
