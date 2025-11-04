package plugin

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/gotify/plugin-api/v2/generated/protobuf"
	"github.com/gotify/server/v2/model"
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
func (c redirectToChannel) SendMessage(msg *protobuf.Message) error {
	extras := make(map[string]interface{})
	outputJson := new(bytes.Buffer)
	cnt := 0
	for k, v := range msg.Extras {
		if cnt > 0 {
			outputJson.WriteByte(',')
		}
		outputJson.WriteByte('"')
		outputJson.WriteString(k)
		outputJson.WriteByte('"')
		outputJson.WriteByte(':')
		outputJson.WriteString(v.GetJson())
		cnt++
	}

	outputJson.WriteByte('}')
	if err := json.Unmarshal(outputJson.Bytes(), &extras); err != nil {
		return err
	}

	intPriority := int(msg.Priority)

	c.Messages <- MessageWithUserID{
		Message: model.MessageExternal{
			ApplicationID: c.ApplicationID,
			Message:       msg.Message,
			Title:         msg.Title,
			Priority:      &intPriority,
			Date:          time.Now(),
			Extras:        extras,
		},
		UserID: c.UserID,
	}
	return nil
}
