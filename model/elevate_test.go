package model

import (
	"testing"
	"time"
)

func TestElevationDuration(t *testing.T) {
	cases := []struct {
		name    string
		seconds int
		want    time.Duration
	}{
		{"normal", 900, 900 * time.Second},
		{"zero falls back to default", 0, DefaultElevationDuration},
		{"negative falls back to default", -5, DefaultElevationDuration},
		{"huge value is clamped", 9223372036854775807, maxElevationDuration},
		{"just under ceiling is unchanged", 100 * 365 * 24 * 3600, maxElevationDuration},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := ElevationDuration(c.seconds)
			if got != c.want {
				t.Fatalf("ElevationDuration(%d) = %v, want %v", c.seconds, got, c.want)
			}
			// The resulting time must always be in the future (no overflow).
			if time.Now().Add(got).Before(time.Now()) {
				t.Fatalf("ElevationDuration(%d) produced a past time", c.seconds)
			}
		})
	}
}
