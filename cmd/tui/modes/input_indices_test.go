package modes

import "testing"

func TestInputIndexConstants(t *testing.T) {
	if InputProject != 0 {
		t.Fatalf("InputProject = %d, expected 0", InputProject)
	}
	if InputTitle != 1 {
		t.Fatalf("InputTitle = %d, expected 1", InputTitle)
	}
	if InputHour != 2 {
		t.Fatalf("InputHour = %d, expected 2", InputHour)
	}
	if InputMinute != 3 {
		t.Fatalf("InputMinute = %d, expected 3", InputMinute)
	}
}
