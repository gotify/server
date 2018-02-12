package database

import (
	"github.com/gotify/server/model"
)

// GetApplicationByID returns the application for the given id or nil.
func (d *GormDatabase) GetApplicationByID(id string) *model.Application {
	app := new(model.Application)
	d.DB.Where("id = ?", id).Find(app)
	if app.ID == id {
		return app
	}
	return nil
}

// CreateApplication creates an application.
func (d *GormDatabase) CreateApplication(application *model.Application) error {
	return d.DB.Create(application).Error
}

// DeleteApplicationByID deletes an application by its id.
func (d *GormDatabase) DeleteApplicationByID(id string) error {
	return d.DB.Where("id = ?", id).Delete(&model.Application{}).Error
}

// GetApplicationsByUser returns all applications from a user.
func (d *GormDatabase) GetApplicationsByUser(userID uint) []*model.Application {
	var apps []*model.Application
	d.DB.Where("user_id = ?", userID).Find(&apps)
	return apps
}
