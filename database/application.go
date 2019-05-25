package database

import (
	"github.com/gotify/server/model"
	"github.com/jinzhu/gorm"
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
	return d.DB.Create(application).Error
}

// DeleteApplicationByID deletes an application by its id.
func (d *GormDatabase) DeleteApplicationByID(id uint) error {
	d.DeleteMessagesByApplication(id)
	return d.DB.Where("id = ?", id).Delete(&model.Application{}).Error
}

// GetApplicationsByUser returns all applications from a user.
func (d *GormDatabase) GetApplicationsByUser(userID uint) ([]*model.Application, error) {
	var apps []*model.Application
	err := d.DB.Where("user_id = ?", userID).Find(&apps).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return apps, err
}

// UpdateApplication updates an application.
func (d *GormDatabase) UpdateApplication(app *model.Application) error {
	return d.DB.Save(app).Error
}
