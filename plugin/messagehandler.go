package plugin

import (
	"github.com/gotify/server/model"
	"github.com/gotify/server/plugin/compat"
)

type redirectToChannel struct {
	ApplicationID uint
	UserID        uint
	Messages      chan MessageWithUserID
}

// MessageWithUserID encapsulates a message with a given user ID
type MessageWithUserID struct {
	Message *model.Message
	UserID  uint
}

// SendMessage sends a message to the underlying message channel
func (c redirectToChannel) SendMessage(msg compat.Message) error {
	c.Messages <- MessageWithUserID{
		Message: msg.ToInternalMessage(c.ApplicationID),
		UserID: c.UserID,
	}
	return nil
}
