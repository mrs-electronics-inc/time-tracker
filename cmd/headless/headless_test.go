package headless

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleState(t *testing.T) {
	// Skip test that requires full initialization
	t.Skip("requires TUI model initialization")
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

func TestHandleInputMethodNotAllowed(t *testing.T) {
	server := NewServer(100)

	req := httptest.NewRequest(http.MethodGet, "/input", nil)
	w := httptest.NewRecorder()

	server.handleInput(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
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

func TestHandleInputValidActions(t *testing.T) {
	// Skip - requires TUI model initialization
	t.Skip("requires TUI model initialization")
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
