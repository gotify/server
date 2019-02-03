package compat

// Plugin is an abstraction of plugin handler
type Plugin interface {
	PluginInfo() Info
	NewPluginInstance(ctx UserContext) PluginInstance
	APIVersion() string
}

// Info is the plugin info
type Info struct {
	Version     string
	Author      string
	Name        string
	Website     string
	Description string
	License     string
	ModulePath  string
}

func (c Info) String() string {
	if c.Name != "" {
		return c.Name
	}
	return c.ModulePath
}

// UserContext is the user context used to create plugin instance.
type UserContext struct {
	ID    uint
	Name  string
	Admin bool
}
