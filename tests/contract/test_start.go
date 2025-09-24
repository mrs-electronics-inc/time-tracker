package contract

import (
	"os/exec"
	"strings"
	"testing"
)

func TestStartCommand(t *testing.T) {
	// Build the binary
	cmd := exec.Command("go", "build", "-o", "time-tracker")
	cmd.Dir = "../../src"
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	// Run start command
	cmd = exec.Command("./time-tracker", "start", "test-project", "Test task")
	cmd.Dir = "../../src"
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Start command failed: %v", err)
	}

	// Check output contains confirmation
	outputStr := string(output)
	if !strings.Contains(outputStr, "Started tracking time") {
		t.Errorf("Expected 'Started tracking time', got: %s", outputStr)
	}
}
