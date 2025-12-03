package models

import "time"

// V0Entry is the original pre-migration format.
type V0Entry struct {
	ID      int        `json:"id"`
	Start   time.Time  `json:"start"`
	End     *time.Time `json:"end,omitempty"`
	Project string     `json:"project"`
	Title   string     `json:"title"`
}

// V1Entry is the format after V0->V1 migration (sorted and gap-filled).
type V1Entry struct {
	ID      int        `json:"id"`
	Start   time.Time  `json:"start"`
	End     *time.Time `json:"end,omitempty"`
	Project string     `json:"project"`
	Title   string     `json:"title"`
}

// V2Entry is the format after V1->V2 migration (short blanks filtered, end field removed).
type V2Entry struct {
	ID      int       `json:"id"`
	Start   time.Time `json:"start"`
	Project string    `json:"project"`
	Title   string    `json:"title"`
}

// V3Entry is the format after V2->V3 migration (ID field removed).
type V3Entry struct {
	Start   time.Time `json:"start"`
	Project string    `json:"project"`
	Title   string    `json:"title"`
}
