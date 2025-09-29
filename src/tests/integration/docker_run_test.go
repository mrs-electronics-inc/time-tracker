package integration

import (
	"os/exec"
	"testing"
)

func TestDockerRun(t *testing.T) {
	cmd := exec.Command("docker", "run", "--rm", "time-tracker", "--help")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Docker run failed: %v", err)
	}
}
