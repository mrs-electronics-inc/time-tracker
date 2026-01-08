//go:build integration

package cmd

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestTestModeIntegration tests the test subcommand end-to-end
// Run with: go test -tags=integration ./cmd -run TestTestModeIntegration
func TestTestModeIntegration(t *testing.T) {
	// Build the binary first
	buildCmd := exec.Command("go", "build", "-o", "time-tracker-test", "..")
	buildCmd.Dir = filepath.Join("..", "cmd")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build binary: %v", err)
	}
	defer os.Remove(filepath.Join("..", "cmd", "time-tracker-test"))

	// Create temp dirs
	renderDir, err := os.MkdirTemp("", "test-renders")
	if err != nil {
		t.Fatalf("failed to create render dir: %v", err)
	}
	defer os.RemoveAll(renderDir)

	configDir, err := os.MkdirTemp("", "test-config")
	if err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}
	defer os.RemoveAll(configDir)

	// Start test mode
	cmd := exec.Command(filepath.Join("..", "cmd", "time-tracker-test"), "test", "--render-dir", renderDir, "--keep-renders")
	cmd.Env = append(os.Environ(), "XDG_CONFIG_HOME="+configDir)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("failed to get stdin pipe: %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("failed to get stdout pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start test mode: %v", err)
	}

	defer func() {
		stdin.Close()
		cmd.Wait()
	}()

	reader := bufio.NewReader(stdout)

	// Read initial render
	resp, err := readResponse(reader)
	if err != nil {
		t.Fatalf("failed to read initial response: %v", err)
	}
	if resp.RenderPath == "" {
		t.Error("initial response has empty render_path")
	}
	if !strings.HasPrefix(resp.RenderPath, renderDir) {
		t.Errorf("render_path %s doesn't start with %s", resp.RenderPath, renderDir)
	}

	// Verify initial render file exists
	if _, err := os.Stat(resp.RenderPath); err != nil {
		t.Errorf("initial render file doesn't exist: %v", err)
	}

	// Send resize command
	writeCommand(stdin, Command{Cmd: "resize", Rows: 20, Cols: 60})
	resp, err = readResponse(reader)
	if err != nil {
		t.Fatalf("failed to read resize response: %v", err)
	}
	if resp.RenderPath == "" {
		t.Error("resize response has empty render_path")
	}

	// Send key command
	writeCommand(stdin, Command{Cmd: "key", Key: "?"})
	resp, err = readResponse(reader)
	if err != nil {
		t.Fatalf("failed to read key response: %v", err)
	}
	if resp.RenderPath == "" {
		t.Error("key response has empty render_path")
	}

	// Close stdin to end test mode
	stdin.Close()

	// Wait for process to exit
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Logf("process exited with: %v", err)
		}
	case <-time.After(5 * time.Second):
		cmd.Process.Kill()
		t.Fatal("test mode didn't exit within timeout")
	}
}

func writeCommand(w io.Writer, cmd Command) {
	data, _ := json.Marshal(cmd)
	w.Write(data)
	w.Write([]byte("\n"))
}

func readResponse(r *bufio.Reader) (Response, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return Response{}, err
	}

	var resp Response
	if err := json.Unmarshal([]byte(line), &resp); err != nil {
		return Response{}, err
	}

	return resp, nil
}
