package utils

import (
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "minutes only",
			duration: 30 * time.Minute,
			expected: "30m",
		},
		{
			name:     "one minute",
			duration: 1 * time.Minute,
			expected: "1m",
		},
		{
			name:     "hours only",
			duration: 2 * time.Hour,
			expected: "2h",
		},
		{
			name:     "one hour",
			duration: 1 * time.Hour,
			expected: "1h",
		},
		{
			name:     "hours and minutes",
			duration: 1*time.Hour + 30*time.Minute,
			expected: "1h 30m",
		},
		{
			name:     "multiple hours and minutes",
			duration: 2*time.Hour + 45*time.Minute,
			expected: "2h 45m",
		},
		{
			name:     "zero duration",
			duration: 0,
			expected: "0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDuration(tt.duration)
			if result != tt.expected {
				t.Errorf("FormatDuration(%v) = %q, want %q", tt.duration, result, tt.expected)
			}
		})
	}
}
