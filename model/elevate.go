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

// maxElevationDuration caps a requested elevation duration. A practically
// "infinite" elevation is an intentionally supported relief valve, so this is
// deliberately large (~100 years) rather than a security limit. Its only
// purpose is to keep the value well below the point where
// time.Duration(seconds) * time.Second overflows (int64 nanoseconds wrap past
// ~292 years), which would otherwise yield an elevatedUntil in the past and
// silently drop the elevation.
const maxElevationDuration = 100 * 365 * 24 * time.Hour

// ElevationDuration converts a requested duration in seconds to a
// time.Duration, clamping absurdly large or negative values to a safe range so
// the subsequent time.Time.Add cannot overflow into a past timestamp.
func ElevationDuration(seconds int) time.Duration {
	if seconds <= 0 {
		return DefaultElevationDuration
	}
	d := time.Duration(seconds) * time.Second
	// Detect overflow (wrap to negative) or an intentionally huge value and
	// clamp to the ceiling.
	if d <= 0 || d > maxElevationDuration {
		return maxElevationDuration
	}
	return d
}
