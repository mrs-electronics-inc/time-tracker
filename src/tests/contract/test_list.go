package contract

import (
	"os/exec"
	"strings"
	"testing"
)

func TestListCommand(t *testing.T) {
	// Build the binary first
	cmd := exec.Command("go", "build", "-o", "time-tracker")
	cmd.Dir = "../../src"
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	// Run list command
	cmd = exec.Command("./time-tracker", "list")
	cmd.Dir = "../../src"
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("List command failed: %v", err)
	}

	// Check output format (should show table header)
	outputStr := string(output)
	if !strings.Contains(outputStr, "ID") || !strings.Contains(outputStr, "Start") {
		t.Errorf("Expected table header with ID and Start, got: %s", outputStr)
	}
}
