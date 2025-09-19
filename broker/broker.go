package broker

import (
	"github.com/gotify/server/v2/model"
)

// MessageBroker defines the interface for distributing messages across multiple instances
type MessageBroker interface {
	// PublishMessage publishes a message to all subscribers
	PublishMessage(userID uint, message *model.MessageExternal) error
	// Subscribe starts listening for messages and calls the callback for each received message
	Subscribe(callback func(userID uint, message *model.MessageExternal)) error
	// Close closes the broker connection
	Close() error
}