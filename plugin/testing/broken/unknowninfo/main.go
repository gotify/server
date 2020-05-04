package main

// GetGotifyPluginInfo returns gotify plugin info
func GetGotifyPluginInfo() string {
	return "github.com/gotify/server/v2/plugin/testing/broken/unknowninfo"
}

func main() {
	panic("this is a broken plugin for testing purposes")
}
