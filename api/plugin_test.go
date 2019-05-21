package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/go-yaml/yaml"

	"github.com/gotify/server/mode"
	"github.com/gotify/server/model"
	"github.com/gotify/server/plugin"
	"github.com/gotify/server/plugin/compat"
	"github.com/gotify/server/plugin/testing/mock"
	"github.com/gotify/server/test"
	"github.com/gotify/server/test/testdb"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestPluginSuite(t *testing.T) {
	suite.Run(t, new(PluginSuite))
}

type PluginSuite struct {
	suite.Suite
	db       *testdb.Database
	a        *PluginAPI
	ctx      *gin.Context
	recorder *httptest.ResponseRecorder
	manager  *plugin.Manager
	notified bool
}

func (s *PluginSuite) BeforeTest(suiteName, testName string) {
	mode.Set(mode.TestDev)
	s.db = testdb.NewDB(s.T())
	s.resetRecorder()
	manager, err := plugin.NewManager(s.db, "", nil, s)
	assert.Nil(s.T(), err)
	s.manager = manager
	withURL(s.ctx, "http", "example.com")
	s.a = &PluginAPI{DB: s.db, Manager: manager, Notifier: s}

	mockPluginCompat := new(mock.Plugin)
	assert.Nil(s.T(), s.manager.LoadPlugin(mockPluginCompat))

	s.db.User(1)
	assert.Nil(s.T(), s.manager.InitializeForUserID(1))
	s.db.User(2)
	assert.Nil(s.T(), s.manager.InitializeForUserID(2))

	s.db.CreatePluginConf(&model.PluginConf{
		UserID:     1,
		ModulePath: "github.com/gotify/server/plugin/example/removed",
		Token:      "P1234",
		Enabled:    false,
	})
}

func (s *PluginSuite) getDanglingConf(uid uint) *model.PluginConf {
	conf, err := s.db.GetPluginConfByUserAndPath(uid, "github.com/gotify/server/plugin/example/removed")
	assert.NoError(s.T(), err)
	return conf
}

func (s *PluginSuite) resetRecorder() {
	s.recorder = httptest.NewRecorder()
	s.ctx, _ = gin.CreateTestContext(s.recorder)
}

func (s *PluginSuite) AfterTest(suiteName, testName string) {
	s.db.Close()
}

func (s *PluginSuite) Notify(userID uint, msg *model.MessageExternal) {
	s.notified = true
}

func (s *PluginSuite) Test_GetPlugins() {
	test.WithUser(s.ctx, 1)

	s.ctx.Request = httptest.NewRequest("GET", "/plugin", nil)
	s.a.GetPlugins(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)

	pluginConfs := make([]model.PluginConfExternal, 0)
	assert.Nil(s.T(), json.Unmarshal(s.recorder.Body.Bytes(), &pluginConfs))

	assert.Equal(s.T(), mock.Name, pluginConfs[0].Name)
	assert.Equal(s.T(), mock.ModulePath, pluginConfs[0].ModulePath)

	assert.False(s.T(), pluginConfs[0].Enabled, "Plugins should be disabled by default")
}

func (s *PluginSuite) Test_EnableDisablePlugin() {

	{
		test.WithUser(s.ctx, 1)

		s.ctx.Request = httptest.NewRequest("POST", "/plugin/1/enable", nil)
		s.ctx.Params = gin.Params{{Key: "id", Value: "1"}}
		s.a.EnablePlugin(s.ctx)

		assert.Equal(s.T(), 200, s.recorder.Code)

		if pluginConf, err := s.db.GetPluginConfByUserAndPath(1, mock.ModulePath); assert.NoError(s.T(), err) {
			assert.True(s.T(), pluginConf.Enabled)
		}
		s.resetRecorder()
	}

	{
		test.WithUser(s.ctx, 1)

		s.ctx.Request = httptest.NewRequest("POST", "/plugin/1/enable", nil)
		s.ctx.Params = gin.Params{{Key: "id", Value: "1"}}
		s.a.EnablePlugin(s.ctx)

		assert.Equal(s.T(), 400, s.recorder.Code)

		if pluginConf, err := s.db.GetPluginConfByUserAndPath(1, mock.ModulePath); assert.NoError(s.T(), err) {
			assert.True(s.T(), pluginConf.Enabled)
		}
		s.resetRecorder()
	}

	{
		test.WithUser(s.ctx, 1)

		s.ctx.Request = httptest.NewRequest("POST", "/plugin/1/disable", nil)
		s.ctx.Params = gin.Params{{Key: "id", Value: "1"}}
		s.a.DisablePlugin(s.ctx)

		assert.Equal(s.T(), 200, s.recorder.Code)

		if pluginConf, err := s.db.GetPluginConfByUserAndPath(1, mock.ModulePath); assert.NoError(s.T(), err) {
			assert.False(s.T(), pluginConf.Enabled)
		}
		s.resetRecorder()
	}

	{
		test.WithUser(s.ctx, 1)

		s.ctx.Request = httptest.NewRequest("POST", "/plugin/1/disable", nil)
		s.ctx.Params = gin.Params{{Key: "id", Value: "1"}}
		s.a.DisablePlugin(s.ctx)

		assert.Equal(s.T(), 400, s.recorder.Code)

		if pluginConf, err := s.db.GetPluginConfByUserAndPath(1, mock.ModulePath); assert.NoError(s.T(), err) {
			assert.False(s.T(), pluginConf.Enabled)
		}
		s.resetRecorder()
	}

}

func (s *PluginSuite) Test_EnableDisablePlugin_EnableReturnsError_expect500() {
	s.db.User(16)
	assert.Nil(s.T(), s.manager.InitializeForUserID(16))
	mock.ReturnErrorOnEnableForUser(16, errors.New("test error"))
	conf, err := s.db.GetPluginConfByUserAndPath(16, mock.ModulePath)
	assert.NoError(s.T(), err)

	{
		test.WithUser(s.ctx, 16)
		s.ctx.Request = httptest.NewRequest("POST", fmt.Sprintf("/plugin/%d/enable", conf.ID), nil)
		s.ctx.Params = gin.Params{{Key: "id", Value: fmt.Sprint(conf.ID)}}
		s.a.EnablePlugin(s.ctx)

		assert.Equal(s.T(), 500, s.recorder.Code)

		if pluginConf, err := s.db.GetPluginConfByUserAndPath(1, mock.ModulePath); assert.NoError(s.T(), err) {
			assert.False(s.T(), pluginConf.Enabled)
		}
		s.resetRecorder()
	}
}

func (s *PluginSuite) Test_EnableDisablePlugin_DisableReturnsError_expect500() {
	s.db.User(17)
	assert.Nil(s.T(), s.manager.InitializeForUserID(17))
	mock.ReturnErrorOnDisableForUser(17, errors.New("test error"))
	conf, err := s.db.GetPluginConfByUserAndPath(17, mock.ModulePath)
	assert.NoError(s.T(), err)
	s.manager.SetPluginEnabled(conf.ID, true)

	{
		test.WithUser(s.ctx, 17)
		s.ctx.Request = httptest.NewRequest("POST", fmt.Sprintf("/plugin/%d/disable", conf.ID), nil)
		s.ctx.Params = gin.Params{{Key: "id", Value: fmt.Sprint(conf.ID)}}
		s.a.DisablePlugin(s.ctx)

		assert.Equal(s.T(), 500, s.recorder.Code)

		if pluginConf, err := s.db.GetPluginConfByUserAndPath(1, mock.ModulePath); assert.NoError(s.T(), err) {
			assert.False(s.T(), pluginConf.Enabled)
		}
		s.resetRecorder()
	}
}

func (s *PluginSuite) Test_EnableDisablePlugin_incorrectUser_expectNotFound() {
	{
		test.WithUser(s.ctx, 2)

		s.ctx.Request = httptest.NewRequest("POST", "/plugin/1/enable", nil)
		s.ctx.Params = gin.Params{{Key: "id", Value: "1"}}
		s.a.EnablePlugin(s.ctx)

		assert.Equal(s.T(), 404, s.recorder.Code)

		if pluginConf, err := s.db.GetPluginConfByUserAndPath(1, mock.ModulePath); assert.NoError(s.T(), err) {
			assert.False(s.T(), pluginConf.Enabled)
		}
		s.resetRecorder()
	}

	{
		test.WithUser(s.ctx, 2)

		s.ctx.Request = httptest.NewRequest("POST", "/plugin/1/disable", nil)
		s.ctx.Params = gin.Params{{Key: "id", Value: "1"}}
		s.a.DisablePlugin(s.ctx)

		assert.Equal(s.T(), 404, s.recorder.Code)

		if pluginConf, err := s.db.GetPluginConfByUserAndPath(1, mock.ModulePath); assert.NoError(s.T(), err) {
			assert.False(s.T(), pluginConf.Enabled)
		}
		s.resetRecorder()
	}

}

func (s *PluginSuite) Test_EnableDisablePlugin_nonExistPlugin_expectNotFound() {
	{
		test.WithUser(s.ctx, 1)

		s.ctx.Request = httptest.NewRequest("POST", "/plugin/99/enable", nil)
		s.ctx.Params = gin.Params{{Key: "id", Value: "99"}}
		s.a.EnablePlugin(s.ctx)

		assert.Equal(s.T(), 404, s.recorder.Code)
		s.resetRecorder()
	}

	{
		test.WithUser(s.ctx, 1)

		s.ctx.Request = httptest.NewRequest("POST", "/plugin/99/disable", nil)
		s.ctx.Params = gin.Params{{Key: "id", Value: "99"}}
		s.a.DisablePlugin(s.ctx)

		assert.Equal(s.T(), 404, s.recorder.Code)
		s.resetRecorder()
	}

}

func (s *PluginSuite) Test_EnableDisablePlugin_danglingConf_expectNotFound() {
	conf := s.getDanglingConf(1)

	{
		test.WithUser(s.ctx, 1)

		s.ctx.Request = httptest.NewRequest("POST", fmt.Sprintf("/plugin/%d/enable", conf.ID), nil)
		s.ctx.Params = gin.Params{{Key: "id", Value: fmt.Sprint(conf.ID)}}
		s.a.EnablePlugin(s.ctx)

		assert.Equal(s.T(), 404, s.recorder.Code)
		s.resetRecorder()
	}

	{
		test.WithUser(s.ctx, 1)

		s.ctx.Request = httptest.NewRequest("POST", fmt.Sprintf("/plugin/%d/disable", conf.ID), nil)
		s.ctx.Params = gin.Params{{Key: "id", Value: fmt.Sprint(conf.ID)}}
		s.a.DisablePlugin(s.ctx)

		assert.Equal(s.T(), 404, s.recorder.Code)
		s.resetRecorder()
	}
}

func (s *PluginSuite) Test_GetDisplay() {
	conf, err := s.db.GetPluginConfByUserAndPath(1, mock.ModulePath)
	assert.NoError(s.T(), err)
	inst, err := s.manager.Instance(conf.ID)
	assert.Nil(s.T(), err)
	mockInst := inst.(*mock.PluginInstance)

	mockInst.DisplayString = "test string"

	{
		test.WithUser(s.ctx, 1)

		s.ctx.Request = httptest.NewRequest("GET", fmt.Sprintf("/plugin/%d/display", conf.ID), nil)
		s.ctx.Params = gin.Params{{Key: "id", Value: fmt.Sprint(conf.ID)}}
		s.a.GetDisplay(s.ctx)

		assert.Equal(s.T(), 200, s.recorder.Code)
		test.JSONEquals(s.T(), mockInst.DisplayString, s.recorder.Body.String())
	}
}

func (s *PluginSuite) Test_GetDisplay_NotImplemented_expectEmptyString() {
	conf, err := s.db.GetPluginConfByUserAndPath(1, mock.ModulePath)
	assert.NoError(s.T(), err)
	inst, err := s.manager.Instance(conf.ID)
	assert.Nil(s.T(), err)
	mockInst := inst.(*mock.PluginInstance)

	mockInst.SetCapability(compat.Displayer, false)
	defer mockInst.SetCapability(compat.Displayer, true)

	{
		test.WithUser(s.ctx, 1)

		s.ctx.Request = httptest.NewRequest("GET", fmt.Sprintf("/plugin/%d/display", conf.ID), nil)
		s.ctx.Params = gin.Params{{Key: "id", Value: fmt.Sprint(conf.ID)}}
		s.a.GetDisplay(s.ctx)

		assert.Equal(s.T(), 200, s.recorder.Code)
		test.JSONEquals(s.T(), "", s.recorder.Body.String())
	}
}

func (s *PluginSuite) Test_GetDisplay_incorrectUser_expectNotFound() {
	conf, err := s.db.GetPluginConfByUserAndPath(1, mock.ModulePath)
	assert.NoError(s.T(), err)
	inst, err := s.manager.Instance(conf.ID)
	assert.Nil(s.T(), err)
	mockInst := inst.(*mock.PluginInstance)

	mockInst.DisplayString = "test string"

	{
		test.WithUser(s.ctx, 2)

		s.ctx.Request = httptest.NewRequest("GET", fmt.Sprintf("/plugin/%d/display", conf.ID), nil)
		s.ctx.Params = gin.Params{{Key: "id", Value: fmt.Sprint(conf.ID)}}
		s.a.GetDisplay(s.ctx)

		assert.Equal(s.T(), 404, s.recorder.Code)
	}
}

func (s *PluginSuite) Test_GetDisplay_danglingConf_expectNotFound() {
	conf := s.getDanglingConf(1)

	{
		test.WithUser(s.ctx, 1)

		s.ctx.Request = httptest.NewRequest("GET", fmt.Sprintf("/plugin/%d/display", conf.ID), nil)
		s.ctx.Params = gin.Params{{Key: "id", Value: fmt.Sprint(conf.ID)}}
		s.a.GetDisplay(s.ctx)

		assert.Equal(s.T(), 404, s.recorder.Code)
	}
}

func (s *PluginSuite) Test_GetDisplay_nonExistPlugin_expectNotFound() {
	{
		test.WithUser(s.ctx, 1)

		s.ctx.Request = httptest.NewRequest("GET", "/plugin/99/display", nil)
		s.ctx.Params = gin.Params{{Key: "id", Value: "99"}}
		s.a.GetDisplay(s.ctx)

		assert.Equal(s.T(), 404, s.recorder.Code)
	}
}

func (s *PluginSuite) Test_GetConfig() {
	conf, err := s.db.GetPluginConfByUserAndPath(1, mock.ModulePath)
	assert.NoError(s.T(), err)
	inst, err := s.manager.Instance(conf.ID)
	assert.Nil(s.T(), err)
	mockInst := inst.(*mock.PluginInstance)

	assert.Equal(s.T(), mockInst.DefaultConfig(), mockInst.Config, "Initial config should be default config")
	{
		test.WithUser(s.ctx, 1)

		s.ctx.Request = httptest.NewRequest("GET", fmt.Sprintf("/plugin/%d/config", conf.ID), nil)
		s.ctx.Params = gin.Params{{Key: "id", Value: fmt.Sprint(conf.ID)}}
		s.a.GetConfig(s.ctx)

		assert.Equal(s.T(), 200, s.recorder.Code)
		returnedConfig := new(mock.PluginConfig)
		assert.Nil(s.T(), yaml.Unmarshal(s.recorder.Body.Bytes(), returnedConfig))
		assert.Equal(s.T(), mockInst.Config, returnedConfig)
	}
}

func (s *PluginSuite) Test_GetConfg_notImplemeted_expect400() {
	conf, err := s.db.GetPluginConfByUserAndPath(1, mock.ModulePath)
	assert.NoError(s.T(), err)
	inst, err := s.manager.Instance(conf.ID)
	assert.Nil(s.T(), err)
	mockInst := inst.(*mock.PluginInstance)

	mockInst.SetCapability(compat.Configurer, false)
	defer mockInst.SetCapability(compat.Configurer, true)

	{
		test.WithUser(s.ctx, 1)

		s.ctx.Request = httptest.NewRequest("GET", fmt.Sprintf("/plugin/%d/config", conf.ID), nil)
		s.ctx.Params = gin.Params{{Key: "id", Value: fmt.Sprint(conf.ID)}}
		s.a.GetConfig(s.ctx)

		assert.Equal(s.T(), 400, s.recorder.Code)
	}
}

func (s *PluginSuite) Test_GetConfig_incorrectUser_expectNotFound() {
	conf, err := s.db.GetPluginConfByUserAndPath(1, mock.ModulePath)
	assert.NoError(s.T(), err)

	{
		test.WithUser(s.ctx, 2)

		s.ctx.Request = httptest.NewRequest("GET", fmt.Sprintf("/plugin/%d/config", conf.ID), nil)
		s.ctx.Params = gin.Params{{Key: "id", Value: fmt.Sprint(conf.ID)}}
		s.a.GetConfig(s.ctx)

		assert.Equal(s.T(), 404, s.recorder.Code)
	}
}

func (s *PluginSuite) Test_GetConfig_danglingConf_expectNotFound() {
	conf := s.getDanglingConf(1)

	{
		test.WithUser(s.ctx, 1)

		s.ctx.Request = httptest.NewRequest("GET", fmt.Sprintf("/plugin/%d/config", conf.ID), nil)
		s.ctx.Params = gin.Params{{Key: "id", Value: fmt.Sprint(conf.ID)}}
		s.a.GetConfig(s.ctx)

		assert.Equal(s.T(), 404, s.recorder.Code)
	}
}

func (s *PluginSuite) Test_GetConfig_nonExistPlugin_expectNotFound() {
	{
		test.WithUser(s.ctx, 1)

		s.ctx.Request = httptest.NewRequest("GET", "/plugin/99/config", nil)
		s.ctx.Params = gin.Params{{Key: "id", Value: "99"}}
		s.a.GetConfig(s.ctx)

		assert.Equal(s.T(), 404, s.recorder.Code)
	}
}

func (s *PluginSuite) Test_UpdateConfig() {
	conf, err := s.db.GetPluginConfByUserAndPath(1, mock.ModulePath)
	assert.NoError(s.T(), err)
	inst, err := s.manager.Instance(conf.ID)
	assert.Nil(s.T(), err)
	mockInst := inst.(*mock.PluginInstance)

	newConfig := &mock.PluginConfig{
		TestKey: "test__new__config",
	}
	newConfigYAML, err := yaml.Marshal(newConfig)
	assert.Nil(s.T(), err)

	{
		test.WithUser(s.ctx, 1)

		s.ctx.Request = httptest.NewRequest("POST", fmt.Sprintf("/plugin/%d/config", conf.ID), bytes.NewReader(newConfigYAML))
		s.ctx.Header("Content-Type", "application/x-yaml")
		s.ctx.Params = gin.Params{{Key: "id", Value: fmt.Sprint(conf.ID)}}
		s.a.UpdateConfig(s.ctx)

		assert.Equal(s.T(), 200, s.recorder.Code)
		assert.Equal(s.T(), newConfig, mockInst.Config, "config should be received by plugin")

		var pluginFromDBBytes []byte
		if pluginConf, err := s.db.GetPluginConfByID(conf.ID); assert.NoError(s.T(), err) {
			pluginFromDBBytes = pluginConf.Config
		}
		pluginFromDB := new(mock.PluginConfig)
		err := yaml.Unmarshal(pluginFromDBBytes, pluginFromDB)
		assert.Nil(s.T(), err)
		assert.Equal(s.T(), newConfig, pluginFromDB, "config should be updated in database")
	}
}

func (s *PluginSuite) Test_UpdateConfig_invalidConfig_expect400() {
	conf, err := s.db.GetPluginConfByUserAndPath(1, mock.ModulePath)
	assert.NoError(s.T(), err)
	inst, err := s.manager.Instance(conf.ID)
	assert.Nil(s.T(), err)
	mockInst := inst.(*mock.PluginInstance)
	origConfig := mockInst.Config

	newConfig := &mock.PluginConfig{
		TestKey:    "test__new__config__invalid",
		IsNotValid: true,
	}
	newConfigYAML, err := yaml.Marshal(newConfig)
	assert.Nil(s.T(), err)

	{
		test.WithUser(s.ctx, 1)

		s.ctx.Request = httptest.NewRequest("POST", fmt.Sprintf("/plugin/%d/config", conf.ID), bytes.NewReader(newConfigYAML))
		s.ctx.Header("Content-Type", "application/x-yaml")
		s.ctx.Params = gin.Params{{Key: "id", Value: fmt.Sprint(conf.ID)}}
		s.a.UpdateConfig(s.ctx)

		assert.Equal(s.T(), 400, s.recorder.Code)
		assert.Equal(s.T(), origConfig, mockInst.Config, "config should not be received by plugin")

		var pluginFromDBBytes []byte
		if pluginConf, err := s.db.GetPluginConfByID(conf.ID); assert.NoError(s.T(), err) {
			pluginFromDBBytes = pluginConf.Config
		}
		pluginFromDB := new(mock.PluginConfig)
		err := yaml.Unmarshal(pluginFromDBBytes, pluginFromDB)
		assert.Nil(s.T(), err)
		assert.Equal(s.T(), origConfig, pluginFromDB, "config should not be updated in database")
	}
}

func (s *PluginSuite) Test_UpdateConfig_malformedYAML_expect400() {
	conf, err := s.db.GetPluginConfByUserAndPath(1, mock.ModulePath)
	assert.NoError(s.T(), err)
	inst, err := s.manager.Instance(conf.ID)
	assert.Nil(s.T(), err)
	mockInst := inst.(*mock.PluginInstance)
	origConfig := mockInst.Config

	newConfigYAML := []byte(`--- "rg e""`)

	{
		test.WithUser(s.ctx, 1)

		s.ctx.Request = httptest.NewRequest("POST", fmt.Sprintf("/plugin/%d/config", conf.ID), bytes.NewReader(newConfigYAML))
		s.ctx.Header("Content-Type", "application/x-yaml")
		s.ctx.Params = gin.Params{{Key: "id", Value: fmt.Sprint(conf.ID)}}
		s.a.UpdateConfig(s.ctx)

		assert.Equal(s.T(), 400, s.recorder.Code)
		assert.Equal(s.T(), origConfig, mockInst.Config, "config should not be received by plugin")

		var pluginFromDBBytes []byte
		if pluginConf, err := s.db.GetPluginConfByID(conf.ID); assert.NoError(s.T(), err) {
			pluginFromDBBytes = pluginConf.Config
		}
		pluginFromDB := new(mock.PluginConfig)
		err := yaml.Unmarshal(pluginFromDBBytes, pluginFromDB)
		assert.Nil(s.T(), err)
		assert.Equal(s.T(), origConfig, pluginFromDB, "config should not be updated in database")
	}
}

func (s *PluginSuite) Test_UpdateConfig_ioError_expect500() {
	conf, err := s.db.GetPluginConfByUserAndPath(1, mock.ModulePath)
	assert.NoError(s.T(), err)
	inst, err := s.manager.Instance(conf.ID)
	assert.Nil(s.T(), err)
	mockInst := inst.(*mock.PluginInstance)
	origConfig := mockInst.Config

	{
		test.WithUser(s.ctx, 1)

		s.ctx.Request = httptest.NewRequest("POST", fmt.Sprintf("/plugin/%d/config", conf.ID), test.UnreadableReader())
		s.ctx.Header("Content-Type", "application/x-yaml")
		s.ctx.Params = gin.Params{{Key: "id", Value: fmt.Sprint(conf.ID)}}
		s.a.UpdateConfig(s.ctx)

		assert.Equal(s.T(), 500, s.recorder.Code)
		assert.Equal(s.T(), origConfig, mockInst.Config, "config should not be received by plugin")

		var pluginFromDBBytes []byte
		if pluginConf, err := s.db.GetPluginConfByID(conf.ID); assert.NoError(s.T(), err) {
			pluginFromDBBytes = pluginConf.Config
		}
		pluginFromDB := new(mock.PluginConfig)
		err := yaml.Unmarshal(pluginFromDBBytes, pluginFromDB)
		assert.Nil(s.T(), err)
		assert.Equal(s.T(), origConfig, pluginFromDB, "config should not be updated in database")
	}
}

func (s *PluginSuite) Test_UpdateConfig_notImplemented_expect400() {
	conf, err := s.db.GetPluginConfByUserAndPath(1, mock.ModulePath)
	assert.NoError(s.T(), err)
	inst, err := s.manager.Instance(conf.ID)
	assert.Nil(s.T(), err)
	mockInst := inst.(*mock.PluginInstance)

	newConfig := &mock.PluginConfig{
		TestKey: "test__new__config",
	}
	newConfigYAML, err := yaml.Marshal(newConfig)
	assert.Nil(s.T(), err)

	mockInst.SetCapability(compat.Configurer, false)
	defer mockInst.SetCapability(compat.Configurer, true)

	{
		test.WithUser(s.ctx, 1)

		s.ctx.Request = httptest.NewRequest("POST", fmt.Sprintf("/plugin/%d/config", conf.ID), bytes.NewReader(newConfigYAML))
		s.ctx.Header("Content-Type", "application/x-yaml")
		s.ctx.Params = gin.Params{{Key: "id", Value: fmt.Sprint(conf.ID)}}
		s.a.UpdateConfig(s.ctx)

		assert.Equal(s.T(), 400, s.recorder.Code)
	}
}

func (s *PluginSuite) Test_UpdateConfig_incorrectUser_expectNotFound() {
	conf, err := s.db.GetPluginConfByUserAndPath(1, mock.ModulePath)
	assert.NoError(s.T(), err)
	inst, err := s.manager.Instance(conf.ID)
	assert.Nil(s.T(), err)
	mockInst := inst.(*mock.PluginInstance)
	origConfig := mockInst.Config

	newConfig := &mock.PluginConfig{
		TestKey: "test__new__config",
	}
	newConfigYAML, err := yaml.Marshal(newConfig)
	assert.Nil(s.T(), err)

	{
		test.WithUser(s.ctx, 2)

		s.ctx.Request = httptest.NewRequest("POST", fmt.Sprintf("/plugin/%d/config", conf.ID), bytes.NewReader(newConfigYAML))
		s.ctx.Header("Content-Type", "application/x-yaml")
		s.ctx.Params = gin.Params{{Key: "id", Value: fmt.Sprint(conf.ID)}}
		s.a.UpdateConfig(s.ctx)

		assert.Equal(s.T(), 404, s.recorder.Code)
		assert.Equal(s.T(), origConfig, mockInst.Config, "config should not be received by plugin")

		var pluginFromDBBytes []byte
		if pluginConf, err := s.db.GetPluginConfByID(conf.ID); assert.NoError(s.T(), err) {
			pluginFromDBBytes = pluginConf.Config
		}
		pluginFromDB := new(mock.PluginConfig)
		err := yaml.Unmarshal(pluginFromDBBytes, pluginFromDB)
		assert.Nil(s.T(), err)
		assert.Equal(s.T(), origConfig, pluginFromDB, "config should not be updated in database")
	}
}

func (s *PluginSuite) Test_UpdateConfig_danglingConf_expectNotFound() {
	conf := s.getDanglingConf(1)

	newConfig := &mock.PluginConfig{
		TestKey: "test__new__config",
	}
	newConfigYAML, err := yaml.Marshal(newConfig)
	assert.Nil(s.T(), err)

	{
		test.WithUser(s.ctx, 1)

		s.ctx.Request = httptest.NewRequest("POST", fmt.Sprintf("/plugin/%d/config", conf.ID), bytes.NewReader(newConfigYAML))
		s.ctx.Header("Content-Type", "application/x-yaml")
		s.ctx.Params = gin.Params{{Key: "id", Value: fmt.Sprint(conf.ID)}}
		s.a.UpdateConfig(s.ctx)

		assert.Equal(s.T(), 404, s.recorder.Code)
	}
}

func (s *PluginSuite) Test_UpdateConfig_nonExistPlugin_expectNotFound() {
	newConfig := &mock.PluginConfig{
		TestKey: "test__new__config",
	}
	newConfigYAML, err := yaml.Marshal(newConfig)
	assert.Nil(s.T(), err)

	{
		test.WithUser(s.ctx, 1)

		s.ctx.Request = httptest.NewRequest("POST", "/plugin/99/config", bytes.NewReader(newConfigYAML))
		s.ctx.Header("Content-Type", "application/x-yaml")
		s.ctx.Params = gin.Params{{Key: "id", Value: "99"}}
		s.a.UpdateConfig(s.ctx)

		assert.Equal(s.T(), 404, s.recorder.Code)
	}
}
