package database

import "github.com/jmattheis/memo/model"

// GetClientByID returns the client for the given id or nil.
func (d *GormDatabase) GetClientByID(id string) *model.Client {
	client := new(model.Client)
	d.DB.Where("id = ?", id).Find(client)
	if client.ID == id {
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
func (d *GormDatabase) DeleteClientByID(id string) error {
	return d.DB.Where("id = ?", id).Delete(&model.Client{}).Error
}
