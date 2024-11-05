package models

import "time"

type TaskStatus string

const (
	StatusNotStarted TaskStatus = "not_started"
	StatusActive    TaskStatus = "active"
	StatusPaused    TaskStatus = "paused"
	StatusCompleted TaskStatus = "completed"
)

type Task struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Status    TaskStatus `json:"status"`
	StartTime time.Time  `json:"startTime,omitempty"`
	EndTime   time.Time  `json:"endTime,omitempty"`
	Duration  string     `json:"duration"`
}
