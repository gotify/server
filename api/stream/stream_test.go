package stream

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/gotify/server/auth"
	"github.com/gotify/server/mode"
	"github.com/gotify/server/model"
	"github.com/stretchr/testify/assert"
)

func TestFailureOnNormalHttpRequest(t *testing.T) {
	mode.Set(mode.TestDev)

	defer leaktest.Check(t)()

	server, api := bootTestServer(staticUserID())
	defer server.Close()
	defer api.Close()

	resp, err := http.Get(server.URL)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	resp.Body.Close()
}

func TestWriteMessageFails(t *testing.T) {
	mode.Set(mode.TestDev)
	oldWrite := writeJSON
	// try emulate an write error, mostly this should kill the ReadMessage goroutine first but you'll never know.
	writeJSON = func(conn *websocket.Conn, v interface{}) error {
		return errors.New("asd")
	}
	defer func() {
		writeJSON = oldWrite
	}()
	defer leaktest.Check(t)()

	server, api := bootTestServer(func(context *gin.Context) {
		auth.RegisterAuthentication(context, nil, 1, "")
	})
	defer server.Close()
	defer api.Close()

	wsURL := wsURL(server.URL)
	user := testClient(t, wsURL)

	// the server may take some time to register the client
	time.Sleep(100 * time.Millisecond)
	clients := clients(api, 1)
	assert.NotEmpty(t, clients)

	api.Notify(1, &model.MessageExternal{Message: "HI"})
	user.expectNoMessage()
}

func TestWritePingFails(t *testing.T) {
	mode.Set(mode.TestDev)
	oldPing := ping
	// try emulate an write error, mostly this should kill the ReadMessage gorouting first but you'll never know.
	ping = func(conn *websocket.Conn) error {
		return errors.New("asd")
	}
	defer func() {
		ping = oldPing
	}()

	defer leaktest.CheckTimeout(t, 10*time.Second)()

	server, api := bootTestServer(staticUserID())
	defer api.Close()
	defer server.Close()

	wsURL := wsURL(server.URL)
	user := testClient(t, wsURL)
	defer user.conn.Close()

	// the server may take some time to register the client
	time.Sleep(100 * time.Millisecond)
	clients := clients(api, 1)

	assert.NotEmpty(t, clients)

	time.Sleep(api.pingPeriod) // waiting for ping

	api.Notify(1, &model.MessageExternal{Message: "HI"})
	user.expectNoMessage()
}

func TestPing(t *testing.T) {
	mode.Set(mode.TestDev)

	server, api := bootTestServer(staticUserID())
	defer server.Close()
	defer api.Close()

	wsURL := wsURL(server.URL)

	user := createClient(t, wsURL)
	defer user.conn.Close()

	ping := make(chan bool)
	oldPingHandler := user.conn.PingHandler()
	user.conn.SetPingHandler(func(appData string) error {
		err := oldPingHandler(appData)
		ping <- true
		return err
	})

	startReading(user)

	expectNoMessage(user)

	select {
	case <-time.After(2 * time.Second):
		assert.Fail(t, "Expected ping but there was one :(")
	case <-ping:
		// expected
	}

	expectNoMessage(user)
	api.Notify(1, &model.MessageExternal{Message: "HI"})
	user.expectMessage(&model.MessageExternal{Message: "HI"})
}

func TestCloseClientOnNotReading(t *testing.T) {
	mode.Set(mode.TestDev)

	server, api := bootTestServer(staticUserID())
	defer server.Close()
	defer api.Close()

	wsURL := wsURL(server.URL)

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.Nil(t, err)
	defer ws.Close()

	// the server may take some time to register the client
	time.Sleep(100 * time.Millisecond)
	assert.NotEmpty(t, clients(api, 1))

	time.Sleep(api.pingPeriod + api.pongTimeout)

	assert.Empty(t, clients(api, 1))
}

func TestMessageDirectlyAfterConnect(t *testing.T) {
	mode.Set(mode.Prod)
	defer leaktest.Check(t)()
	server, api := bootTestServer(staticUserID())
	defer server.Close()
	defer api.Close()

	wsURL := wsURL(server.URL)

	user := testClient(t, wsURL)
	defer user.conn.Close()
	// the server may take some time to register the client
	time.Sleep(100 * time.Millisecond)
	api.Notify(1, &model.MessageExternal{Message: "msg"})
	user.expectMessage(&model.MessageExternal{Message: "msg"})
}

func TestDeleteClientShouldCloseConnection(t *testing.T) {
	mode.Set(mode.Prod)
	defer leaktest.Check(t)()
	server, api := bootTestServer(staticUserID())
	defer server.Close()
	defer api.Close()

	wsURL := wsURL(server.URL)

	user := testClient(t, wsURL)
	defer user.conn.Close()
	// the server may take some time to register the client
	time.Sleep(100 * time.Millisecond)
	api.Notify(1, &model.MessageExternal{Message: "msg"})
	user.expectMessage(&model.MessageExternal{Message: "msg"})

	api.NotifyDeletedClient(1, "customtoken")

	api.Notify(1, &model.MessageExternal{Message: "msg"})
	user.expectNoMessage()
}

func TestDeleteMultipleClients(t *testing.T) {
	mode.Set(mode.TestDev)

	defer leaktest.Check(t)()
	userIDs := []uint{1, 1, 1, 1, 2, 2, 3}
	tokens := []string{"1-1", "1-2", "1-2", "1-3", "2-1", "2-2", "3"}
	i := 0
	server, api := bootTestServer(func(context *gin.Context) {
		auth.RegisterAuthentication(context, nil, userIDs[i], tokens[i])
		i++
	})
	defer server.Close()

	wsURL := wsURL(server.URL)

	userOneIPhone := testClient(t, wsURL)
	defer userOneIPhone.conn.Close()
	userOneAndroid := testClient(t, wsURL)
	defer userOneAndroid.conn.Close()
	userOneBrowser := testClient(t, wsURL)
	defer userOneBrowser.conn.Close()
	userOneOther := testClient(t, wsURL)
	defer userOneOther.conn.Close()
	userOne := []*testingClient{userOneAndroid, userOneBrowser, userOneIPhone, userOneOther}

	userTwoBrowser := testClient(t, wsURL)
	defer userTwoBrowser.conn.Close()
	userTwoAndroid := testClient(t, wsURL)
	defer userTwoAndroid.conn.Close()
	userTwo := []*testingClient{userTwoAndroid, userTwoBrowser}

	userThreeAndroid := testClient(t, wsURL)
	defer userThreeAndroid.conn.Close()
	userThree := []*testingClient{userThreeAndroid}

	// the server may take some time to register the client
	time.Sleep(100 * time.Millisecond)

	api.Notify(1, &model.MessageExternal{ID: 4, Message: "there"})
	expectMessage(&model.MessageExternal{ID: 4, Message: "there"}, userOne...)
	expectNoMessage(userTwo...)
	expectNoMessage(userThree...)

	api.NotifyDeletedClient(1, "1-2")

	api.Notify(1, &model.MessageExternal{ID: 2, Message: "there"})
	expectMessage(&model.MessageExternal{ID: 2, Message: "there"}, userOneIPhone, userOneOther)
	expectNoMessage(userOneBrowser, userOneAndroid)
	expectNoMessage(userThree...)
	expectNoMessage(userTwo...)

	api.Notify(2, &model.MessageExternal{ID: 2, Message: "there"})
	expectNoMessage(userOne...)
	expectMessage(&model.MessageExternal{ID: 2, Message: "there"}, userTwo...)
	expectNoMessage(userThree...)

	api.Notify(3, &model.MessageExternal{ID: 5, Message: "there"})
	expectNoMessage(userOne...)
	expectNoMessage(userTwo...)
	expectMessage(&model.MessageExternal{ID: 5, Message: "there"}, userThree...)

	api.Close()
}

func TestDeleteUser(t *testing.T) {
	mode.Set(mode.TestDev)

	defer leaktest.Check(t)()
	userIDs := []uint{1, 1, 1, 1, 2, 2, 3}
	tokens := []string{"1-1", "1-2", "1-2", "1-3", "2-1", "2-2", "3"}
	i := 0
	server, api := bootTestServer(func(context *gin.Context) {
		auth.RegisterAuthentication(context, nil, userIDs[i], tokens[i])
		i++
	})
	defer server.Close()

	wsURL := wsURL(server.URL)

	userOneIPhone := testClient(t, wsURL)
	defer userOneIPhone.conn.Close()
	userOneAndroid := testClient(t, wsURL)
	defer userOneAndroid.conn.Close()
	userOneBrowser := testClient(t, wsURL)
	defer userOneBrowser.conn.Close()
	userOneOther := testClient(t, wsURL)
	defer userOneOther.conn.Close()
	userOne := []*testingClient{userOneAndroid, userOneBrowser, userOneIPhone, userOneOther}

	userTwoBrowser := testClient(t, wsURL)
	defer userTwoBrowser.conn.Close()
	userTwoAndroid := testClient(t, wsURL)
	defer userTwoAndroid.conn.Close()
	userTwo := []*testingClient{userTwoAndroid, userTwoBrowser}

	userThreeAndroid := testClient(t, wsURL)
	defer userThreeAndroid.conn.Close()
	userThree := []*testingClient{userThreeAndroid}

	// the server may take some time to register the client
	time.Sleep(100 * time.Millisecond)

	api.Notify(1, &model.MessageExternal{ID: 4, Message: "there"})
	expectMessage(&model.MessageExternal{ID: 4, Message: "there"}, userOne...)
	expectNoMessage(userTwo...)
	expectNoMessage(userThree...)

	api.NotifyDeletedUser(1)

	api.Notify(1, &model.MessageExternal{ID: 2, Message: "there"})
	expectNoMessage(userOne...)
	expectNoMessage(userThree...)
	expectNoMessage(userTwo...)

	api.Notify(2, &model.MessageExternal{ID: 2, Message: "there"})
	expectNoMessage(userOne...)
	expectMessage(&model.MessageExternal{ID: 2, Message: "there"}, userTwo...)
	expectNoMessage(userThree...)

	api.Notify(3, &model.MessageExternal{ID: 5, Message: "there"})
	expectNoMessage(userOne...)
	expectNoMessage(userTwo...)
	expectMessage(&model.MessageExternal{ID: 5, Message: "there"}, userThree...)

	api.Close()
}

func TestMultipleClients(t *testing.T) {
	mode.Set(mode.TestDev)

	defer leaktest.Check(t)()
	userIDs := []uint{1, 1, 1, 2, 2, 3}
	i := 0
	server, api := bootTestServer(func(context *gin.Context) {
		auth.RegisterAuthentication(context, nil, userIDs[i], "t"+string(userIDs[i]))
		i++
	})
	defer server.Close()

	wsURL := wsURL(server.URL)

	userOneIPhone := testClient(t, wsURL)
	defer userOneIPhone.conn.Close()
	userOneAndroid := testClient(t, wsURL)
	defer userOneAndroid.conn.Close()
	userOneBrowser := testClient(t, wsURL)
	defer userOneBrowser.conn.Close()
	userOne := []*testingClient{userOneAndroid, userOneBrowser, userOneIPhone}

	userTwoBrowser := testClient(t, wsURL)
	defer userTwoBrowser.conn.Close()
	userTwoAndroid := testClient(t, wsURL)
	defer userTwoAndroid.conn.Close()
	userTwo := []*testingClient{userTwoAndroid, userTwoBrowser}

	userThreeAndroid := testClient(t, wsURL)
	defer userThreeAndroid.conn.Close()
	userThree := []*testingClient{userThreeAndroid}

	// the server may take some time to register the client
	time.Sleep(100 * time.Millisecond)

	// there should not be messages at the beginning
	expectNoMessage(userOne...)
	expectNoMessage(userTwo...)
	expectNoMessage(userThree...)

	api.Notify(1, &model.MessageExternal{ID: 1, Message: "hello"})
	time.Sleep(500 * time.Millisecond)
	expectMessage(&model.MessageExternal{ID: 1, Message: "hello"}, userOne...)
	expectNoMessage(userTwo...)
	expectNoMessage(userThree...)

	api.Notify(2, &model.MessageExternal{ID: 2, Message: "there"})
	expectNoMessage(userOne...)
	expectMessage(&model.MessageExternal{ID: 2, Message: "there"}, userTwo...)
	expectNoMessage(userThree...)

	userOneIPhone.conn.Close()

	expectNoMessage(userOne...)
	expectNoMessage(userTwo...)
	expectNoMessage(userThree...)

	api.Notify(1, &model.MessageExternal{ID: 3, Message: "how"})
	expectMessage(&model.MessageExternal{ID: 3, Message: "how"}, userOneAndroid, userOneBrowser)
	expectNoMessage(userOneIPhone)
	expectNoMessage(userTwo...)
	expectNoMessage(userThree...)

	api.Notify(2, &model.MessageExternal{ID: 4, Message: "are"})

	expectNoMessage(userOne...)
	expectMessage(&model.MessageExternal{ID: 4, Message: "are"}, userTwo...)
	expectNoMessage(userThree...)

	api.Close()

	api.Notify(2, &model.MessageExternal{ID: 5, Message: "you"})

	expectNoMessage(userOne...)
	expectNoMessage(userTwo...)
	expectNoMessage(userThree...)
}

func Test_sameOrigin_returnsTrue(t *testing.T) {
	mode.Set(mode.Prod)
	req := httptest.NewRequest("GET", "http://example.com/stream", nil)
	req.Header.Set("Origin", "http://example.com")
	actual := isAllowedOrigin(req, nil)
	assert.True(t, actual)
}

func Test_sameOrigin_returnsTrue_withCustomPort(t *testing.T) {
	mode.Set(mode.Prod)
	req := httptest.NewRequest("GET", "http://example.com:8080/stream", nil)
	req.Header.Set("Origin", "http://example.com:8080")
	actual := isAllowedOrigin(req, nil)
	assert.True(t, actual)
}

func Test_isAllowedOrigin_withoutAllowedOrigins_failsWhenNotSameOrigin(t *testing.T) {
	mode.Set(mode.Prod)
	req := httptest.NewRequest("GET", "http://example.com/stream", nil)
	req.Header.Set("Origin", "http://gorify.example.com")
	actual := isAllowedOrigin(req, nil)
	assert.False(t, actual)
}

func Test_isAllowedOriginMatching(t *testing.T) {
	mode.Set(mode.Prod)
	compiledAllowedOrigins := compileAllowedWebSocketOrigins([]string{"go.{4}\\.example\\.com", "go\\.example\\.com"})

	req := httptest.NewRequest("GET", "http://example.me/stream", nil)
	req.Header.Set("Origin", "http://gorify.example.com")
	assert.True(t, isAllowedOrigin(req, compiledAllowedOrigins))

	req.Header.Set("Origin", "http://go.example.com")
	assert.True(t, isAllowedOrigin(req, compiledAllowedOrigins))

	req.Header.Set("Origin", "http://hello.example.com")
	assert.False(t, isAllowedOrigin(req, compiledAllowedOrigins))
}

func Test_emptyOrigin_returnsTrue(t *testing.T) {
	mode.Set(mode.Prod)
	req := httptest.NewRequest("GET", "http://example.com/stream", nil)
	actual := isAllowedOrigin(req, nil)
	assert.True(t, actual)
}

func Test_otherOrigin_returnsFalse(t *testing.T) {
	mode.Set(mode.Prod)
	req := httptest.NewRequest("GET", "http://example.com/stream", nil)
	req.Header.Set("Origin", "http://otherexample.de")
	actual := isAllowedOrigin(req, nil)
	assert.False(t, actual)
}

func Test_invalidOrigin_returnsFalse(t *testing.T) {
	mode.Set(mode.Prod)
	req := httptest.NewRequest("GET", "http://example.com/stream", nil)
	req.Header.Set("Origin", "http\\://otherexample.de")
	actual := isAllowedOrigin(req, nil)
	assert.False(t, actual)
}

func Test_compileAllowedWebSocketOrigins(t *testing.T) {
	assert.Equal(t, 0, len(compileAllowedWebSocketOrigins([]string{})))
	assert.Equal(t, 3, len(compileAllowedWebSocketOrigins([]string{"^.*$", "", "abc"})))
}

func clients(api *API, user uint) []*client {
	api.lock.RLock()
	defer api.lock.RUnlock()

	return api.clients[user]
}

func testClient(t *testing.T, url string) *testingClient {
	client := createClient(t, url)
	startReading(client)
	return client
}

func startReading(client *testingClient) {
	go func() {
		for {
			_, payload, err := client.conn.ReadMessage()

			if err != nil {
				return
			}

			actual := &model.MessageExternal{}
			json.NewDecoder(bytes.NewBuffer(payload)).Decode(actual)
			client.readMessage <- *actual
		}
	}()
}

func createClient(t *testing.T, url string) *testingClient {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.Nil(t, err)

	readMessages := make(chan model.MessageExternal)

	return &testingClient{conn: ws, readMessage: readMessages, t: t}
}

type testingClient struct {
	conn        *websocket.Conn
	readMessage chan model.MessageExternal
	t           *testing.T
}

func (c *testingClient) expectMessage(expected *model.MessageExternal) {
	select {
	case <-time.After(50 * time.Millisecond):
		assert.Fail(c.t, "Expected message but none was send :(")
	case actual := <-c.readMessage:
		assert.Equal(c.t, *expected, actual)
	}
}

func expectMessage(expected *model.MessageExternal, clients ...*testingClient) {
	for _, client := range clients {
		client.expectMessage(expected)
	}
}

func expectNoMessage(clients ...*testingClient) {
	for _, client := range clients {
		client.expectNoMessage()
	}
}

func (c *testingClient) expectNoMessage() {
	select {
	case <-time.After(50 * time.Millisecond):
		// no message == as expected
	case msg := <-c.readMessage:
		assert.Fail(c.t, "Expected NO message but there was one :(", fmt.Sprint(msg))
	}
}

func bootTestServer(handlerFunc gin.HandlerFunc) (*httptest.Server, *API) {
	r := gin.New()
	r.Use(handlerFunc)
	// ping every 500 ms, and the client has 500 ms to respond
	api := New(500*time.Millisecond, 500*time.Millisecond, []string{})

	r.GET("/", api.Handle)
	server := httptest.NewServer(r)
	return server, api
}

func wsURL(httpURL string) string {
	return "ws" + strings.TrimPrefix(httpURL, "http")
}

func staticUserID() gin.HandlerFunc {
	return func(context *gin.Context) {
		auth.RegisterAuthentication(context, nil, 1, "customtoken")
	}
}
