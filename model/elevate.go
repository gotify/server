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
	DurationSeconds int `form:"durationSeconds" query:"durationSeconds" json:"durationSeconds" binding:"required,min=1,max=2592000"`
}

var DefaultElevationDuration = time.Hour

// MaxElevationDurationSeconds is the upper bound (30 days) accepted for a
// requested elevation duration. It bounds ElevateRequest.DurationSeconds so a
// request can neither make elevation effectively permanent nor overflow
// time.Duration (which wraps past ~292 years, silently yielding a past
// elevatedUntil timestamp). Requests above this are rejected with 400.
const MaxElevationDurationSeconds = 2592000
