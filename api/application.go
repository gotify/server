package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gotify/location"
	"github.com/gotify/server/auth"
	"github.com/gotify/server/model"
	"github.com/h2non/filetype"
)

// The ApplicationDatabase interface for encapsulating database access.
type ApplicationDatabase interface {
	CreateApplication(application *model.Application) error
	GetApplicationByToken(token string) *model.Application
	GetApplicationByID(id uint) *model.Application
	GetApplicationsByUser(userID uint) []*model.Application
	DeleteApplicationByID(id uint) error
	UpdateApplication(application *model.Application) error
}

// The ApplicationAPI provides handlers for managing applications.
type ApplicationAPI struct {
	DB       ApplicationDatabase
	ImageDir string
}

// CreateApplication creates an application and returns the access token.
func (a *ApplicationAPI) CreateApplication(ctx *gin.Context) {
	app := model.Application{}
	if err := ctx.Bind(&app); err == nil {
		app.Token = generateNotExistingToken(auth.GenerateApplicationToken, a.applicationExists)
		app.UserID = auth.GetUserID(ctx)
		a.DB.CreateApplication(&app)
		ctx.JSON(200, withAbsoluteURL(ctx, &app))
	}
}

// GetApplications returns all applications a user has.
func (a *ApplicationAPI) GetApplications(ctx *gin.Context) {
	userID := auth.GetUserID(ctx)
	apps := a.DB.GetApplicationsByUser(userID)
	for _, app := range apps {
		withAbsoluteURL(ctx, app)
	}
	ctx.JSON(200, apps)
}

// DeleteApplication deletes an application by its id.
func (a *ApplicationAPI) DeleteApplication(ctx *gin.Context) {
	withID(ctx, "id", func(id uint) {
		if app := a.DB.GetApplicationByID(id); app != nil && app.UserID == auth.GetUserID(ctx) {
			a.DB.DeleteApplicationByID(id)
			if app.Image != "" {
				os.Remove(a.ImageDir + app.Image)
			}
		} else {
			ctx.AbortWithError(404, fmt.Errorf("app with id %d doesn't exists", id))
		}
	})
}

// UpdateApplication updates an application info by its id.
func (a *ApplicationAPI) UpdateApplication(ctx *gin.Context) {
	withID(ctx, "id", func(id uint) {
		if app := a.DB.GetApplicationByID(id); app != nil && app.UserID == auth.GetUserID(ctx) {
			newValues := &model.Application{}
			if err := ctx.Bind(newValues); err == nil {
				app.Description = newValues.Description
				app.Name = newValues.Name

				a.DB.UpdateApplication(app)

				ctx.JSON(200, withAbsoluteURL(ctx, app))
			}
		} else {
			ctx.AbortWithError(404, fmt.Errorf("app with id %d doesn't exists", id))
		}
	})
}

// UploadApplicationImage uploads an image for an application.
func (a *ApplicationAPI) UploadApplicationImage(ctx *gin.Context) {
	withID(ctx, "id", func(id uint) {
		if app := a.DB.GetApplicationByID(id); app != nil && app.UserID == auth.GetUserID(ctx) {
			file, err := ctx.FormFile("file")
			if err == http.ErrMissingFile {
				ctx.AbortWithError(400, errors.New("file with key 'file' must be present"))
				return
			} else if err != nil {
				ctx.AbortWithError(500, err)
				return
			}
			head := make([]byte, 261)
			open, _ := file.Open()
			open.Read(head)
			if !filetype.IsImage(head) {
				ctx.AbortWithError(400, errors.New("file must be an image"))
				return
			}

			ext := filepath.Ext(file.Filename)

			name := auth.GenerateImageName()
			for exist(a.ImageDir + name + ext) {
				name = auth.GenerateImageName()
			}

			err = ctx.SaveUploadedFile(file, a.ImageDir+name+ext)
			if err != nil {
				ctx.AbortWithError(500, err)
				return
			}

			if app.Image != "" {
				os.Remove(a.ImageDir + app.Image)
			}

			app.Image = name + ext
			a.DB.UpdateApplication(app)
			ctx.JSON(200, withAbsoluteURL(ctx, app))
		} else {
			ctx.AbortWithError(404, fmt.Errorf("client with id %d doesn't exists", id))
		}
	})
}

func (a *ApplicationAPI) applicationExists(token string) bool {
	return a.DB.GetApplicationByToken(token) != nil
}

func exist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func withAbsoluteURL(ctx *gin.Context, app *model.Application) *model.Application {
	url := location.Get(ctx)

	if app.Image == "" {
		url.Path = "static/defaultapp.png"
	} else {
		url.Path = "image/" + app.Image
	}
	app.Image = url.String()
	return app
}
