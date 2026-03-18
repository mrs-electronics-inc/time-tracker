package modes

import "testing"

func TestInputIndexConstants(t *testing.T) {
	if InputProject != 0 {
		t.Fatalf("InputProject = %d, expected 0", InputProject)
	}
	if InputTitle != 1 {
		t.Fatalf("InputTitle = %d, expected 1", InputTitle)
	}
	if InputYear != 2 {
		t.Fatalf("InputYear = %d, expected 2", InputYear)
	}
	if InputMonth != 3 {
		t.Fatalf("InputMonth = %d, expected 3", InputMonth)
	}
	if InputDay != 4 {
		t.Fatalf("InputDay = %d, expected 4", InputDay)
	}
	if InputHour != 5 {
		t.Fatalf("InputHour = %d, expected 5", InputHour)
	}
	if InputMinute != 6 {
		t.Fatalf("InputMinute = %d, expected 6", InputMinute)
	}
}
