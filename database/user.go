package database

import (
	"github.com/gotify/server/v3/model"
	"gorm.io/gorm"
)

// GetUserByName returns the user by the given name or nil.
func (d *GormDatabase) GetUserByName(name string) (*model.User, error) {
	user := new(model.User)
	err := d.DB.Where("name = ?", name).Find(user).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	if user.Name == name {
		return user, err
	}
	return nil, err
}

// GetUserByOIDC returns the user bound to the given oidc id or nil.
func (d *GormDatabase) GetUserByOIDC(oidcID string) (*model.User, error) {
	user := new(model.User)
	err := d.DB.Where("oidc_id = ?", oidcID).Find(user).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	if user.OIDCID != nil && *user.OIDCID == oidcID {
		return user, err
	}
	return nil, err
}

// GetUserByID returns the user by the given id or nil.
func (d *GormDatabase) GetUserByID(id uint) (*model.User, error) {
	user := new(model.User)
	err := d.DB.Find(user, id).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	if user.ID == id {
		return user, err
	}
	return nil, err
}

// GetUsers returns the users which satisfy the given condition.
func (d *GormDatabase) GetUsers(condition ...any) ([]*model.User, error) {
	users := make([]*model.User, 0)
	handle := d.DB.Model(new(model.User))
	if len(condition) == 1 {
		handle = handle.Where(condition[0])
	} else if len(condition) > 1 {
		handle = handle.Where(condition[0], condition[1:]...)
	}
	err := handle.Find(&users).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return users, err
}

// DeleteUserByID deletes a user by its id.
func (d *GormDatabase) DeleteUserByID(id uint) error {
	apps, _ := d.GetApplicationsByUser(id)
	for _, app := range apps {
		d.DeleteApplicationByID(app.ID)
	}
	clients, _ := d.GetClientsByUser(id)
	for _, client := range clients {
		d.DeleteClientByID(client.ID)
	}
	pluginConfs, _ := d.GetPluginConfByUser(id)
	for _, conf := range pluginConfs {
		d.DeletePluginConfByID(conf.ID)
	}
	return d.DB.Where("id = ?", id).Delete(&model.User{}).Error
}

// UpdateUser updates a user.
func (d *GormDatabase) UpdateUser(user *model.User) error {
	return d.DB.Save(user).Error
}

// CreateUser creates a user.
func (d *GormDatabase) CreateUser(user *model.User) error {
	return d.DB.Create(user).Error
}
