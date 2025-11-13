package models

type Storage interface {
	Load() ([]TimeEntry, error)
	Save([]TimeEntry) error
}

