package database

import (
	"github.com/gotify/server/model"
)

// GetUserByName returns the user by the given name or nil.
func (d *GormDatabase) GetUserByName(name string) *model.User {
	user := new(model.User)
	d.DB.Where("name = ?", name).Find(user)
	if user.Name == name {
		return user
	}
	return nil
}

// GetUserByID returns the user by the given id or nil.
func (d *GormDatabase) GetUserByID(id uint) *model.User {
	user := new(model.User)
	d.DB.Find(user, id)
	if user.ID == id {
		return user
	}
	return nil
}

// CountUser returns the user count which satisfies the given condition.
func (d *GormDatabase) CountUser(condition ...interface{}) int {
	c := -1
	handle := d.DB.Model(new(model.User))
	if len(condition) == 1 {
		handle = handle.Where(condition[0])
	} else if len(condition) > 1 {
		handle = handle.Where(condition[0], condition[1:]...)
	}
	handle.Count(&c)
	return c
}

// GetUsers returns all users.
func (d *GormDatabase) GetUsers() []*model.User {
	var users []*model.User
	d.DB.Find(&users)
	return users
}

// DeleteUserByID deletes a user by its id.
func (d *GormDatabase) DeleteUserByID(id uint) error {
	for _, app := range d.GetApplicationsByUser(id) {
		d.DeleteApplicationByID(app.ID)
	}
	for _, client := range d.GetClientsByUser(id) {
		d.DeleteClientByID(client.ID)
	}
	for _, conf := range d.GetPluginConfByUser(id) {
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
