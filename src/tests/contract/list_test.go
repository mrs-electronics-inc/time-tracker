package contract

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestListCommand(t *testing.T) {
	// Build the binary first
	os.Remove("../../time-tracker")
	cmd := exec.Command("go", "build", "-o", "time-tracker")
	cmd.Dir = "../../"
	buildOutput, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build binary: %v, output: %s", err, string(buildOutput))
	}

	// Run list command
	cmd = exec.Command("./time-tracker", "list")
	cmd.Dir = "../../"
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("List command failed: %v", err)
	}

	// Check output (should show no entries message)
	outputStr := string(output)
	if !strings.Contains(outputStr, "No time entries found") {
		t.Errorf("Expected 'No time entries found', got: %s", outputStr)
	}
}
