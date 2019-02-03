package main

import (
	"time"

	"github.com/gotify/plugin-api"
	"github.com/robfig/cron"
)

// GetGotifyPluginInfo returns gotify plugin info
func GetGotifyPluginInfo() plugin.Info {
	return plugin.Info{
		Name:        "clock",
		Description: "Sends an hourly reminder",
		ModulePath:  "github.com/gotify/server/example/clock",
	}
}

// Plugin is plugin instance
type Plugin struct {
	msgHandler  plugin.MessageHandler
	enabled     bool
	cronHandler *cron.Cron
}

// Enable implements plugin.Plugin
func (c *Plugin) Enable() error {
	c.enabled = true
	c.cronHandler = cron.New()
	c.cronHandler.AddFunc("0 0 * * *", func() {
		c.msgHandler.SendMessage(plugin.Message{
			Title:   "Tick Tock!",
			Message: time.Now().Format("It is 15:04:05 now."),
		})
	})
	c.cronHandler.Start()
	return nil
}

// Disable implements plugin.Plugin
func (c *Plugin) Disable() error {
	if c.cronHandler != nil {
		c.cronHandler.Stop()
	}
	c.enabled = false
	return nil
}

// SetMessageHandler implements plugin.Messenger.
func (c *Plugin) SetMessageHandler(h plugin.MessageHandler) {
	c.msgHandler = h
}

// NewGotifyPluginInstance creates a plugin instance for a user context.
func NewGotifyPluginInstance(ctx plugin.UserContext) plugin.Plugin {
	p := &Plugin{}

	return p
}

func main() {
	panic("this should be built as go plugin")
}
