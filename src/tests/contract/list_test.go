package contract

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestListCommand(t *testing.T) {
	// Clean up
	os.Remove("../../data.json")
	os.Remove("../../data.json")
os.Remove("../../time-tracker")

	// Build the binary first
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
	if !strings.Contains(outputStr, "ID") && strings.Contains(outputStr, "Start") && strings.Contains(outputStr, "test-project") {
		t.Errorf("Expected table with entry, got: %s", outputStr)
	}
}
