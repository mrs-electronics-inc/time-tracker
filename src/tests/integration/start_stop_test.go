package integration

import (
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestStartStopScenario(t *testing.T) {
	// Clean up any existing data.json
	os.Remove("../../data.json")

	// Build the binary
	os.Remove("../../time-tracker")
	cmd := exec.Command("go", "build", "-o", "time-tracker")
	cmd.Dir = "../../"
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	// Start tracking
	cmd = exec.Command("./time-tracker", "start", "test-project", "Test task")
	cmd.Dir = "../../"
	_, err = cmd.Output()
	if err != nil {
		t.Fatalf("Start command failed: %v", err)
	}

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// Stop tracking
	cmd = exec.Command("./time-tracker", "stop")
	cmd.Dir = "../../"
	stopOutput, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Stop command failed: %v, output: %s", err, string(stopOutput))
	}

	// List entries
	cmd = exec.Command("./time-tracker", "list")
	cmd.Dir = "../../"
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("List command failed: %v", err)
	}

	outputStr := string(output)
	// Should contain the entry with duration
	if !strings.Contains(outputStr, "test-project") || !strings.Contains(outputStr, "Test task") {
		t.Errorf("Expected entry in list, got: %s", outputStr)
	}
}
