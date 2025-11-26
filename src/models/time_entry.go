package models

import (
	"time"
)

// TimeEntry represents a period of tracked time
type TimeEntry struct {
	Start   time.Time  `json:"start"`
	End     *time.Time `json:"end,omitempty"`
	Project string     `json:"project"`
	Title   string     `json:"title"`
}

// SavedTimeEntry represents the format used to persist TimeEntry to disk (v3+)
// It excludes the ID field (auto-generated on load) and End field (reconstructed on load)
type SavedTimeEntry struct {
	Start   time.Time `json:"start"`
	Project string    `json:"project"`
	Title   string    `json:"title"`
}

// IsRunning returns true if the entry is currently active
func (te *TimeEntry) IsRunning() bool {
	return te.End == nil
}

// Duration returns the duration of the entry
func (te *TimeEntry) Duration() time.Duration {
	if te.End == nil {
		return time.Since(te.Start)
	}
	return te.End.Sub(te.Start)
}

// IsBlank returns true if the entry is a blank entry (empty project and title)
func (te *TimeEntry) IsBlank() bool {
	return te.Project == "" && te.Title == ""
}
