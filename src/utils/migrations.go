package utils

import (
	"encoding/json"
	"reflect"
	"sort"
	"time"

	"time-tracker/models"
)

// Transformations maps source version to transformation function
var Transformations = map[int]any{
	0: TransformV0ToV1,
	1: TransformV1ToV2,
	2: TransformV2ToV3,
}

// callTransformWithMarshal handles the common pattern of unmarshal -> transform -> marshal
// It uses reflection to work with any transformation function signature
func callTransformWithMarshal(data []byte, transformFunc any) ([]byte, error) {
	// First, try to unmarshal into a generic any
	var entries any
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}

	// Call the transform function via reflection
	rf := reflect.ValueOf(transformFunc)
	result := rf.Call([]reflect.Value{reflect.ValueOf(entries)})

	// Check for error return
	if !result[1].IsNil() {
		return nil, result[1].Interface().(error)
	}

	transformed := result[0].Interface()

	// Marshal the result
	marshalledData, err := json.Marshal(transformed)
	if err != nil {
		return nil, err
	}
	return marshalledData, nil
}

// MigrateWithTransform is a generic helper for transformation functions
func MigrateWithTransform[In, Out any](data []byte, transform func([]In) ([]Out, error)) ([]byte, error) {
	var entries []In
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}

	transformed, err := transform(entries)
	if err != nil {
		return nil, err
	}

	result, err := json.Marshal(transformed)
	if err != nil {
		return nil, err
	}
	return result, nil
}

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

func TransformV1ToV2(entries []models.V1Entry) ([]models.V1Entry, error) {
	// Filter out blank entries that are less than 5 seconds long
	var filtered []models.V1Entry
	for _, entry := range entries {
		// Skip if it's a blank entry with duration < 5 seconds
		if entry.Project == "" && entry.Title == "" && entry.End != nil {
			duration := entry.End.Sub(entry.Start)
			if duration < 5*time.Second {
				continue
			}
		}
		filtered = append(filtered, entry)
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
