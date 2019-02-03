package main

import (
	"github.com/gotify/plugin-api"
)

// GetGotifyPluginInfo returns gotify plugin info
func GetGotifyPluginInfo() plugin.Info {
	return plugin.Info{
		ModulePath: "github.com/gotify/server/plugin/testing/broken/noinstance",
	}
}

func main() {
	panic("this is a broken plugin for testing purposes")
}
