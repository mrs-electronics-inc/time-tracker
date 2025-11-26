package unit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"time-tracker/utils"
)

func TestLoadMigrations(t *testing.T) {
	tempDir := t.TempDir()
	
	t.Run("Migrate V0 to Current", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "data_v0.json")
		v0Data := map[string]interface{}{
			"version": 0,
			"time-entries": []map[string]interface{}{
				{
					"id": 1,
					"start": "2023-01-01T10:00:00Z",
					"end": "2023-01-01T11:00:00Z",
					"project": "p1",
					"title": "t1",
				},
				{
					"id": 2,
					"start": "2023-01-01T12:00:00Z",
					"project": "p2",
					"title": "t2",
				},
			},
		}
		
		dataBytes, _ := json.Marshal(v0Data)
		if err := os.WriteFile(filePath, dataBytes, 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
		
		fs, err := utils.NewFileStorage(filePath)
		if err != nil {
			t.Fatalf("NewFileStorage failed: %v", err)
		}
		
		entries, err := fs.Load()
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}
		
		// Expected:
		// 1. 10:00-11:00 p1
		// 2. 11:00-12:00 blank (gap filled)
		// 3. 12:00-now p2
		
		if len(entries) != 3 {
			t.Fatalf("expected 3 entries, got %d", len(entries))
		}
		
		if entries[1].Project != "" || entries[1].Title != "" {
			t.Errorf("expected second entry to be blank, got %v %v", entries[1].Project, entries[1].Title)
		}
		
		// Verify End times are reconstructed
		if entries[0].End == nil || !entries[0].End.Equal(entries[1].Start) {
			t.Errorf("entry 0 end mismatch")
		}
		if entries[1].End == nil || !entries[1].End.Equal(entries[2].Start) {
			t.Errorf("entry 1 end mismatch")
		}
		if entries[2].End != nil {
			t.Errorf("entry 2 end should be nil")
		}
	})

	t.Run("Migrate V1 to Current (Filter Short Blanks)", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "data_v1.json")
		v1Data := map[string]interface{}{
			"version": 1,
			"time-entries": []map[string]interface{}{
				{
					"id": 1,
					"start": "2023-01-01T10:00:00Z",
					"end": "2023-01-01T11:00:00Z",
					"project": "p1",
					"title": "t1",
				},
				{
					"id": 2,
					"start": "2023-01-01T11:00:00Z",
					"end": "2023-01-01T11:00:02Z", // 2 seconds blank
					"project": "",
					"title": "",
				},
				{
					"id": 3,
					"start": "2023-01-01T12:00:00Z",
					"project": "p2",
					"title": "t2",
				},
			},
		}
		
		dataBytes, _ := json.Marshal(v1Data)
		if err := os.WriteFile(filePath, dataBytes, 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
		
		fs, err := utils.NewFileStorage(filePath)
		if err != nil {
			t.Fatalf("NewFileStorage failed: %v", err)
		}
		
		entries, err := fs.Load()
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}
		
		// Expected:
		// 1. 10:00-12:00 p1 (because the short blank is removed, and Load reconstructs End from next Start)
		// Wait. If blank is removed, then entry 1 is p1 (start 10:00), next entry is p2 (start 12:00).
		// So p1 End becomes 12:00.
		// The original end was 11:00.
		// This seems to be a consequence of the system: "End times for all entries based on the start time of the next entry".
		// So yes, removing the blank entry EXTENDS the previous entry to cover the gap.
		// Unless there was another mechanism?
		// But `Load` overwrites `End`.
		
		if len(entries) != 2 {
			t.Fatalf("expected 2 entries, got %d", len(entries))
		}
		
		if entries[0].Project != "p1" {
			t.Errorf("expected first entry p1")
		}
		if entries[1].Project != "p2" {
			t.Errorf("expected second entry p2")
		}
		
		// Check that p1 end is p2 start (12:00)
		expectedEnd, _ := time.Parse(time.RFC3339, "2023-01-01T12:00:00Z")
		if entries[0].End == nil || !entries[0].End.Equal(expectedEnd) {
			t.Errorf("expected p1 end to be %v, got %v", expectedEnd, entries[0].End)
		}
	})

	t.Run("Migrate V2 to Current (No End Field)", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "data_v2.json")
		v2Data := map[string]interface{}{
			"version": 2,
			"time-entries": []map[string]interface{}{
				{
					"id": 1,
					"start": "2023-01-01T10:00:00Z",
					"project": "p1",
					"title": "t1",
				},
				{
					"id": 2,
					"start": "2023-01-01T11:00:00Z",
					"project": "p2",
					"title": "t2",
				},
			},
		}
		
		dataBytes, _ := json.Marshal(v2Data)
		if err := os.WriteFile(filePath, dataBytes, 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
		
		fs, err := utils.NewFileStorage(filePath)
		if err != nil {
			t.Fatalf("NewFileStorage failed: %v", err)
		}
		
		entries, err := fs.Load()
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}
		
		if len(entries) != 2 {
			t.Fatalf("expected 2 entries, got %d", len(entries))
		}
		
		// Verify End times are reconstructed
		if entries[0].End == nil || !entries[0].End.Equal(entries[1].Start) {
			t.Errorf("entry 0 end mismatch")
		}
		if entries[1].End != nil {
			t.Errorf("entry 1 end should be nil")
		}
	})
}
