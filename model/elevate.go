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

const maxElevation = 100 * 365 * 24 * time.Hour

// ElevationTime calculates the time at which an elevation will expire.
func ElevationTime(base time.Time, deltaSeconds int) time.Time {
	if deltaSeconds < 0 {
		return time.Time{}
	}

	// Prevent int64 overflow by limiting to an arbitrary high number.
	seconds := min(time.Duration(deltaSeconds), maxElevation/time.Second)
	return base.Add(seconds * time.Second)
}
