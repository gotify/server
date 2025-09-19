package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gotify/server/v2/model"
	"github.com/redis/go-redis/v9"
)

const (
	defaultChannel         = "gotify:messages"
	defaultReconnectDelay  = 5 * time.Second
	defaultPingInterval    = 30 * time.Second
)

// RedisMessage represents the message format sent through Redis
type RedisMessage struct {
	UserID  uint                     `json:"user_id"`
	Message *model.MessageExternal   `json:"message"`
}

// RedisBroker implements MessageBroker using Redis pub/sub
type RedisBroker struct {
	client    *redis.Client
	pubsub    *redis.PubSub
	channel   string
	ctx       context.Context
	cancel    context.CancelFunc
	closed    bool
}

// NewRedisBroker creates a new Redis broker instance
func NewRedisBroker(redisURL, channelPrefix string) (*RedisBroker, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("invalid Redis URL: %v", err)
	}

	client := redis.NewClient(opts)
	
	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	ctx, cancel = context.WithCancel(context.Background())
	
	channel := defaultChannel
	if channelPrefix != "" {
		channel = channelPrefix + ":messages"
	}

	return &RedisBroker{
		client:  client,
		channel: channel,
		ctx:     ctx,
		cancel:  cancel,
		closed:  false,
	}, nil
}

// PublishMessage publishes a message to Redis for distribution to all subscribers
func (r *RedisBroker) PublishMessage(userID uint, message *model.MessageExternal) error {
	if r.closed {
		return fmt.Errorf("broker is closed")
	}

	redisMsg := RedisMessage{
		UserID:  userID,
		Message: message,
	}

	data, err := json.Marshal(redisMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	return r.client.Publish(ctx, r.channel, data).Err()
}

// Subscribe starts listening for messages on the Redis channel
func (r *RedisBroker) Subscribe(callback func(userID uint, message *model.MessageExternal)) error {
	if r.closed {
		return fmt.Errorf("broker is closed")
	}

	r.pubsub = r.client.Subscribe(r.ctx, r.channel)

	// Wait for confirmation that subscription is created
	_, err := r.pubsub.Receive(r.ctx)
	if err != nil {
		return fmt.Errorf("failed to subscribe to Redis channel: %v", err)
	}

	// Start processing messages in a goroutine
	go r.processMessages(callback)

	return nil
}

// processMessages handles incoming Redis messages
func (r *RedisBroker) processMessages(callback func(userID uint, message *model.MessageExternal)) {
	ch := r.pubsub.Channel()

	for {
		select {
		case <-r.ctx.Done():
			return
		case msg := <-ch:
			if msg == nil {
				continue
			}

			var redisMsg RedisMessage
			if err := json.Unmarshal([]byte(msg.Payload), &redisMsg); err != nil {
				log.Printf("Failed to unmarshal Redis message: %v", err)
				continue
			}

			// Call the callback with the parsed message
			callback(redisMsg.UserID, redisMsg.Message)
		}
	}
}

// Close closes the Redis connection and stops the subscriber
func (r *RedisBroker) Close() error {
	if r.closed {
		return nil
	}

	r.closed = true
	r.cancel()

	var err error
	if r.pubsub != nil {
		err = r.pubsub.Close()
	}

	if closeErr := r.client.Close(); closeErr != nil && err == nil {
		err = closeErr
	}

	return err
}