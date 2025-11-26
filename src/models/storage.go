package models

// This needs incremented when we change the data format
const CurrentVersion = 2

type Storage interface {
	Load() ([]TimeEntry, error)
	Save([]TimeEntry) error
}
