package database

import (
	"github.com/gotify/server/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	pluginConf, err := s.db.GetPluginConfByUserAndPath(1, "github.com/gotify/example-plugin")
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "Pabc", pluginConf.Token)

	pluginConf, err = s.db.GetPluginConfByToken("Pabc")
	require.NoError(s.T(), err)
	assert.Equal(s.T(), true, pluginConf.Enabled)

	pluginConf, err = s.db.GetPluginConfByApplicationID(2)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "Pabc", pluginConf.Token)

	pluginConf, err = s.db.GetPluginConfByID(1)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "github.com/gotify/example-plugin", pluginConf.ModulePath)

	pluginConf, err = s.db.GetPluginConfByToken("Pnotexist")
	require.NoError(s.T(), err)
	assert.Nil(s.T(), pluginConf)

	pluginConf, err = s.db.GetPluginConfByID(12)
	require.NoError(s.T(), err)
	assert.Nil(s.T(), pluginConf)

	pluginConf, err = s.db.GetPluginConfByUserAndPath(1, "not/exist")
	require.NoError(s.T(), err)
	assert.Nil(s.T(), pluginConf)

	pluginConf, err = s.db.GetPluginConfByApplicationID(99)
	require.NoError(s.T(), err)
	assert.Nil(s.T(), pluginConf)

	pluginConfs, err := s.db.GetPluginConfByUser(1)
	require.NoError(s.T(), err)
	assert.Len(s.T(), pluginConfs, 1)

	pluginConfs, err = s.db.GetPluginConfByUser(0)
	require.NoError(s.T(), err)
	assert.Len(s.T(), pluginConfs, 0)

	testConf := `{"test_config_key":"hello"}`
	plugin.Enabled = false
	plugin.Config = []byte(testConf)
	assert.Nil(s.T(), s.db.UpdatePluginConf(&plugin))
	pluginConf, err = s.db.GetPluginConfByToken("Pabc")
	require.NoError(s.T(), err)
	assert.Equal(s.T(), false, pluginConf.Enabled)
	assert.Equal(s.T(), testConf, string(pluginConf.Config))

}
