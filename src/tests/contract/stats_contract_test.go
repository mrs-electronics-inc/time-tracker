package contract

import (
	"os/exec"
	"strings"
	"testing"
)

func TestStatsCommandOutputFormat(t *testing.T) {
	// Test that stats command produces expected output format
	// This test will fail until the stats command is implemented

	// Run the stats command with default flags
	cmd := exec.Command("go", "run", "main.go", "stats")
	cmd.Dir = "../../"
	output, err := cmd.CombinedOutput()

	// Now that implemented, expect it to succeed
	if err != nil {
		t.Errorf("Expected command to succeed, but failed: %v, output: %s", err, string(output))
	}

	// Check that output contains expected headers
	outputStr := string(output)
	if !strings.Contains(outputStr, "DATE") || !strings.Contains(outputStr, "TOTAL TIME") {
		t.Errorf("Expected output to contain headers, got: %s", outputStr)
	}
}
