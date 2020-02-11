package eventencoder

import (
	"github.com/gotify/server/model"
)

// V2 encodes given events within the constraints of API version 1.
type V2 struct {
	Version uint
}

// NewV2 creates a new instance of V2.
func NewV2() *V2 {
	return &V2{
		2,
	}
}

// Encode encodes the given event.
func (encoder V2) Encode(event model.Event) (interface{}, error) {
	var encoded eventMessageV2

	switch event := event.(type) {
	case *model.Message:
		encoded.Type = "message"
		encoded.Content = event.ToExternal()
	case *model.MessageDeletions:
		encoded.Type = "message_deletions"
		encoded.Content = event.ToExternal()
	default:
		err := newTypeNotSupportedError(event, encoder.Version)
		return nil, err
	}

	return encoded, nil
}

// eventMessageV2 Model
//
// eventMessageV2 holds information about a message which will be sent to the clients.
//
// swagger:model eventMessageV2
type eventMessageV2 struct {
	// The event type.
	//
	// read only: true
	// required: true
	// example: message_deletion
	Type string `json:"type"`
	// The event content.
	//
	// read only: true
	// required: true
	// example: [14,15,16]
	Content interface{} `json:"content"`
}
