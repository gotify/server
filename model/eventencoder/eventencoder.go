package eventencoder

import (
	"fmt"
	"reflect"

	"github.com/gotify/server/model"
)

// TypeNotSupportedError indicates that the client does not support the specified event type.
type TypeNotSupportedError struct {
	version uint
	eventType reflect.Type
	msg string
}

// Error returns the error message associated with the given error.
func (err TypeNotSupportedError) Error() string {
	return err.msg
}

func newTypeNotSupportedError(event interface{}, version uint) error {
	return &TypeNotSupportedError {
		version,
		reflect.TypeOf(event),
		"type not supported in API version",
	}
}

type eventEncoder interface {
	Encode(event model.Event) (interface{}, error)
}

var eventEncoders = map[uint]eventEncoder{
	1: NewV1(),
	2: NewV2(),
}

// Encode the given event within the constraints of the given API version.
func Encode(apiVersion uint, event model.Event) (interface{}, error) {
	if encoder, ok := eventEncoders[apiVersion]; ok {
		return encoder.Encode(event)
	}

	err := fmt.Errorf("API version %d is not supported", apiVersion)
	return nil, err
}

// IsSupported checks if a encoder for the given API version exists.
func IsSupported(apiVersion uint) bool {
	_, ok := eventEncoders[apiVersion]
	return ok
}
