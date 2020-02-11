package model

// Event is an interface for all event types.
type Event interface {
	ToExternal() interface{}
}
