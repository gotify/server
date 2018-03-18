package stream

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/gotify/server/auth"
	"github.com/gotify/server/model"
	"net/http"
	"github.com/gotify/server/mode"
	"net/url"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		if mode.IsDev() {
			return true
		}
		return checkSameOrigin(r)
	},
}

func checkSameOrigin(r *http.Request) bool {
	origin := r.Header["Origin"]
	if len(origin) == 0 {
		return true
	}
	u, err := url.Parse(origin[0])
	if err != nil {
		return false
	}
	return u.Host == r.Host
}

// The API provides a handler for a WebSocket stream API.
type API struct {
	clients     map[uint][]*client
	lock        sync.RWMutex
	pingPeriod  time.Duration
	pongTimeout time.Duration
}

// New creates a new instance of API.
// pingPeriod: is the interval, in which is server sends the a ping to the client.
// pongTimeout: is the duration after the connection will be terminated, when the client does not respond with the
// pong command.
func New(pingPeriod, pongTimeout time.Duration) *API {
	return &API{
		clients:     make(map[uint][]*client),
		pingPeriod:  pingPeriod,
		pongTimeout: pingPeriod + pongTimeout,
	}
}

func (a *API) getClients(userID uint) ([]*client, bool) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	clients, ok := a.clients[userID]
	return clients, ok
}

// Notify notifies the clients with the given userID that a new messages was created.
func (a *API) Notify(userID uint, msg *model.Message) {
	if clients, ok := a.getClients(userID); ok {
		go func() {
			for _, c := range clients {
				c.write <- msg
			}
		}()
	}
}

func (a *API) remove(remove *client) {
	a.lock.Lock()
	defer a.lock.Unlock()
	if userIDClients, ok := a.clients[remove.userID]; ok {
		for i, client := range userIDClients {
			if client == remove {
				a.clients[remove.userID] = append(userIDClients[:i], userIDClients[i+1:]...)
				break
			}
		}
	}
}

func (a *API) register(client *client) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.clients[client.userID] = append(a.clients[client.userID], client)
}

// Handle handles incoming requests. First it upgrades the protocol to the WebSocket protocol and then starts listening
// for read and writes.
func (a *API) Handle(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}

	client := newClient(conn, auth.GetUserID(ctx), a.remove)
	a.register(client)
	go client.startReading(a.pongTimeout)
	go client.startWriteHandler(a.pingPeriod)
}

// Close closes all client connections and stops answering new connections.
func (a *API) Close() {
	a.lock.Lock()
	defer a.lock.Unlock()

	for _, clients := range a.clients {
		for _, client := range clients {
			client.Close()
		}
	}
	for k := range a.clients {
		delete(a.clients, k)
	}
}
