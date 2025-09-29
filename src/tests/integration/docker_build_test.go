package integration

import (
	"os/exec"
	"testing"
)

func TestDockerBuild(t *testing.T) {
	cmd := exec.Command("docker", "build", "-t", "time-tracker", ".")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Docker build failed: %v", err)
	}
}
