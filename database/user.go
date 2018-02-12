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

// GetUsers returns all users.
func (d *GormDatabase) GetUsers() []*model.User {
	var users []*model.User
	d.DB.Find(&users)
	return users
}

// DeleteUserByID deletes a user by its id.
func (d *GormDatabase) DeleteUserByID(id uint) error {
	return d.DB.Where("id = ?", id).Delete(&model.User{}).Error
}

// UpdateUser updates a user.
func (d *GormDatabase) UpdateUser(user *model.User) {
	d.DB.Save(user)
}

// CreateUser creates a user.
func (d *GormDatabase) CreateUser(user *model.User) error {
	return d.DB.Create(user).Error
}
