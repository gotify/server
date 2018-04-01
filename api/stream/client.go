package stream

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/gotify/server/model"
)

const (
	writeWait = 2 * time.Second
)

type client struct {
	conn    *websocket.Conn
	onClose func(*client)
	write   chan *model.Message
	userID  uint
	token   string
	once    sync.Once
}

func newClient(conn *websocket.Conn, userID uint, token string, onClose func(*client)) *client {
	return &client{
		conn:    conn,
		write:   make(chan *model.Message, 1),
		userID:  userID,
		token:   token,
		onClose: onClose,
	}
}

// Close closes the connection.
func (c *client) Close() {
	c.once.Do(func() {
		c.conn.Close()
		close(c.write)
	})
}

// NotifyClose closes the connection and notifies that the connection was closed.
func (c *client) NotifyClose() {
	c.once.Do(func() {
		c.conn.Close()
		close(c.write)
		c.onClose(c)
	})
}

// startWriteHandler starts listening on the client connection. As we do not need anything from the client,
// we ignore incoming messages. Leaves the loop on errors.
func (c *client) startReading(pongWait time.Duration) {
	defer c.NotifyClose()
	c.conn.SetReadLimit(64)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		if _, _, err := c.conn.NextReader(); err != nil {
			return
		}
	}
}

// startWriteHandler starts the write loop. The method has the following tasks:
// * ping the client in the interval provided as parameter
// * write messages send by the channel to the client
// * on errors exit the loop
func (c *client) startWriteHandler(pingPeriod time.Duration) {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		c.NotifyClose()
		pingTicker.Stop()
	}()

	for {
		select {
		case message, ok := <-c.write:
			if !ok {
				return
			}

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteJSON(message); err != nil {
				return
			}
		case <-pingTicker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
