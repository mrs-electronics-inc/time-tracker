package integration

import (
	"strings"
	"testing"

	"time-tracker/utils"
)

func TestStopNoActiveScenario(t *testing.T) {
	// Create task manager with memory storage
	storage := utils.NewMemoryStorage()
	tm := utils.NewTaskManager(storage)

	// Try to stop when no active entry
	_, err := tm.StopEntry()
	if err == nil {
		t.Fatalf("Expected stop to fail")
	}

	// Check error message
	if !strings.Contains(err.Error(), "no active time entry") {
		t.Errorf("Expected 'no active time entry', got: %s", err.Error())
	}
}
