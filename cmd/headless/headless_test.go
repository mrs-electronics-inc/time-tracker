package headless

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleState(t *testing.T) {
	server := NewServer(100)

	req := httptest.NewRequest(http.MethodGet, "/state", nil)
	w := httptest.NewRecorder()

	server.handleState(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp StateResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Width != 160 {
		t.Errorf("expected width 160, got %d", resp.Width)
	}
	if resp.Height != 40 {
		t.Errorf("expected height 40, got %d", resp.Height)
	}
}

func TestHandleStateMethodNotAllowed(t *testing.T) {
	server := NewServer(100)

	req := httptest.NewRequest(http.MethodPost, "/state", nil)
	w := httptest.NewRecorder()

	server.handleState(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}

func TestHandleInputValidActions(t *testing.T) {
	server := NewServer(100)

	tests := []struct {
		name string
		body InputRequest
	}{
		{"key action", InputRequest{Action: "key", Key: "j"}},
		{"type action", InputRequest{Action: "type", Text: "hello"}},
		{"resize action", InputRequest{Action: "resize", Rows: 24, Cols: 80}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/input", bytes.NewReader(body))
			w := httptest.NewRecorder()

			server.handleInput(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status 200, got %d", w.Code)
			}

			var resp StateResponse
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
		})
	}
}

func TestHandleInputInvalidAction(t *testing.T) {
	server := NewServer(100)

	body, _ := json.Marshal(InputRequest{Action: "invalid"})
	req := httptest.NewRequest(http.MethodPost, "/input", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.handleInput(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error != "invalid action: invalid" {
		t.Errorf("unexpected error message: %s", resp.Error)
	}
}

func TestHandleInputInvalidJSON(t *testing.T) {
	server := NewServer(100)

	req := httptest.NewRequest(http.MethodPost, "/input", bytes.NewReader([]byte("not json")))
	w := httptest.NewRecorder()

	server.handleInput(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandleInputMethodNotAllowed(t *testing.T) {
	server := NewServer(100)

	req := httptest.NewRequest(http.MethodGet, "/input", nil)
	w := httptest.NewRecorder()

	server.handleInput(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}

func TestHandleRenderLatestNoRenders(t *testing.T) {
	server := NewServer(100)

	req := httptest.NewRequest(http.MethodGet, "/render/latest", nil)
	w := httptest.NewRecorder()

	server.handleRenderLatest(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestHandleRenderLatestWithRender(t *testing.T) {
	server := NewServer(100)
	server.renders["2026-01-08T10-45-32-123"] = []byte("png data")
	server.latest = "2026-01-08T10-45-32-123"

	req := httptest.NewRequest(http.MethodGet, "/render/latest", nil)
	w := httptest.NewRecorder()

	server.handleRenderLatest(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("expected status 302, got %d", w.Code)
	}

	loc := w.Header().Get("Location")
	if loc != "/render/2026-01-08T10-45-32-123.png" {
		t.Errorf("unexpected location: %s", loc)
	}
}

func TestHandleRenderNotFound(t *testing.T) {
	server := NewServer(100)

	req := httptest.NewRequest(http.MethodGet, "/render/nonexistent.png", nil)
	w := httptest.NewRecorder()

	server.handleRender(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestHandleRenderFound(t *testing.T) {
	server := NewServer(100)
	server.renders["2026-01-08T10-45-32-123"] = []byte("fake png data")

	req := httptest.NewRequest(http.MethodGet, "/render/2026-01-08T10-45-32-123.png", nil)
	w := httptest.NewRecorder()

	server.handleRender(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "image/png" {
		t.Errorf("expected Content-Type image/png, got %s", w.Header().Get("Content-Type"))
	}

	if w.Body.String() != "fake png data" {
		t.Errorf("unexpected body: %s", w.Body.String())
	}
}
