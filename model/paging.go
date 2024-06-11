package model

// Paging Model
//
// The Paging holds information about the limit and making requests to the next page.
//
// swagger:model Paging
type Paging struct {
	// The request url for the next page. Empty/Null when no next page is available.
	//
	// read only: true
	// required: false
	// example: http://example.com/message?limit=50&since=123456
	Next string `json:"next,omitempty"`
	// The amount of messages that got returned in the current request.
	//
	// read only: true
	// required: true
	// example: 5
	Size int `json:"size"`
	// The ID of the last message returned in the current request. Use this as alternative to the next link.
	//
	// read only: true
	// required: true
	// example: 5
	// min: 0
	Since uint `json:"since"`
	// The limit of the messages for the current request.
	//
	// read only: true
	// required: true
	// min: 1
	// max: 200
	// example: 123
	Limit int `json:"limit"`
}

// PagedMessages Model
//
// Wrapper for the paging and the messages.
//
// swagger:model PagedMessages
type PagedMessages struct {
	// The paging of the messages.
	//
	// read only: true
	// required: true
	Paging Paging `json:"paging"`
	// The messages.
	//
	// read only: true
	// required: true
	Messages []*MessageExternal `json:"messages"`
}
