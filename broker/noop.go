package broker

import (
	"github.com/gotify/server/v2/model"
)

// NoopBroker is a no-operation broker that does nothing
// Used when Redis is disabled and we fall back to local-only notifications
type NoopBroker struct{}

// NewNoopBroker creates a new no-op broker
func NewNoopBroker() *NoopBroker {
	return &NoopBroker{}
}

// PublishMessage does nothing in the no-op broker
func (n *NoopBroker) PublishMessage(userID uint, message *model.MessageExternal) error {
	// No-op: messages will be handled locally only
	return nil
}

// Subscribe does nothing in the no-op broker
func (n *NoopBroker) Subscribe(callback func(userID uint, message *model.MessageExternal)) error {
	// No-op: no external messages to subscribe to
	return nil
}

// Close does nothing in the no-op broker
func (n *NoopBroker) Close() error {
	return nil
}