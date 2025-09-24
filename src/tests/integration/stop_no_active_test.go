package integration

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestStopNoActiveScenario(t *testing.T) {
	// Clean up
	os.Remove("../../data.json")

	// Build
	cmd := exec.Command("sh", "-c", "cd ../../ && go build -o time-tracker")
	buildOutput, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build binary: %v, output: %s", err, string(buildOutput))
	}

	// Try to stop when no active entry
	cmd = exec.Command("./time-tracker", "stop")
	cmd.Dir = "../../"
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected stop to fail")
	}

	outputStr := string(output)
	// Should show error message
	if !strings.Contains(outputStr, "no active time entry") {
		t.Errorf("Expected 'no active time entry', got: %s", outputStr)
	}
}
