package models

import (
	"time"
)

// TimeEntry represents a period of tracked time
type TimeEntry struct {
	ID      int        `json:"id"`
	Start   time.Time  `json:"start"`
	End     *time.Time `json:"end,omitempty"`
	Project string     `json:"project"`
	Title   string     `json:"title"`
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
