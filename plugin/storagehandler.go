package plugin

type dbStorageHandler struct {
	pluginID uint
	db       Database
}

func (c dbStorageHandler) Save(b []byte) error {
	conf := c.db.GetPluginConfByID(c.pluginID)
	conf.Storage = b
	return c.db.UpdatePluginConf(conf)
}

func (c dbStorageHandler) Load() ([]byte, error) {
	return c.db.GetPluginConfByID(c.pluginID).Storage, nil
}
