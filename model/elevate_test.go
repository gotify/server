package model

import (
	"math"
	"testing"
	"time"
)

func TestElevationTime(t *testing.T) {
	// Use runtime variables so arithmetic doesn't happen as a constant
	// expression and overflow Go's constant checker on either arch.
	maxInt := math.MaxInt
	minInt := math.MinInt
	maxElevationSeconds := int64(maxElevation / time.Second)

	// On 32-bit math.MaxInt (~2.1B) is below maxElevationSeconds (~3.15B),
	// so no clamp triggers and the expected output is the requested
	// duration; on 64-bit math.MaxInt is far above, so the expected output
	// is the cap.
	var expectedForHuge time.Duration
	if int64(maxInt) < maxElevationSeconds {
		expectedForHuge = time.Duration(maxInt) * time.Second
	} else {
		expectedForHuge = maxElevation
	}

	baseTime := time.Now()

	cases := []struct {
		name    string
		seconds int
		want    time.Time
	}{
		{"normal", 900, baseTime.Add(900 * time.Second)},
		{"zero stays zero", 0, baseTime},
		// Negative values collapse to a one-second past duration so the
		// client can cancel an elevation (UI sends -1); they must NOT grant
		// a fresh positive elevation.
		{"negative value cancels", -1, baseTime.Add(-time.Second)},
		{"large negative cancels", minInt, baseTime.Add(-time.Second)},
		// A huge positive request is clamped on 64-bit; on 32-bit no clamp
		// is reachable since math.MaxInt < maxElevationSeconds there.
		{"huge clamped to max", maxInt, baseTime.Add(expectedForHuge)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := ElevationTime(baseTime, c.seconds); got != c.want {
				t.Fatalf("ElevationTime(%d) = %v, want %v", c.seconds, got, c.want)
			}
		})
	}
}
