package plugin

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"plugin"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/v2/auth"
	"github.com/gotify/server/v2/model"
	"github.com/gotify/server/v2/plugin/compat"
	"gopkg.in/yaml.v2"
)

// The Database interface for encapsulating database access.
type Database interface {
	GetUsers() ([]*model.User, error)
	GetPluginConfByUserAndPath(userid uint, path string) (*model.PluginConf, error)
	CreatePluginConf(p *model.PluginConf) error
	GetPluginConfByApplicationID(appid uint) (*model.PluginConf, error)
	UpdatePluginConf(p *model.PluginConf) error
	CreateMessage(message *model.Message) error
	GetPluginConfByID(id uint) (*model.PluginConf, error)
	GetPluginConfByToken(token string) (*model.PluginConf, error)
	GetUserByID(id uint) (*model.User, error)
	CreateApplication(application *model.Application) error
	UpdateApplication(app *model.Application) error
	GetApplicationsByUser(userID uint) ([]*model.Application, error)
	GetApplicationByToken(token string) (*model.Application, error)
}

// Notifier notifies when a new message was created.
type Notifier interface {
	Notify(userID uint, message *model.MessageExternal)
}

// Manager is an encapsulating layer for plugins and manages all plugins and its instances.
type Manager struct {
	mutex     *sync.RWMutex
	instances map[uint]compat.PluginInstance
	plugins   map[string]compat.Plugin
	messages  chan MessageWithUserID
	db        Database
	mux       *gin.RouterGroup
}

// NewManager created a Manager from configurations.
func NewManager(db Database, directory string, mux *gin.RouterGroup, notifier Notifier) (*Manager, error) {
	manager := &Manager{
		mutex:     &sync.RWMutex{},
		instances: map[uint]compat.PluginInstance{},
		plugins:   map[string]compat.Plugin{},
		messages:  make(chan MessageWithUserID),
		db:        db,
		mux:       mux,
	}

	go func() {
		for {
			message := <-manager.messages
			internalMsg := &model.Message{
				ApplicationID: message.Message.ApplicationID,
				Title:         message.Message.Title,
				Priority:      message.Message.Priority,
				Date:          message.Message.Date,
				Message:       message.Message.Message,
			}
			if message.Message.Extras != nil {
				internalMsg.Extras, _ = json.Marshal(message.Message.Extras)
			}
			db.CreateMessage(internalMsg)
			message.Message.ID = internalMsg.ID
			notifier.Notify(message.UserID, &message.Message)
		}
	}()

	if err := manager.loadPlugins(directory); err != nil {
		return nil, err
	}

	users, err := manager.db.GetUsers()
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		if err := manager.initializeForUser(*user); err != nil {
			return nil, err
		}
	}

	return manager, nil
}

// ErrAlreadyEnabledOrDisabled is returned on SetPluginEnabled call when a plugin is already enabled or disabled.
var ErrAlreadyEnabledOrDisabled = errors.New("config is already enabled/disabled")

func (m *Manager) applicationExists(token string) bool {
	app, _ := m.db.GetApplicationByToken(token)
	return app != nil
}

func (m *Manager) pluginConfExists(token string) bool {
	pluginConf, _ := m.db.GetPluginConfByToken(token)
	return pluginConf != nil
}

// SetPluginEnabled sets the plugins enabled state.
func (m *Manager) SetPluginEnabled(pluginID uint, enabled bool) error {
	instance, err := m.Instance(pluginID)
	if err != nil {
		return errors.New("instance not found")
	}
	conf, err := m.db.GetPluginConfByID(pluginID)
	if err != nil {
		return err
	}

	if conf.Enabled == enabled {
		return ErrAlreadyEnabledOrDisabled
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if enabled {
		err = instance.Enable()
	} else {
		err = instance.Disable()
	}
	if err != nil {
		return err
	}

	if newConf, err := m.db.GetPluginConfByID(pluginID); /* conf might be updated by instance */ err == nil {
		conf = newConf
	}
	conf.Enabled = enabled
	return m.db.UpdatePluginConf(conf)
}

// PluginInfo returns plugin info.
func (m *Manager) PluginInfo(modulePath string) compat.Info {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if p, ok := m.plugins[modulePath]; ok {
		return p.PluginInfo()
	}
	fmt.Println("Could not get plugin info for", modulePath)
	return compat.Info{
		Name:        "UNKNOWN",
		ModulePath:  modulePath,
		Description: "Oops something went wrong",
	}
}

// Instance returns an instance with the given ID.
func (m *Manager) Instance(pluginID uint) (compat.PluginInstance, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if instance, ok := m.instances[pluginID]; ok {
		return instance, nil
	}
	return nil, errors.New("instance not found")
}

// HasInstance returns whether the given plugin ID has a corresponding instance.
func (m *Manager) HasInstance(pluginID uint) bool {
	instance, err := m.Instance(pluginID)
	return err == nil && instance != nil
}

// RemoveUser disabled all plugins of a user when the user is disabled.
func (m *Manager) RemoveUser(userID uint) error {
	for _, p := range m.plugins {
		pluginConf, err := m.db.GetPluginConfByUserAndPath(userID, p.PluginInfo().ModulePath)
		if err != nil {
			return err
		}
		if pluginConf == nil {
			continue
		}
		if pluginConf.Enabled {
			inst, err := m.Instance(pluginConf.ID)
			if err != nil {
				continue
			}
			m.mutex.Lock()
			err = inst.Disable()
			m.mutex.Unlock()
			if err != nil {
				return err
			}
		}
		delete(m.instances, pluginConf.ID)
	}
	return nil
}

type pluginFileLoadError struct {
	Filename        string
	UnderlyingError error
}

func (c pluginFileLoadError) Error() string {
	return fmt.Sprintf("error while loading plugin %s: %s", c.Filename, c.UnderlyingError)
}

func (m *Manager) loadPlugins(directory string) error {
	if directory == "" {
		return nil
	}

	pluginFiles, err := ioutil.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("error while reading directory %s", err)
	}
	for _, f := range pluginFiles {
		pluginPath := filepath.Join(directory, "./", f.Name())

		fmt.Println("Loading plugin", pluginPath)
		pRaw, err := plugin.Open(pluginPath)
		if err != nil {
			return pluginFileLoadError{f.Name(), err}
		}
		compatPlugin, err := compat.Wrap(pRaw)
		if err != nil {
			return pluginFileLoadError{f.Name(), err}
		}
		if err := m.LoadPlugin(compatPlugin); err != nil {
			return pluginFileLoadError{f.Name(), err}
		}
	}
	return nil
}

// LoadPlugin loads a compat plugin, exported to sideload plugins for testing purposes.
func (m *Manager) LoadPlugin(compatPlugin compat.Plugin) error {
	modulePath := compatPlugin.PluginInfo().ModulePath
	if _, ok := m.plugins[modulePath]; ok {
		return fmt.Errorf("plugin with module path %s is present at least twice", modulePath)
	}
	m.plugins[modulePath] = compatPlugin
	return nil
}

// InitializeForUserID initializes all plugin instances for a given user.
func (m *Manager) InitializeForUserID(userID uint) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	user, err := m.db.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user != nil {
		return m.initializeForUser(*user)
	}
	return fmt.Errorf("user with id %d not found", userID)
}

func (m *Manager) initializeForUser(user model.User) error {
	userCtx := compat.UserContext{
		ID:    user.ID,
		Name:  user.Name,
		Admin: user.Admin,
	}

	for _, p := range m.plugins {
		if err := m.initializeSingleUserPlugin(userCtx, p); err != nil {
			return err
		}
	}

	apps, err := m.db.GetApplicationsByUser(user.ID)
	if err != nil {
		return err
	}
	for _, app := range apps {
		conf, err := m.db.GetPluginConfByApplicationID(app.ID)
		if err != nil {
			return err
		}
		if conf != nil {
			_, compatExist := m.plugins[conf.ModulePath]
			app.Internal = compatExist
		} else {
			app.Internal = false
		}
		m.db.UpdateApplication(app)
	}

	return nil
}

func (m *Manager) initializeSingleUserPlugin(userCtx compat.UserContext, p compat.Plugin) error {
	info := p.PluginInfo()
	instance := p.NewPluginInstance(userCtx)
	userID := userCtx.ID

	pluginConf, err := m.db.GetPluginConfByUserAndPath(userID, info.ModulePath)
	if err != nil {
		return err
	}

	if pluginConf == nil {
		var err error
		pluginConf, err = m.createPluginConf(instance, info, userID)
		if err != nil {
			return err
		}
	}

	m.instances[pluginConf.ID] = instance

	if compat.HasSupport(instance, compat.Messenger) {
		instance.SetMessageHandler(redirectToChannel{
			ApplicationID: pluginConf.ApplicationID,
			UserID:        pluginConf.UserID,
			Messages:      m.messages,
		})
	}
	if compat.HasSupport(instance, compat.Storager) {
		instance.SetStorageHandler(dbStorageHandler{pluginConf.ID, m.db})
	}
	if compat.HasSupport(instance, compat.Configurer) {
		m.initializeConfigurerForSingleUserPlugin(instance, pluginConf)
	}
	if compat.HasSupport(instance, compat.Webhooker) {
		id := pluginConf.ID
		g := m.mux.Group(pluginConf.Token+"/", requirePluginEnabled(id, m.db))
		instance.RegisterWebhook(strings.Replace(g.BasePath(), ":id", strconv.Itoa(int(id)), 1), g)
	}
	if pluginConf.Enabled {
		err := instance.Enable()
		if err != nil {
			// Single user plugin cannot be enabled
			// Don't panic, disable for now and wait for user to update config
			log.Printf("Plugin initialize failed for user %s: %s. Disabling now...", userCtx.Name, err.Error())
			pluginConf.Enabled = false
			m.db.UpdatePluginConf(pluginConf)
		}
	}
	return nil
}

func (m *Manager) initializeConfigurerForSingleUserPlugin(instance compat.PluginInstance, pluginConf *model.PluginConf) {
	if len(pluginConf.Config) == 0 {
		// The Configurer is newly implemented
		// Use the default config
		pluginConf.Config, _ = yaml.Marshal(instance.DefaultConfig())
		m.db.UpdatePluginConf(pluginConf)
	}
	c := instance.DefaultConfig()
	if yaml.Unmarshal(pluginConf.Config, c) != nil || instance.ValidateAndSetConfig(c) != nil {
		pluginConf.Enabled = false

		log.Printf("Plugin %s for user %d failed to initialize because it rejected the current config. It might be outdated. A default config is used and the user would need to enable it again.", pluginConf.ModulePath, pluginConf.UserID)
		newConf := bytes.NewBufferString("# Plugin initialization failed because it rejected the current config. It might be outdated.\r\n# A default plugin configuration is used:\r\n")

		d, _ := yaml.Marshal(c)
		newConf.Write(d)
		newConf.WriteString("\r\n")

		newConf.WriteString("# The original configuration: \r\n")
		oldConf := bufio.NewScanner(bytes.NewReader(pluginConf.Config))
		for oldConf.Scan() {
			newConf.WriteString("# ")
			newConf.WriteString(oldConf.Text())
			newConf.WriteString("\r\n")
		}

		pluginConf.Config = newConf.Bytes()

		m.db.UpdatePluginConf(pluginConf)
		instance.ValidateAndSetConfig(instance.DefaultConfig())
	}
}

func (m *Manager) createPluginConf(instance compat.PluginInstance, info compat.Info, userID uint) (*model.PluginConf, error) {
	pluginConf := &model.PluginConf{
		UserID:     userID,
		ModulePath: info.ModulePath,
		Token:      auth.GenerateNotExistingToken(auth.GeneratePluginToken, m.pluginConfExists),
	}
	if compat.HasSupport(instance, compat.Configurer) {
		pluginConf.Config, _ = yaml.Marshal(instance.DefaultConfig())
	}
	if compat.HasSupport(instance, compat.Messenger) {
		app := &model.Application{
			Token:       auth.GenerateNotExistingToken(auth.GenerateApplicationToken, m.applicationExists),
			Name:        info.String(),
			UserID:      userID,
			Internal:    true,
			Description: fmt.Sprintf("auto generated application for %s", info.ModulePath),
		}
		if err := m.db.CreateApplication(app); err != nil {
			return nil, err
		}
		pluginConf.ApplicationID = app.ID
	}
	if err := m.db.CreatePluginConf(pluginConf); err != nil {
		return nil, err
	}
	return pluginConf, nil
}
