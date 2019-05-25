package database

import (
	"github.com/gotify/server/model"
	"github.com/jinzhu/gorm"
)

// GetClientByID returns the client for the given id or nil.
func (d *GormDatabase) GetClientByID(id uint) (*model.Client, error) {
	client := new(model.Client)
	err := d.DB.Where("id = ?", id).Find(client).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	if client.ID == id {
		return client, err
	}
	return nil, err
}

// GetClientByToken returns the client for the given token or nil.
func (d *GormDatabase) GetClientByToken(token string) (*model.Client, error) {
	client := new(model.Client)
	err := d.DB.Where("token = ?", token).Find(client).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	if client.Token == token {
		return client, err
	}
	return nil, err
}

// CreateClient creates a client.
func (d *GormDatabase) CreateClient(client *model.Client) error {
	return d.DB.Create(client).Error
}

// GetClientsByUser returns all clients from a user.
func (d *GormDatabase) GetClientsByUser(userID uint) ([]*model.Client, error) {
	var clients []*model.Client
	err := d.DB.Where("user_id = ?", userID).Find(&clients).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return clients, err
}

// DeleteClientByID deletes a client by its id.
func (d *GormDatabase) DeleteClientByID(id uint) error {
	return d.DB.Where("id = ?", id).Delete(&model.Client{}).Error
}

// UpdateClient updates a client.
func (d *GormDatabase) UpdateClient(client *model.Client) error {
	return d.DB.Save(client).Error
}
