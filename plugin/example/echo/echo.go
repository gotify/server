package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gotify/plugin-api"
)

// GetGotifyPluginInfo returns gotify plugin info.
func GetGotifyPluginInfo() plugin.Info {
	return plugin.Info{
		ModulePath: "github.com/gotify/server/plugin/example/echo",
		Name:       "test plugin",
	}
}

// EchoPlugin is the gotify plugin instance.
type EchoPlugin struct {
	msgHandler     plugin.MessageHandler
	storageHandler plugin.StorageHandler
	config         *Config
	basePath       string
}

// SetStorageHandler implements plugin.Storager
func (c *EchoPlugin) SetStorageHandler(h plugin.StorageHandler) {
	c.storageHandler = h
}

// SetMessageHandler implements plugin.Messenger.
func (c *EchoPlugin) SetMessageHandler(h plugin.MessageHandler) {
	c.msgHandler = h
}

// Storage defines the plugin storage scheme
type Storage struct {
	CalledTimes int `json:"called_times"`
}

// Config defines the plugin config scheme
type Config struct {
	MagicString string `yaml:"magic_string"`
}

// DefaultConfig implements plugin.Configurer
func (c *EchoPlugin) DefaultConfig() interface{} {
	return &Config{
		MagicString: "hello world",
	}
}

// ValidateAndSetConfig implements plugin.Configurer
func (c *EchoPlugin) ValidateAndSetConfig(config interface{}) error {
	c.config = config.(*Config)
	return nil
}

// Enable enables the plugin.
func (c *EchoPlugin) Enable() error {
	log.Println("echo plugin enabled")
	return nil
}

// Disable disables the plugin.
func (c *EchoPlugin) Disable() error {
	log.Println("echo plugin disbled")
	return nil
}

// RegisterWebhook implements plugin.Webhooker.
func (c *EchoPlugin) RegisterWebhook(baseURL string, g *gin.RouterGroup) {
	c.basePath = baseURL
	g.GET("/echo", func(ctx *gin.Context) {

		storage, _ := c.storageHandler.Load()
		conf := new(Storage)
		json.Unmarshal(storage, conf)
		conf.CalledTimes++
		newStorage, _ := json.Marshal(conf)
		c.storageHandler.Save(newStorage)

		c.msgHandler.SendMessage(plugin.Message{
			Title:    "Hello received",
			Message:  fmt.Sprintf("echo server received a hello message %d times", conf.CalledTimes),
			Priority: 2,
			Extras: map[string]interface{}{
				"plugin::name": "echo",
			},
		})
		ctx.Writer.WriteString(fmt.Sprintf("Magic string is: %s\r\nEcho server running at %secho", c.config.MagicString, c.basePath))
	})
}

// GetDisplay implements plugin.Displayer.
func (c *EchoPlugin) GetDisplay(location *url.URL) string {
	loc := &url.URL{
		Path: c.basePath,
	}
	if location != nil {
		loc.Scheme = location.Scheme
		loc.Host = location.Host
	}
	loc = loc.ResolveReference(&url.URL{
		Path: "echo",
	})
	return "Echo plugin running at: " + loc.String()
}

// NewGotifyPluginInstance creates a plugin instance for a user context.
func NewGotifyPluginInstance(ctx plugin.UserContext) plugin.Plugin {
	return &EchoPlugin{}
}

func main() {
	panic("this should be built as go plugin")
}
