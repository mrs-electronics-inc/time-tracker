package integration

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestAutoStopScenario(t *testing.T) {
	// Clean up
	os.Remove("../../data.json")

	// Build
	cmd := exec.Command("go", "build", "-o", "time-tracker")
	cmd.Dir = "../../"
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	// Start first task
	cmd = exec.Command("./time-tracker", "start", "project1", "Task 1")
	cmd.Dir = "../../"
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Start command 1 failed: %v", err)
	}

	// Start second task (should auto-stop first)
	cmd = exec.Command("./time-tracker", "start", "project2", "Task 2")
	cmd.Dir = "../../"
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Start command 2 failed: %v", err)
	}

	// List
	cmd = exec.Command("./time-tracker", "list")
	cmd.Dir = "../../"
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("List command failed: %v", err)
	}

	outputStr := string(output)
	// Should have both entries, first stopped, second running
	if !strings.Contains(outputStr, "project1") || !strings.Contains(outputStr, "project2") || !strings.Contains(outputStr, "running") {
		t.Errorf("Expected both entries with first stopped and second running, got: %s", outputStr)
	}
}
