package database

import "github.com/gotify/server/model"

// GetClientByID returns the client for the given id or nil.
func (d *GormDatabase) GetClientByID(id uint) *model.Client {
	client := new(model.Client)
	d.DB.Where("id = ?", id).Find(client)
	if client.ID == id {
		return client
	}
	return nil
}

// GetClientByToken returns the client for the given token or nil.
func (d *GormDatabase) GetClientByToken(token string) *model.Client {
	client := new(model.Client)
	d.DB.Where("token = ?", token).Find(client)
	if client.Token == token {
		return client
	}
	return nil
}

// CreateClient creates a client.
func (d *GormDatabase) CreateClient(client *model.Client) error {
	return d.DB.Create(client).Error
}

// GetClientsByUser returns all clients from a user.
func (d *GormDatabase) GetClientsByUser(userID uint) []*model.Client {
	var clients []*model.Client
	d.DB.Where("user_id = ?", userID).Find(&clients)
	return clients
}

// DeleteClientByID deletes a client by its id.
func (d *GormDatabase) DeleteClientByID(id uint) error {
	return d.DB.Where("id = ?", id).Delete(&model.Client{}).Error
}

// UpdateClient updates a client.
func (d *GormDatabase) UpdateClient(client *model.Client) error {
	return d.DB.Save(client).Error
}
