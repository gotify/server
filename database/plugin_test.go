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
	assert.Equal(s.T(), "Pabc", s.db.GetPluginConfByUserAndPath(1, "github.com/gotify/example-plugin").Token)
	assert.Equal(s.T(), true, s.db.GetPluginConfByToken("Pabc").Enabled)
	assert.Equal(s.T(), "Pabc", s.db.GetPluginConfByApplicationID(2).Token)
	assert.Equal(s.T(), "github.com/gotify/example-plugin", s.db.GetPluginConfByID(1).ModulePath)

	assert.Nil(s.T(), s.db.GetPluginConfByToken("Pnotexist"))
	assert.Nil(s.T(), s.db.GetPluginConfByID(12))
	assert.Nil(s.T(), s.db.GetPluginConfByUserAndPath(1, "not/exist"))
	assert.Nil(s.T(), s.db.GetPluginConfByApplicationID(99))

	assert.Len(s.T(), s.db.GetPluginConfByUser(1), 1)
	assert.Len(s.T(), s.db.GetPluginConfByUser(0), 0)

	testConf := `{"test_config_key":"hello"}`
	plugin.Enabled = false
	plugin.Config = []byte(testConf)
	assert.Nil(s.T(), s.db.UpdatePluginConf(&plugin))
	assert.Equal(s.T(), false, s.db.GetPluginConfByToken("Pabc").Enabled)
	assert.Equal(s.T(), testConf, string(s.db.GetPluginConfByToken("Pabc").Config))
}
