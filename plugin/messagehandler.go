package plugin

import (
	"errors"
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
	if c.ApplicationID == 0 {
		// Final safety net: the internal application should always be set up by
		// Manager.initializeSingleUserPlugin. If it somehow isn't, refuse the
		// message instead of storing it with application_id = 0, where it would
		// be orphaned (not shown in the UI and not deletable).
		return errors.New("plugin messenger has no associated internal application")
	}
	c.Messages <- MessageWithUserID{
		Message: model.MessageExternal{
			ApplicationID: c.ApplicationID,
			Message:       msg.Message,
			Title:         msg.Title,
			Priority:      &msg.Priority,
			Date:          time.Now(),
			Extras:        msg.Extras,
		},
		UserID: c.UserID,
	}
	return nil
}
