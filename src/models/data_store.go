package models

// This needs incremented when we change the data format
CurrentVersion = 0

// DataStore represents the JSON structure for time entries
type DataStore struct {
	Version     int                `json:"version"`
	TimeEntries []models.TimeEntry `json:"time-entries"`
}
