package plugin

import (
	"time"

	"github.com/gotify/server/v2/model"
	"github.com/gotify/server/v2/plugin/compat"
)

type redirectToChannel struct {
	ApplicationID uint
	UserID        uint
	Messages      chan MessageWithUserID
}

// MessageWithUserID encapsulates a message with a given user ID.
type MessageWithUserID struct {
	Message model.MessageExternal
	UserID  uint
}

// SendMessage sends a message to the underlying message channel.
func (c redirectToChannel) SendMessage(msg compat.Message) error {
	c.Messages <- MessageWithUserID{
		Message: model.MessageExternal{
			ApplicationID: c.ApplicationID,
			Message:       msg.Message,
			Title:         msg.Title,
			Priority:      msg.Priority,
			Date:          time.Now(),
			Extras:        msg.Extras,
		},
		UserID: c.UserID,
	}
	return nil
}
