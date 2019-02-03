package compat

import (
	"errors"
	"fmt"
	"plugin"

	papiv1 "github.com/gotify/plugin-api"
)

// Wrap wraps around a raw go plugin to provide typesafe access.
func Wrap(p *plugin.Plugin) (Plugin, error) {
	getInfoHandle, err := p.Lookup("GetGotifyPluginInfo")
	if err != nil {
		return nil, errors.New("missing GetGotifyPluginInfo symbol")
	}
	switch getInfoHandle := getInfoHandle.(type) {
	case func() papiv1.Info:
		v1 := PluginV1{}

		v1.Info = getInfoHandle()
		newInstanceHandle, err := p.Lookup("NewGotifyPluginInstance")
		if err != nil {
			return nil, errors.New("missing NewGotifyPluginInstance symbol")
		}
		constructor, ok := newInstanceHandle.(func(ctx papiv1.UserContext) papiv1.Plugin)
		if !ok {
			return nil, fmt.Errorf("NewGotifyPluginInstance signature mismatch, func(ctx plugin.UserContext) plugin.Plugin expected, got %T", newInstanceHandle)
		}
		v1.Constructor = constructor
		return v1, nil
	default:
		return nil, fmt.Errorf("unknown plugin version (unrecogninzed GetGotifyPluginInfo signature %T)", getInfoHandle)
	}
}
