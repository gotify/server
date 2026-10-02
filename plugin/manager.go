package plugin

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"plugin"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/v3/auth"
	"github.com/gotify/server/v3/database"
	"github.com/gotify/server/v3/model"
	"github.com/gotify/server/v3/plugin/compat"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// The Database interface for encapsulating database access.
type Database interface {
	GetUsers(condition ...any) ([]*model.User, error)
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

type InstanceWrapper struct {
	enabled  bool
	userID   uint
	instance compat.PluginInstance
}

func (i *InstanceWrapper) Instance() compat.PluginInstance {
	return i.instance
}

func (i *InstanceWrapper) Enable() error {
	i.enabled = true
	err := i.instance.Enable()
	if err != nil {
		i.enabled = false
		return err
	}
	return nil
}

func (i *InstanceWrapper) Disable() error {
	i.enabled = false
	err := i.instance.Disable()
	if err != nil {
		i.enabled = true
		return err
	}
	return nil
}

// Manager is an encapsulating layer for plugins and manages all plugins and its instances.
type Manager struct {
	mutex     *sync.RWMutex
	instances map[uint]*InstanceWrapper
	plugins   map[string]compat.Plugin
	messages  chan MessageWithUserID
	db        *database.GormDatabase
	mux       *gin.RouterGroup
}

// NewManager created a Manager from configurations.
func NewManager(db *database.GormDatabase, directory string, mux *gin.RouterGroup, notifier Notifier) (*Manager, error) {
	manager := &Manager{
		mutex:     &sync.RWMutex{},
		instances: map[uint]*InstanceWrapper{},
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
				Priority:      *message.Message.Priority,
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
		if err := manager.initializeForUser(manager.db, *user); err != nil {
			return nil, err
		}
	}

	return manager, nil
}

// ErrAlreadyEnabledOrDisabled is returned on SetPluginEnabled call when a plugin is already enabled or disabled.
var ErrAlreadyEnabledOrDisabled = errors.New("config is already enabled/disabled")

// SetPluginEnabled sets the plugins enabled state.
func (m *Manager) SetPluginEnabled(pluginID uint, enabled bool) error {
	instanceWrapper, err := m.Instance(pluginID)
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
		err = instanceWrapper.Enable()
	} else {
		err = instanceWrapper.Disable()
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
	log.Warn().Str("module_path", modulePath).Msg("Could not get plugin info")
	return compat.Info{
		Name:        "UNKNOWN",
		ModulePath:  modulePath,
		Description: "Oops something went wrong",
	}
}

// Instance returns an instance with the given ID.
func (m *Manager) Instance(pluginID uint) (*InstanceWrapper, error) {
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
func (m *Manager) RemoveUser(tx *database.GormDatabase, userID uint) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	maps.DeleteFunc(m.instances, func(id uint, instance *InstanceWrapper) (delete bool) {
		delete = instance.userID == userID
		if delete {
			// ignore errors and force delete to prevent
			// leftover instances for deleted users
			err := instance.Disable()
			if err != nil {
				log.Warn().Err(err).Uint("user_id", userID).Uint("plugin_id", id).Msg("Plugin disable failed")
			}
		}
		return delete
	})

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

	pluginFiles, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("error while reading directory %s", err)
	}
	for _, f := range pluginFiles {
		if f.IsDir() {
			continue
		}

		name := f.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}

		pluginPath := filepath.Join(directory, "./", name)

		log.Info().Str("path", pluginPath).Msg("Loading plugin")
		pRaw, err := plugin.Open(pluginPath)
		if err != nil {
			return pluginFileLoadError{name, err}
		}
		compatPlugin, err := compat.Wrap(pRaw)
		if err != nil {
			return pluginFileLoadError{name, err}
		}
		if err := m.LoadPlugin(compatPlugin); err != nil {
			return pluginFileLoadError{name, err}
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
func (m *Manager) InitializeForUserID(tx *database.GormDatabase, userID uint) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	user, err := tx.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user != nil {
		return m.initializeForUser(tx, *user)
	}
	return fmt.Errorf("user with id %d not found", userID)
}

func (m *Manager) initializeForUser(tx *database.GormDatabase, user model.User) error {
	userCtx := compat.UserContext{
		ID:    user.ID,
		Name:  user.Name,
		Admin: user.Admin,
	}

	for _, p := range m.plugins {
		if err := m.initializeSingleUserPlugin(tx, userCtx, p); err != nil {
			return err
		}
	}

	apps, err := tx.GetApplicationsByUser(user.ID)
	if err != nil {
		return err
	}
	for _, app := range apps {
		conf, err := tx.GetPluginConfByApplicationID(app.ID)
		if err != nil {
			return err
		}
		if conf != nil {
			_, compatExist := m.plugins[conf.ModulePath]
			app.Internal = compatExist
		} else {
			app.Internal = false
		}
		tx.UpdateApplication(app)
	}

	return nil
}

func (m *Manager) initializeSingleUserPlugin(tx *database.GormDatabase, userCtx compat.UserContext, p compat.Plugin) error {
	info := p.PluginInfo()
	instance := p.NewPluginInstance(userCtx)
	userID := userCtx.ID

	pluginConf, err := tx.GetPluginConfByUserAndPath(userID, info.ModulePath)
	if err != nil {
		return err
	}

	if pluginConf == nil {
		var err error
		pluginConf, err = m.createPluginConf(tx, instance, info, userID)
		if err != nil {
			return err
		}
	}

	instanceWrapper := &InstanceWrapper{
		userID:   userID,
		instance: instance,
	}

	m.instances[pluginConf.ID] = instanceWrapper

	if compat.HasSupport(instance, compat.Messenger) {
		if pluginConf.ApplicationID == 0 {
			// The Messenger capability was added after this plugin was first
			// initialized for the user, so no internal application exists yet.
			// Create one now, otherwise messages would be stored with
			// application_id = 0 and become orphaned (not shown, not deletable).
			app, err := m.createInternalApplication(tx, info, userID)
			if err != nil {
				return err
			}
			pluginConf.ApplicationID = app.ID
			if err := tx.UpdatePluginConf(pluginConf); err != nil {
				return err
			}
		}
		instance.SetMessageHandler(redirectToChannel{
			ApplicationID: pluginConf.ApplicationID,
			UserID:        pluginConf.UserID,
			Messages:      m.messages,
		})
	}
	if compat.HasSupport(instance, compat.Storager) {
		instance.SetStorageHandler(dbStorageHandler{pluginConf.ID, tx})
	}
	if compat.HasSupport(instance, compat.Configurer) {
		m.initializeConfigurerForSingleUserPlugin(tx, instance, pluginConf)
	}
	if compat.HasSupport(instance, compat.Webhooker) {
		id := pluginConf.ID
		g := m.mux.Group(pluginConf.Token+"/", requirePluginEnabled(id, tx))
		instance.RegisterWebhook(strings.Replace(g.BasePath(), ":id", strconv.Itoa(int(id)), 1), g)
	}
	if pluginConf.Enabled {
		err := instanceWrapper.Enable()
		if err != nil {
			// Single user plugin cannot be enabled
			// Don't panic, disable for now and wait for user to update config
			log.Warn().Err(err).Str("user", userCtx.Name).Msg("Plugin initialize failed, disabling now")
			pluginConf.Enabled = false
			if err = tx.UpdatePluginConf(pluginConf); err != nil {
				log.Warn().Err(err).Uint("plugin_id", pluginConf.ID).Msg("Plugin enable failed, disabling now")
			}
		}
	}
	return nil
}

func (m *Manager) initializeConfigurerForSingleUserPlugin(tx *database.GormDatabase, instance compat.PluginInstance, pluginConf *model.PluginConf) {
	if len(pluginConf.Config) == 0 {
		// The Configurer is newly implemented
		// Use the default config
		pluginConf.Config, _ = yaml.Marshal(instance.DefaultConfig())
		tx.UpdatePluginConf(pluginConf)
	}
	c := instance.DefaultConfig()
	if yaml.Unmarshal(pluginConf.Config, c) != nil || instance.ValidateAndSetConfig(c) != nil {
		pluginConf.Enabled = false

		log.Warn().
			Str("module_path", pluginConf.ModulePath).
			Uint("user_id", pluginConf.UserID).
			Msg("Plugin failed to initialize because it rejected the current config. It might be outdated. A default config is used and the user would need to enable it again.")
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

		tx.UpdatePluginConf(pluginConf)
		instance.ValidateAndSetConfig(instance.DefaultConfig())
	}
}

func (m *Manager) createPluginConf(tx *database.GormDatabase, instance compat.PluginInstance, info compat.Info, userID uint) (*model.PluginConf, error) {
	pluginConf := &model.PluginConf{
		UserID:     userID,
		ModulePath: info.ModulePath,
		Token:      auth.GeneratePluginToken(),
	}
	if compat.HasSupport(instance, compat.Configurer) {
		pluginConf.Config, _ = yaml.Marshal(instance.DefaultConfig())
	}
	if compat.HasSupport(instance, compat.Messenger) {
		app, err := m.createInternalApplication(tx, info, userID)
		if err != nil {
			return nil, err
		}
		pluginConf.ApplicationID = app.ID
	}
	if err := tx.CreatePluginConf(pluginConf); err != nil {
		return nil, err
	}
	return pluginConf, nil
}

// createInternalApplication creates the auto generated internal application a
// Messenger plugin uses to publish its messages.
func (m *Manager) createInternalApplication(tx *database.GormDatabase, info compat.Info, userID uint) (*model.Application, error) {
	tokenPublic, _ := auth.GenerateApplicationToken()
	app := &model.Application{
		Token:       tokenPublic,
		Name:        info.String(),
		UserID:      userID,
		Internal:    true,
		Description: fmt.Sprintf("auto generated application for %s", info.ModulePath),
	}
	if err := tx.CreateApplication(app); err != nil {
		return nil, err
	}
	return app, nil
}
