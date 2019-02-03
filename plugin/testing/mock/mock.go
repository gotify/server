package mock

import (
	"errors"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/gotify/server/plugin/compat"
)

// ModulePath is for convenient access of the module path of this mock plugin
const ModulePath = "github.com/gotify/server/plugin/testing/mock"

// Name is for convenient access of the module path of the name of this mock plugin
const Name = "mock plugin"

// Plugin is a mock plugin.
type Plugin struct {
	Instances []PluginInstance
}

// PluginInfo implements loader.PluginCompat
func (c *Plugin) PluginInfo() compat.Info {
	return compat.Info{
		ModulePath: ModulePath,
		Name:       Name,
	}
}

// NewPluginInstance implements loader.PluginCompat
func (c *Plugin) NewPluginInstance(ctx compat.UserContext) compat.PluginInstance {
	inst := PluginInstance{UserCtx: ctx, capabilities: compat.Capabilities{compat.Configurer, compat.Storager, compat.Messenger, compat.Displayer}}
	c.Instances = append(c.Instances, inst)
	return &inst
}

// APIVersion implements loader.PluginCompat
func (c *Plugin) APIVersion() string {
	return "v1"
}

// PluginInstance is a mock plugin instance
type PluginInstance struct {
	UserCtx        compat.UserContext
	Enabled        bool
	DisplayString  string
	Config         *PluginConfig
	storageHandler compat.StorageHandler
	messageHandler compat.MessageHandler
	capabilities   compat.Capabilities
	BasePath       string
}

// PluginConfig is a mock plugin config struct
type PluginConfig struct {
	TestKey    string
	IsNotValid bool
}

var disableFailUsers = make(map[uint]error)
var enableFailUsers = make(map[uint]error)

// ReturnErrorOnEnableForUser registers a uid which will throw an error on enabling.
func ReturnErrorOnEnableForUser(uid uint, err error) {
	enableFailUsers[uid] = err
}

// ReturnErrorOnDisableForUser registers a uid which will throw an error on disabling.
func ReturnErrorOnDisableForUser(uid uint, err error) {
	disableFailUsers[uid] = err
}

// Enable implements compat.PluginInstance
func (c *PluginInstance) Enable() error {
	if err, ok := enableFailUsers[c.UserCtx.ID]; ok {
		return err
	}
	c.Enabled = true
	return nil
}

// Disable implements compat.PluginInstance
func (c *PluginInstance) Disable() error {
	if err, ok := disableFailUsers[c.UserCtx.ID]; ok {
		return err
	}
	c.Enabled = false
	return nil
}

// SetMessageHandler implements compat.Messenger
func (c *PluginInstance) SetMessageHandler(h compat.MessageHandler) {
	c.messageHandler = h
}

// SetStorageHandler implements compat.Storager
func (c *PluginInstance) SetStorageHandler(handler compat.StorageHandler) {
	c.storageHandler = handler
}

// SetStorage sets current storage
func (c *PluginInstance) SetStorage(b []byte) error {
	return c.storageHandler.Save(b)
}

// GetStorage sets current storage
func (c *PluginInstance) GetStorage() ([]byte, error) {
	return c.storageHandler.Load()
}

// RegisterWebhook implements compat.Webhooker
func (c *PluginInstance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	c.BasePath = basePath
}

// SetCapability changes the capability of this plugin
func (c *PluginInstance) SetCapability(p compat.Capability, enable bool) {
	if enable {
		for _, cap := range c.capabilities {
			if cap == p {
				return
			}
		}
		c.capabilities = append(c.capabilities, p)
	} else {
		newCap := make(compat.Capabilities, 0)
		for _, cap := range c.capabilities {
			if cap == p {
				continue
			}
			newCap = append(newCap, cap)
		}
		c.capabilities = newCap
	}
}

// Supports implements compat.PluginInstance
func (c *PluginInstance) Supports() compat.Capabilities {
	return c.capabilities
}

// DefaultConfig implements compat.Configuror
func (c *PluginInstance) DefaultConfig() interface{} {
	return &PluginConfig{
		TestKey:    "default",
		IsNotValid: false,
	}
}

// ValidateAndSetConfig implements compat.Configuror
func (c *PluginInstance) ValidateAndSetConfig(config interface{}) error {
	if (config.(*PluginConfig)).IsNotValid {
		return errors.New("conf is not valid")
	}
	c.Config = config.(*PluginConfig)
	return nil
}

// GetDisplay implements compat.Displayer
func (c *PluginInstance) GetDisplay(url *url.URL) string {
	return c.DisplayString
}

// TriggerMessage triggers a test message
func (c *PluginInstance) TriggerMessage() {
	c.messageHandler.SendMessage(compat.Message{
		Title:    "test message",
		Message:  "test",
		Priority: 2,
		Extras: map[string]interface{}{
			"test::string": "test",
		},
	})
}
