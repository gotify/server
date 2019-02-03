package database

import (
	"github.com/gotify/server/model"
)

// GetPluginConfByUser gets plugin configurations from a user
func (d *GormDatabase) GetPluginConfByUser(userid uint) []*model.PluginConf {
	var plugins []*model.PluginConf
	d.DB.Where("user_id = ?", userid).Find(&plugins)
	return plugins
}

// GetPluginConfByUserAndPath gets plugin configuration by user and file name
func (d *GormDatabase) GetPluginConfByUserAndPath(userid uint, path string) *model.PluginConf {
	plugin := new(model.PluginConf)
	d.DB.Where("user_id = ? AND module_path = ?", userid, path).First(plugin)
	if plugin.ModulePath == path {
		return plugin
	}
	return nil
}

// GetPluginConfByApplicationID gets plugin configuration by its internal appid.
func (d *GormDatabase) GetPluginConfByApplicationID(appid uint) *model.PluginConf {
	plugin := new(model.PluginConf)
	d.DB.Where("application_id = ?", appid).First(plugin)
	if plugin.ApplicationID == appid {
		return plugin
	}
	return nil
}

// CreatePluginConf creates a new plugin configuration
func (d *GormDatabase) CreatePluginConf(p *model.PluginConf) error {
	return d.DB.Create(p).Error
}

// GetPluginConfByToken gets plugin configuration by plugin token
func (d *GormDatabase) GetPluginConfByToken(token string) *model.PluginConf {
	plugin := new(model.PluginConf)
	d.DB.Where("token = ?", token).First(plugin)
	if plugin.Token == token {
		return plugin
	}
	return nil
}

// GetPluginConfByID gets plugin configuration by plugin ID
func (d *GormDatabase) GetPluginConfByID(id uint) *model.PluginConf {
	plugin := new(model.PluginConf)
	d.DB.Where("id = ?", id).First(plugin)
	if plugin.ID == id {
		return plugin
	}
	return nil
}

// UpdatePluginConf updates plugin configuration
func (d *GormDatabase) UpdatePluginConf(p *model.PluginConf) error {
	return d.DB.Save(p).Error
}

// DeletePluginConfByID deletes a plugin configuration by its id.
func (d *GormDatabase) DeletePluginConfByID(id uint) error {
	return d.DB.Where("id = ?", id).Delete(&model.PluginConf{}).Error
}
