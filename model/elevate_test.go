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
		{"zero stays zero", 0, 0},
		// Negative values must be preserved so the client can cancel an
		// elevation (UI sends -1); they must NOT become a positive duration.
		{"cancel -1 preserved", -1, -1 * time.Second},
		{"huge clamped to max", 9223372036854775807, maxElevationSeconds * time.Second},
		// Values that would overflow when multiplied but land positive are
		// clamped by clamping seconds first.
		{"overflow-positive clamped", 20023544073, maxElevationSeconds * time.Second},
		{"large negative clamped", -9223372037, -maxElevationSeconds * time.Second},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := ElevationDuration(c.seconds); got != c.want {
				t.Fatalf("ElevationDuration(%d) = %v, want %v", c.seconds, got, c.want)
			}
		})
	}
}
