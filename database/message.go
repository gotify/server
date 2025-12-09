package database

import (
	"time"

	"github.com/gotify/server/v2/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// GetMessagesByUserPaginated returns limited messages from a user.
// If since is 0 it will be ignored.
func (d *GormDatabase) GetMessagesByUserPaginated(userID uint, limit int, since uint64, after uint64, by string) ([]*model.Message, error) {
	var messages []*model.Message
	db := d.DB.Joins("JOIN applications ON applications.user_id = ?", userID).
		Where("messages.application_id = applications.id").Order(clause.OrderBy{Columns: []clause.OrderByColumn{
		{
			Column: clause.Column{
				Table: "messages",
				Name:  by,
			},
			Desc: since != 0 || after == 0,
		},
	}}).Limit(limit)
	if since != 0 {
		sinceVal := any(since)
		if by == "date" {
			sinceVal = time.Unix(int64(since), 0)
		}
		db = db.Where(clause.Lt{Column: clause.Column{Table: "messages", Name: by}, Value: sinceVal})
	}
	if after != 0 {
		afterVal := any(after)
		if by == "date" {
			afterVal = time.Unix(int64(after), 0)
		}
		db = db.Where(clause.Gte{Column: clause.Column{Table: "messages", Name: by}, Value: afterVal})
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

// GetMessagesByApplicationPaginated returns limited messages from an application.
// If since is 0 it will be ignored.
func (d *GormDatabase) GetMessagesByApplicationPaginated(appID uint, limit int, since uint64, after uint64, by string) ([]*model.Message, error) {
	var messages []*model.Message
	db := d.DB.Where("application_id = ?", appID).Order(clause.OrderBy{Columns: []clause.OrderByColumn{
		{
			Column: clause.Column{
				Table: "messages",
				Name:  by,
			},
			Desc: since != 0 || after == 0,
		},
	}}).Limit(limit)
	if since != 0 {
		sinceVal := any(since)
		if by == "date" {
			sinceVal = time.Unix(int64(since), 0)
		}
		db = db.Where(clause.Lt{Column: clause.Column{Table: "messages", Name: by}, Value: sinceVal})
	}
	if after != 0 {
		afterVal := any(after)
		if by == "date" {
			afterVal = time.Unix(int64(after), 0)
		}
		db = db.Where(clause.Gte{Column: clause.Column{Table: "messages", Name: by}, Value: afterVal})
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
