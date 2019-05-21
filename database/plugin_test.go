package database

import (
	"github.com/gotify/server/model"
	"github.com/stretchr/testify/assert"
)

func (s *DatabaseSuite) TestPluginConf() {
	plugin := model.PluginConf{
		ModulePath:    "github.com/gotify/example-plugin",
		Token:         "Pabc",
		UserID:        1,
		Enabled:       true,
		Config:        nil,
		ApplicationID: 2,
	}

	assert.Nil(s.T(), s.db.CreatePluginConf(&plugin))

	assert.Equal(s.T(), uint(1), plugin.ID)
	if pluginConf, err := s.db.GetPluginConfByUserAndPath(1, "github.com/gotify/example-plugin"); assert.NoError(s.T(), err) {
		assert.Equal(s.T(), "Pabc", pluginConf.Token)
	}
	if pluginConf, err := s.db.GetPluginConfByToken("Pabc"); assert.NoError(s.T(), err) {
		assert.Equal(s.T(), true, pluginConf.Enabled)
	}
	if pluginConf, err := s.db.GetPluginConfByApplicationID(2); assert.NoError(s.T(), err) {
		assert.Equal(s.T(), "Pabc", pluginConf.Token)
	}
	if pluginConf, err := s.db.GetPluginConfByID(1); assert.NoError(s.T(), err) {
		assert.Equal(s.T(), "github.com/gotify/example-plugin", pluginConf.ModulePath)
	}

	if pluginConf, err := s.db.GetPluginConfByToken("Pnotexist"); assert.NoError(s.T(), err) {
		assert.Nil(s.T(), pluginConf)
	}
	if pluginConf, err := s.db.GetPluginConfByID(12); assert.NoError(s.T(), err) {
		assert.Nil(s.T(), pluginConf)
	}
	if pluginConf, err := s.db.GetPluginConfByUserAndPath(1, "not/exist"); assert.NoError(s.T(), err) {
		assert.Nil(s.T(), pluginConf)
	}
	if pluginConf, err := s.db.GetPluginConfByApplicationID(99); assert.NoError(s.T(), err) {
		assert.Nil(s.T(), pluginConf)
	}

	if pluginConfs, err := s.db.GetPluginConfByUser(1); assert.NoError(s.T(), err) {
		assert.Len(s.T(), pluginConfs, 1)
	}
	if pluginConfs, err := s.db.GetPluginConfByUser(0); assert.NoError(s.T(), err) {
		assert.Len(s.T(), pluginConfs, 0)
	}

	testConf := `{"test_config_key":"hello"}`
	plugin.Enabled = false
	plugin.Config = []byte(testConf)
	assert.Nil(s.T(), s.db.UpdatePluginConf(&plugin))
	if pluginConf, err := s.db.GetPluginConfByToken("Pabc"); assert.NoError(s.T(), err) {
		assert.Equal(s.T(), false, pluginConf.Enabled)
		assert.Equal(s.T(), testConf, string(pluginConf.Config))
	}
}
