package unit

import (
	"testing"
	"time"

	"time-tracker/models"
)

func TestTimeEntry_IsRunning(t *testing.T) {
	// Running entry
	entry := models.TimeEntry{
		Start:   time.Now(),
		End:     nil,
		Project: "test",
		Title:   "test",
	}
	if !entry.IsRunning() {
		t.Error("Expected entry to be running")
	}

	// Stopped entry
	end := time.Now()
	entry.End = &end
	if entry.IsRunning() {
		t.Error("Expected entry to not be running")
	}
}

func TestTimeEntry_Duration(t *testing.T) {
	start := time.Now().Add(-time.Hour)
	entry := models.TimeEntry{
		Start: start,
		End:   nil,
	}
	duration := entry.Duration()
	if duration < time.Hour-time.Second || duration > time.Hour+time.Second {
		t.Errorf("Expected duration around 1 hour, got %v", duration)
	}

	// Stopped
	end := start.Add(time.Hour)
	entry.End = &end
	duration = entry.Duration()
	if duration != time.Hour {
		t.Errorf("Expected duration 1 hour, got %v", duration)
	}
}
