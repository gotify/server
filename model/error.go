package model

// Error Model
//
// The Error contains error relevant information.
//
// swagger:model Error
type Error struct {
	// The general error message
	//
	// required: true
	// example: Unauthorized
	Error string `json:"error"`
	// The http error code.
	//
	// required: true
	// example: 401
	ErrorCode int `json:"errorCode"`
	// The http error code.
	//
	// required: true
	// example: you need to provide a valid access token or user credentials to access this api
	ErrorDescription string `json:"errorDescription"`
}
