package eventencoder

import (
	"github.com/gotify/server/model"
)

// V1 encodes given events within the constraints of API version 1.
type V1 struct {
	Version uint
}

// NewV1 creates a new instance of V1.
func NewV1() *V1 {
	return &V1{
		1,
	}
}

// Encode encodes the given event.
func (encoder V1) Encode(event model.Event) (interface{}, error) {
	var encoded eventMessageV1

	switch event := event.(type) {
	case *model.Message:
		encoded = event.ToExternal().(eventMessageV1)
	default:
		err := newTypeNotSupportedError(event, encoder.Version)
		return nil, err
	}

	return encoded, nil
}

// eventMessageV1 Model
//
// eventMessageV1 holds information about a message which will be sent to the clients.
//
// swagger:model eventMessageV1
type eventMessageV1 = *model.MessageExternal
