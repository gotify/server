package database

import (
	"github.com/gotify/server/model"
)

// GetMessageByID returns the messages for the given id or nil.
func (d *GormDatabase) GetMessageByID(id uint) *model.Message {
	msg := new(model.Message)
	d.DB.Find(msg, id)
	if msg.ID == id {
		return msg
	}
	return nil
}

// CreateMessage creates a message.
func (d *GormDatabase) CreateMessage(message *model.Message) error {
	return d.DB.Create(message).Error
}

// GetMessagesByUser returns all messages from a user.
func (d *GormDatabase) GetMessagesByUser(userID uint) []*model.Message {
	var messages []*model.Message
	d.DB.Joins("JOIN applications ON applications.user_id = ?", userID).
		Where("messages.application_id = applications.id").Order("id desc").Find(&messages)
	return messages
}

// GetMessagesByUserSince returns limited messages from a user.
// If since is 0 it will be ignored.
func (d *GormDatabase) GetMessagesByUserSince(userID uint, limit int, since uint) []*model.Message {
	var messages []*model.Message
	db := d.DB.Joins("JOIN applications ON applications.user_id = ?", userID).
		Where("messages.application_id = applications.id").Order("id desc").Limit(limit)
	if since != 0 {
		db = db.Where("messages.id < ?", since)
	}
	db.Find(&messages)
	return messages
}

// GetMessagesByApplication returns all messages from an application.
func (d *GormDatabase) GetMessagesByApplication(tokenID uint) []*model.Message {
	var messages []*model.Message
	d.DB.Where("application_id = ?", tokenID).Order("id desc").Find(&messages)
	return messages
}

// GetMessagesByApplicationSince returns limited messages from an application.
// If since is 0 it will be ignored.
func (d *GormDatabase) GetMessagesByApplicationSince(appID uint, limit int, since uint) []*model.Message {
	var messages []*model.Message
	db := d.DB.Where("application_id = ?", appID).Order("id desc").Limit(limit)
	if since != 0 {
		db = db.Where("messages.id < ?", since)
	}
	db.Find(&messages)
	return messages
}

// DeleteMessageByID deletes a message by its id.
func (d *GormDatabase) DeleteMessageByID(id uint) error {
	return d.DB.Where("id = ?", id).Delete(&model.Message{}).Error
}

// DeleteMessagesByApplication deletes all messages from an application.
func (d *GormDatabase) DeleteMessagesByApplication(applicationID uint) error {
	return d.DB.Where("application_id = ?", applicationID).Delete(&model.Message{}).Error
}

// DeleteMessagesByUser deletes all messages from a user.
func (d *GormDatabase) DeleteMessagesByUser(userID uint) error {
	for _, app := range d.GetApplicationsByUser(userID) {
		d.DeleteMessagesByApplication(app.ID)
	}
	return nil
}
