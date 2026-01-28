package database

import (
	"database/sql"
	"time"

	"github.com/gotify/server/v2/fracdex"
	"github.com/gotify/server/v2/model"
	"gorm.io/gorm"
)

// GetApplicationByToken returns the application for the given token or nil.
func (d *GormDatabase) GetApplicationByToken(token string) (*model.Application, error) {
	app := new(model.Application)
	err := d.DB.Where("token = ?", token).Find(app).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	if app.Token == token {
		return app, err
	}
	return nil, err
}

// GetApplicationByID returns the application for the given id or nil.
func (d *GormDatabase) GetApplicationByID(id uint) (*model.Application, error) {
	app := new(model.Application)
	err := d.DB.Where("id = ?", id).Find(app).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	if app.ID == id {
		return app, err
	}
	return nil, err
}

// CreateApplication creates an application.
func (d *GormDatabase) CreateApplication(application *model.Application) error {
	return d.DB.Transaction(func(tx *gorm.DB) error {
		if application.SortKey == "" {
			sortKey := ""
			err := tx.Model(&model.Application{}).Select("sort_key").Where("user_id = ?", application.UserID).Order("sort_key DESC").Limit(1).Find(&sortKey).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}
			application.SortKey, err = fracdex.KeyBetween(sortKey, "")
			if err != nil {
				return err
			}
		}

		return tx.Create(application).Error
	}, &sql.TxOptions{Isolation: sql.LevelSerializable})
}

// DeleteApplicationByID deletes an application by its id.
func (d *GormDatabase) DeleteApplicationByID(id uint) error {
	d.DeleteMessagesByApplication(id)
	return d.DB.Where("id = ?", id).Delete(&model.Application{}).Error
}

// GetApplicationsByUser returns all applications from a user.
func (d *GormDatabase) GetApplicationsByUser(userID uint) ([]*model.Application, error) {
	var apps []*model.Application
	err := d.DB.Where("user_id = ?", userID).Order("sort_key, id ASC").Find(&apps).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return apps, err
}

// UpdateApplication updates an application.
func (d *GormDatabase) UpdateApplication(app *model.Application) error {
	return d.DB.Save(app).Error
}

// UpdateApplicationTokenLastUsed updates the last used time of the application token.
func (d *GormDatabase) UpdateApplicationTokenLastUsed(token string, t *time.Time) error {
	return d.DB.Model(&model.Application{}).Where("token = ?", token).Update("last_used", t).Error
}
