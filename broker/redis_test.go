package broker

import (
	"testing"
	"time"

	"github.com/gotify/server/v2/model"
	"github.com/stretchr/testify/assert"
)

func TestNoopBroker(t *testing.T) {
	broker := NewNoopBroker()
	defer broker.Close()

	priority := 1
	// Test that operations don't fail
	err := broker.PublishMessage(1, &model.MessageExternal{
		ID:            1,
		ApplicationID: 1,
		Message:       "test",
		Title:         "Test",
		Priority:      &priority,
		Date:          time.Now(),
	})
	assert.NoError(t, err)

	err = broker.Subscribe(func(userID uint, message *model.MessageExternal) {})
	assert.NoError(t, err)

	err = broker.Close()
	assert.NoError(t, err)
}

func TestRedisBroker_InvalidURL(t *testing.T) {
	_, err := NewRedisBroker("invalid-url", "test")
	assert.Error(t, err)
}

func TestRedisBroker_ConnectionFailure(t *testing.T) {
	// This will fail to connect since there's no Redis server
	_, err := NewRedisBroker("redis://localhost:9999/0", "test")
	assert.Error(t, err)
}