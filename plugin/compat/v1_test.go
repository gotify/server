package compat

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	papiv1 "github.com/gotify/plugin-api"
)

type v1MockInstance struct {
	Enabled bool
}

func (c *v1MockInstance) Enable() error {
	c.Enabled = true
	return nil
}

func (c *v1MockInstance) Disable() error {
	c.Enabled = false
	return nil
}

type V1WrapperSuite struct {
	suite.Suite
	i PluginV1Instance
}

func (s *V1WrapperSuite) SetupSuite() {
	inst := new(v1MockInstance)
	s.i.instance = inst
}

func (s *V1WrapperSuite) TestConfigurer_notSupported_expectEmpty() {
	assert.Equal(s.T(), struct{}{}, s.i.DefaultConfig())
	assert.Nil(s.T(), s.i.ValidateAndSetConfig(struct{}{}))
}

func (s *V1WrapperSuite) TestDisplayer_notSupported_expectEmpty() {
	assert.Equal(s.T(), "", s.i.GetDisplay(nil))
}

type v1StorageHandler struct {
	storage []byte
}

func (c *v1StorageHandler) Save(b []byte) error {
	c.storage = b
	return nil
}

func (c *v1StorageHandler) Load() ([]byte, error) {
	return c.storage, nil
}

type v1Storager struct {
	handler papiv1.StorageHandler
}

func (c *v1Storager) Enable() error {
	return nil
}

func (c *v1Storager) Disable() error {
	return nil
}

func (c *v1Storager) SetStorageHandler(h papiv1.StorageHandler) {
	c.handler = h
}

func (s *V1WrapperSuite) TestStorager() {
	storager := new(v1Storager)
	s.i.storager = storager

	s.i.SetStorageHandler(new(v1StorageHandler))

	assert.Nil(s.T(), storager.handler.Save([]byte("test")))
	storage, err := storager.handler.Load()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "test", string(storage))
}

type v1MessengerHandler struct {
	msgSent Message
}

func (c *v1MessengerHandler) SendMessage(msg Message) error {
	c.msgSent = msg
	return nil
}

type v1Messenger struct {
	handler papiv1.MessageHandler
}

func (c *v1Messenger) Enable() error {
	return nil
}

func (c *v1Messenger) Disable() error {
	return nil
}

func (c *v1Messenger) SetMessageHandler(h papiv1.MessageHandler) {
	c.handler = h
}

func (s *V1WrapperSuite) TestMessenger_sendMessageWithExtras() {
	messenger := new(v1Messenger)
	s.i.messenger = messenger

	handler := new(v1MessengerHandler)
	s.i.SetMessageHandler(handler)

	msg := papiv1.Message{
		Title:    "test message",
		Message:  "test",
		Priority: 2,
		Extras: map[string]interface{}{
			"test::string": "test",
		},
	}
	assert.Nil(s.T(), messenger.handler.SendMessage(msg))
	assert.Equal(s.T(), Message{
		Title:    "test message",
		Message:  "test",
		Priority: 2,
		Extras: map[string]interface{}{
			"test::string": "test",
		},
	}, handler.msgSent)
}

func (s *V1WrapperSuite) TestMessenger_sendMessageWithoutExtras() {
	messenger := new(v1Messenger)
	s.i.messenger = messenger

	handler := new(v1MessengerHandler)
	s.i.SetMessageHandler(handler)

	msg := papiv1.Message{
		Title:    "test message",
		Message:  "test",
		Priority: 2,
		Extras:   nil,
	}
	assert.Nil(s.T(), messenger.handler.SendMessage(msg))
	assert.Equal(s.T(), Message{
		Title:    "test message",
		Message:  "test",
		Priority: 2,
		Extras:   nil,
	}, handler.msgSent)
}
func TestV1Wrapper(t *testing.T) {
	suite.Run(t, new(V1WrapperSuite))
}
