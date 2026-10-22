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

// maxElevation caps a requested elevation duration (~100 years). A
// practically "infinite" elevation is intentionally supported, so this is not a
// security limit; its only purpose is to keep the value well below the point
// where time.Duration(seconds) * time.Second overflows int64 nanoseconds
// (~292 years) and would otherwise wrap to a bogus (even past) timestamp.
const maxElevation = 100 * 365 * 24 * time.Hour

// ElevationTime calculates the time at which an elevation will expire.
func ElevationTime(base time.Time, deltaSeconds int) time.Time {
	if deltaSeconds > int(maxElevation/time.Second) {
		deltaSeconds = int(maxElevation / time.Second)
	}
	if deltaSeconds < 0 {
		return time.Time{}
	}
	return base.Add(time.Duration(deltaSeconds) * time.Second)
}
