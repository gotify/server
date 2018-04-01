package stream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"reflect"

	"errors"

	"github.com/bouk/monkey"
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

	// try emulate an write error, mostly this should kill the ReadMessage goroutine first but you'll never know.
	patch := monkey.PatchInstanceMethod(reflect.TypeOf(clients[0].conn), "WriteJSON", func(*websocket.Conn, interface{}) error {
		return errors.New("could not do something")
	})
	defer patch.Unpatch()

	api.Notify(1, &model.Message{Message: "HI"})
	user.expectNoMessage()
}

func TestWritePingFails(t *testing.T) {
	mode.Set(mode.TestDev)

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
	// try emulate an write error, mostly this should kill the ReadMessage gorouting first but you'll never know.
	patch := monkey.PatchInstanceMethod(reflect.TypeOf(clients[0].conn), "WriteMessage", func(*websocket.Conn, int, []byte) error {
		return errors.New("could not do something")
	})
	defer patch.Unpatch()

	time.Sleep(5 * time.Second) // waiting for ping

	api.Notify(1, &model.Message{Message: "HI"})
	user.expectNoMessage()
}

func TestPing(t *testing.T) {
	mode.Set(mode.TestDev)

	server, api := bootTestServer(staticUserID())
	defer server.Close()
	defer api.Close()

	wsURL := wsURL(server.URL)

	user := testClient(t, wsURL)
	defer user.conn.Close()

	ping := make(chan bool)
	oldPingHandler := user.conn.PingHandler()
	user.conn.SetPingHandler(func(appData string) error {
		err := oldPingHandler(appData)
		ping <- true
		return err
	})

	expectNoMessage(user)

	select {
	case <-time.After(5 * time.Second):
		assert.Fail(t, "Expected ping but there was one :(")
	case <-ping:
		// expected
	}

	expectNoMessage(user)
	api.Notify(1, &model.Message{Message: "HI"})
	user.expectMessage(&model.Message{Message: "HI"})
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

	time.Sleep(7 * time.Second)

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
	api.Notify(1, &model.Message{Message: "msg"})
	user.expectMessage(&model.Message{Message: "msg"})
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
	api.Notify(1, &model.Message{Message: "msg"})
	user.expectMessage(&model.Message{Message: "msg"})

	api.NotifyDeleted(1, "customtoken")

	api.Notify(1, &model.Message{Message: "msg"})
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

	api.Notify(1, &model.Message{ID: 4, Message: "there"})
	expectMessage(&model.Message{ID: 4, Message: "there"}, userOne...)
	expectNoMessage(userTwo...)
	expectNoMessage(userThree...)

	api.NotifyDeleted(1, "1-2")

	api.Notify(1, &model.Message{ID: 2, Message: "there"})
	expectMessage(&model.Message{ID: 2, Message: "there"}, userOneIPhone, userOneOther)
	expectNoMessage(userOneBrowser, userOneAndroid)
	expectNoMessage(userThree...)
	expectNoMessage(userTwo...)

	api.Notify(2, &model.Message{ID: 2, Message: "there"})
	expectNoMessage(userOne...)
	expectMessage(&model.Message{ID: 2, Message: "there"}, userTwo...)
	expectNoMessage(userThree...)

	api.Notify(3, &model.Message{ID: 5, Message: "there"})
	expectNoMessage(userOne...)
	expectNoMessage(userTwo...)
	expectMessage(&model.Message{ID: 5, Message: "there"}, userThree...)

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

	api.Notify(1, &model.Message{ID: 1, Message: "hello"})
	time.Sleep(1 * time.Second)
	expectMessage(&model.Message{ID: 1, Message: "hello"}, userOne...)
	expectNoMessage(userTwo...)
	expectNoMessage(userThree...)

	api.Notify(2, &model.Message{ID: 2, Message: "there"})
	expectNoMessage(userOne...)
	expectMessage(&model.Message{ID: 2, Message: "there"}, userTwo...)
	expectNoMessage(userThree...)

	userOneIPhone.conn.Close()

	expectNoMessage(userOne...)
	expectNoMessage(userTwo...)
	expectNoMessage(userThree...)

	api.Notify(1, &model.Message{ID: 3, Message: "how"})
	expectMessage(&model.Message{ID: 3, Message: "how"}, userOneAndroid, userOneBrowser)
	expectNoMessage(userOneIPhone)
	expectNoMessage(userTwo...)
	expectNoMessage(userThree...)

	api.Notify(2, &model.Message{ID: 4, Message: "are"})

	expectNoMessage(userOne...)
	expectMessage(&model.Message{ID: 4, Message: "are"}, userTwo...)
	expectNoMessage(userThree...)

	api.Close()

	api.Notify(2, &model.Message{ID: 5, Message: "you"})

	expectNoMessage(userOne...)
	expectNoMessage(userTwo...)
	expectNoMessage(userThree...)
}

func Test_sameOrigin_returnsTrue(t *testing.T) {
	mode.Set(mode.Prod)
	req := httptest.NewRequest("GET", "http://example.com/stream", nil)
	req.Header.Set("Origin", "http://example.com")
	actual := checkSameOrigin(req)
	assert.True(t, actual)
}

func Test_emptyOrigin_returnsTrue(t *testing.T) {
	mode.Set(mode.Prod)
	req := httptest.NewRequest("GET", "http://example.com/stream", nil)
	actual := checkSameOrigin(req)
	assert.True(t, actual)
}

func Test_otherOrigin_returnsFalse(t *testing.T) {
	mode.Set(mode.Prod)
	req := httptest.NewRequest("GET", "http://example.com/stream", nil)
	req.Header.Set("Origin", "http://otherexample.de")
	actual := checkSameOrigin(req)
	assert.False(t, actual)
}

func Test_invalidOrigin_returnsFalse(t *testing.T) {
	mode.Set(mode.Prod)
	req := httptest.NewRequest("GET", "http://example.com/stream", nil)
	req.Header.Set("Origin", "http\\://otherexample.de")
	actual := checkSameOrigin(req)
	assert.False(t, actual)
}

func clients(api *API, user uint) []*client {
	api.lock.RLock()
	defer api.lock.RUnlock()

	return api.clients[user]
}

func testClient(t *testing.T, url string) *testingClient {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.Nil(t, err)

	readMessages := make(chan model.Message)

	go func() {
		for {
			_, payload, err := ws.ReadMessage()

			if err != nil {
				return
			}

			actual := &model.Message{}
			json.NewDecoder(bytes.NewBuffer(payload)).Decode(actual)
			readMessages <- *actual
		}
	}()

	return &testingClient{conn: ws, readMessage: readMessages, t: t}
}

type testingClient struct {
	conn        *websocket.Conn
	readMessage <-chan model.Message
	t           *testing.T
}

func (c *testingClient) expectMessage(expected *model.Message) {
	select {
	case <-time.After(50 * time.Millisecond):
		assert.Fail(c.t, "Expected message but none was send :(")
	case actual := <-c.readMessage:
		assert.Equal(c.t, *expected, actual)
	}
}

func expectMessage(expected *model.Message, clients ...*testingClient) {
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
	// all 4 seconds a ping, and the client has 1 second to respond
	api := New(4*time.Second, 1*time.Second)
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
