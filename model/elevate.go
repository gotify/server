package model

import "time"

// ElevateRequest parameters for client elevation.
//
// swagger:model ElevateRequest
type ElevateRequest struct {
	// How long the elevation should last, in seconds.
	//
	// required: true
	// example: 900
	DurationSeconds int `form:"durationSeconds" query:"durationSeconds" json:"durationSeconds" binding:"required"`
}

var DefaultElevationDuration = time.Hour

// maxElevationSeconds caps a requested elevation duration (~100 years). A
// practically "infinite" elevation is intentionally supported, so this is not a
// security limit; its only purpose is to keep the value well below the point
// where time.Duration(seconds) * time.Second overflows int64 nanoseconds
// (~292 years) and would otherwise wrap to a bogus (even past) timestamp.
const maxElevationSeconds = 100 * 365 * 24 * 60 * 60

// ElevationDuration converts a requested duration in seconds to a
// time.Duration, clamping the seconds (before multiplying) to a safe range so
// the subsequent multiply cannot overflow. Negative values are preserved so the
// client can still cancel an elevation (the UI sends durationSeconds: -1, which
// must set elevatedUntil in the past rather than grant a fresh elevation).
func ElevationDuration(seconds int) time.Duration {
	if seconds > maxElevationSeconds {
		seconds = maxElevationSeconds
	} else if seconds < -maxElevationSeconds {
		seconds = -maxElevationSeconds
	}
	return time.Duration(seconds) * time.Second
}
