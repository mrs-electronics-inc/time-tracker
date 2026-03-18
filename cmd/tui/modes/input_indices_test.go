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
	if InputYear != 4 {
		t.Fatalf("InputYear = %d, expected 4", InputYear)
	}
	if InputMonth != 5 {
		t.Fatalf("InputMonth = %d, expected 5", InputMonth)
	}
	if InputDay != 6 {
		t.Fatalf("InputDay = %d, expected 6", InputDay)
	}
}
