package model

// Error Model
//
// The Error contains error relevant information.
//
// swagger:model Error
type Error struct {
	Error            string `json:"error"`
	ErrorCode        int    `json:"errorCode"`
	ErrorDescription string `json:"errorDescription"`
}
