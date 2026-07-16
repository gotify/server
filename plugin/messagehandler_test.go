package plugin

import (
	"testing"

	"github.com/gotify/server/v2/plugin/compat"
	"github.com/stretchr/testify/assert"
)

func TestRedirectToChannel_SendMessage_rejectsMissingApplication(t *testing.T) {
	messages := make(chan MessageWithUserID, 1)
	handler := redirectToChannel{ApplicationID: 0, UserID: 1, Messages: messages}

	err := handler.SendMessage(compat.Message{Message: "orphan"})

	assert.Error(t, err)
	assert.Empty(t, messages, "no message should be queued when the internal application is missing")
}

func TestRedirectToChannel_SendMessage_forwardsWithApplication(t *testing.T) {
	messages := make(chan MessageWithUserID, 1)
	handler := redirectToChannel{ApplicationID: 7, UserID: 3, Messages: messages}

	assert.NoError(t, handler.SendMessage(compat.Message{Message: "hi", Title: "t"}))

	got := <-messages
	assert.Equal(t, uint(7), got.Message.ApplicationID)
	assert.Equal(t, uint(3), got.UserID)
	assert.Equal(t, "hi", got.Message.Message)
	assert.Equal(t, "t", got.Message.Title)
}
