package contract

import (
	"os/exec"
	"strings"
	"testing"
)

func TestStopCommand(t *testing.T) {
	// Build the binary
	cmd := exec.Command("go", "build", "-o", "time-tracker")
	cmd.Dir = "../../src"
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	// First start a task
	cmd = exec.Command("./time-tracker", "start", "test-project", "Test task")
	cmd.Dir = "../../src"
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to start task: %v", err)
	}

	// Run stop command
	cmd = exec.Command("./time-tracker", "stop")
	cmd.Dir = "../../src"
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Stop command failed: %v", err)
	}

	// Check output contains confirmation
	outputStr := string(output)
	if !strings.Contains(outputStr, "Stopped tracking time") {
		t.Errorf("Expected 'Stopped tracking time', got: %s", outputStr)
	}
}
