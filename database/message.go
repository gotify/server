package database

import (
	"time"

	"github.com/gotify/server/v3/model"
	"gorm.io/gorm"
)

// GetMessageByID returns the messages for the given id or nil.
func (d *GormDatabase) GetMessageByID(id uint) (*model.Message, error) {
	msg := new(model.Message)
	err := d.DB.Find(msg, id).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	if msg.ID == id {
		return msg, err
	}
	return nil, err
}

// CreateMessage creates a message.
func (d *GormDatabase) CreateMessage(message *model.Message) error {
	return d.DB.Create(message).Error
}

// GetMessagesByUser returns all messages from a user.
func (d *GormDatabase) GetMessagesByUser(userID uint) ([]*model.Message, error) {
	var messages []*model.Message
	err := d.DB.Joins("JOIN applications ON applications.user_id = ?", userID).
		Where("messages.application_id = applications.id").Order("messages.id desc").Find(&messages).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return messages, err
}

// GetMessagesByUserSince returns limited messages from a user.
// If since is 0 it will be ignored.
func (d *GormDatabase) GetMessagesByUserSince(userID uint, limit int, since uint) ([]*model.Message, error) {
	var messages []*model.Message
	db := d.DB.Joins("JOIN applications ON applications.user_id = ?", userID).
		Where("messages.application_id = applications.id").Order("messages.id desc").Limit(limit)
	if since != 0 {
		db = db.Where("messages.id < ?", since)
	}
	err := db.Find(&messages).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return messages, err
}

// GetMessagesByApplication returns all messages from an application.
func (d *GormDatabase) GetMessagesByApplication(tokenID uint) ([]*model.Message, error) {
	var messages []*model.Message
	err := d.DB.Where("application_id = ?", tokenID).Order("messages.id desc").Find(&messages).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return messages, err
}

// GetMessagesByApplicationSince returns limited messages from an application.
// If since is 0 it will be ignored.
func (d *GormDatabase) GetMessagesByApplicationSince(appID uint, limit int, since uint) ([]*model.Message, error) {
	var messages []*model.Message
	db := d.DB.Where("application_id = ?", appID).Order("messages.id desc").Limit(limit)
	if since != 0 {
		db = db.Where("messages.id < ?", since)
	}
	err := db.Find(&messages).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return messages, err
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
	app, _ := d.GetApplicationsByUser(userID)
	for _, app := range app {
		d.DeleteMessagesByApplication(app.ID)
	}
	return nil
}

// PruneMessages deletes messages older than their application's retention
// period. Applications with RetentionSeconds == 0 fall back to
// globalDefaultRetentionSeconds; if that is also 0, those messages are kept
// forever. Returns the number of deleted messages.
func (d *GormDatabase) PruneMessages(now time.Time, globalDefaultRetentionSeconds int) (int64, error) {
	var apps []*model.Application
	if err := d.DB.Select("id", "retention_seconds").Find(&apps).Error; err != nil {
		return 0, err
	}

	var deleted int64
	var defaultRetentionAppIDs []uint
	for _, app := range apps {
		if app.RetentionSeconds == 0 {
			defaultRetentionAppIDs = append(defaultRetentionAppIDs, app.ID)
			continue
		}
		cutoff := now.Add(-time.Duration(app.RetentionSeconds) * time.Second)
		result := d.DB.Where("application_id = ? AND date <= ?", app.ID, cutoff).Delete(&model.Message{})
		if result.Error != nil {
			return deleted, result.Error
		}
		deleted += result.RowsAffected
	}

	if globalDefaultRetentionSeconds > 0 && len(defaultRetentionAppIDs) > 0 {
		cutoff := now.Add(-time.Duration(globalDefaultRetentionSeconds) * time.Second)
		result := d.DB.Where("application_id IN ? AND date <= ?", defaultRetentionAppIDs, cutoff).Delete(&model.Message{})
		if result.Error != nil {
			return deleted, result.Error
		}
		deleted += result.RowsAffected
	}

	return deleted, nil
}
