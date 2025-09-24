package integration

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestStopNoActiveScenario(t *testing.T) {
	// Clean up
	os.Remove("../../src/data.json")

	// Build
	cmd := exec.Command("go", "build", "-o", "time-tracker")
	cmd.Dir = "../../src"
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	// Try to stop when no active entry
	cmd = exec.Command("./time-tracker", "stop")
	cmd.Dir = "../../src"
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Stop command failed: %v", err)
	}

	outputStr := string(output)
	// Should show error message
	if !strings.Contains(outputStr, "No active time entry") {
		t.Errorf("Expected 'No active time entry', got: %s", outputStr)
	}
}
