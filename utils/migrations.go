package utils

import (
	"sort"
	"time"

	"time-tracker/models"
)

func TransformV0ToV1(entries []models.V0Entry) ([]models.V1Entry, error) {
	if len(entries) == 0 {
		return nil, nil
	}

	// Make a shallow copy to avoid mutating the input
	copied := append([]models.V0Entry(nil), entries...)

	// Sort by start time
	sort.Slice(copied, func(i, j int) bool {
		return copied[i].Start.Before(copied[j].Start)
	})

	// Find max ID
	maxID := 0
	for _, e := range copied {
		if e.ID > maxID {
			maxID = e.ID
		}
	}

	var newEntries []models.V1Entry
	for i, entry := range copied {
		newEntries = append(newEntries, models.V1Entry{
			ID:      entry.ID,
			Start:   entry.Start,
			End:     entry.End,
			Project: entry.Project,
			Title:   entry.Title,
		})
		if i < len(copied)-1 && entry.End != nil && entry.End.Before(copied[i+1].Start) {
			end := copied[i+1].Start
			maxID++
			blank := models.V1Entry{
				ID:      maxID,
				Start:   *entry.End,
				End:     &end,
				Project: "",
				Title:   "",
			}
			newEntries = append(newEntries, blank)
		}
	}
	return newEntries, nil
}

func TransformV1ToV2(entries []models.V1Entry) ([]models.V2Entry, error) {
	// Filter out blank entries that are less than 5 seconds long
	var filtered []models.V2Entry
	for _, entry := range entries {
		// Skip if it's a blank entry with duration < 5 seconds
		if entry.Project == "" && entry.Title == "" && entry.End != nil {
			duration := entry.End.Sub(entry.Start)
			if duration < 5*time.Second {
				continue
			}
		}
		filtered = append(filtered, models.V2Entry{
			ID:      entry.ID,
			Start:   entry.Start,
			Project: entry.Project,
			Title:   entry.Title,
		})
	}
	return filtered, nil
}

func TransformV2ToV3(entries []models.V2Entry) ([]models.V3Entry, error) {
	v3Entries := make([]models.V3Entry, len(entries))
	for i, entry := range entries {
		v3Entries[i] = models.V3Entry{
			Start:   entry.Start,
			Project: entry.Project,
			Title:   entry.Title,
		}
	}
	return v3Entries, nil
}
