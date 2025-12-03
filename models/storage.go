package models

// This needs incremented when we change the data format
const CurrentVersion = 3

type Storage interface {
	Load() ([]TimeEntry, error)
	Save([]TimeEntry) error
}
