package cmd

import (
	"strings"
	"testing"
)

func TestValidateStatsFlagsRejectsWeeklyAndMonthly(t *testing.T) {
	err := validateStatsFlags(true, true)
	if err == nil {
		t.Fatal("expected error when weekly and monthly flags are both set")
	}

	if !strings.Contains(err.Error(), "cannot be used together") {
		t.Fatalf("expected conflicting flags error, got: %v", err)
	}
}

func TestValidateStatsFlagsAcceptsSingleMode(t *testing.T) {
	tests := []struct {
		name        string
		weeklyFlag  bool
		monthlyFlag bool
	}{
		{name: "daily default"},
		{name: "weekly", weeklyFlag: true},
		{name: "monthly", monthlyFlag: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateStatsFlags(tt.weeklyFlag, tt.monthlyFlag); err != nil {
				t.Fatalf("expected flags to be accepted, got: %v", err)
			}
		})
	}
}
