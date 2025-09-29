package integration

import (
	"os/exec"
	"testing"
)

func TestDockerCommands(t *testing.T) {
	// Test start command
	cmd := exec.Command("docker", "run", "--rm", "-v", "/tmp:/data", "time-tracker", "start", "test-task")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Docker start command failed: %v", err)
	}

	// Test stop command
	cmd = exec.Command("docker", "run", "--rm", "-v", "/tmp:/data", "time-tracker", "stop")
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Docker stop command failed: %v", err)
	}
}
