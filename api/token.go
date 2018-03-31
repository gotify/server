package api

import (
	"fmt"

	"errors"
	"net/http"
	"path/filepath"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/gotify/location"
	"github.com/gotify/server/auth"
	"github.com/gotify/server/model"
	"github.com/h2non/filetype"
)

// The TokenDatabase interface for encapsulating database access.
type TokenDatabase interface {
	CreateApplication(application *model.Application) error
	GetApplicationByToken(token string) *model.Application
	GetApplicationByID(id uint) *model.Application
	GetApplicationsByUser(userID uint) []*model.Application
	DeleteApplicationByID(id uint) error
	UpdateApplication(application *model.Application)

	CreateClient(client *model.Client) error
	GetClientByToken(token string) *model.Client
	GetClientByID(id uint) *model.Client
	GetClientsByUser(userID uint) []*model.Client
	DeleteClientByID(id uint) error
}

// The TokenAPI provides handlers for managing clients and applications.
type TokenAPI struct {
	DB       TokenDatabase
	ImageDir string
}

// CreateApplication creates an application and returns the access token.
func (a *TokenAPI) CreateApplication(ctx *gin.Context) {
	app := model.Application{}
	if err := ctx.Bind(&app); err == nil {
		app.Token = generateNotExistingToken(auth.GenerateApplicationToken, a.applicationExists)
		app.UserID = auth.GetUserID(ctx)
		a.DB.CreateApplication(&app)
		ctx.JSON(200, withAbsoluteURL(ctx, &app))
	}
}

// CreateClient creates a client and returns the access token.
func (a *TokenAPI) CreateClient(ctx *gin.Context) {
	client := model.Client{}
	if err := ctx.Bind(&client); err == nil {
		client.Token = generateNotExistingToken(auth.GenerateClientToken, a.clientExists)
		client.UserID = auth.GetUserID(ctx)
		a.DB.CreateClient(&client)
		ctx.JSON(200, client)
	}
}

// GetApplications returns all applications a user has.
func (a *TokenAPI) GetApplications(ctx *gin.Context) {
	userID := auth.GetUserID(ctx)
	apps := a.DB.GetApplicationsByUser(userID)
	for _, app := range apps {
		withAbsoluteURL(ctx, app)
	}
	ctx.JSON(200, apps)
}

// GetClients returns all clients a user has.
func (a *TokenAPI) GetClients(ctx *gin.Context) {
	userID := auth.GetUserID(ctx)
	clients := a.DB.GetClientsByUser(userID)
	ctx.JSON(200, clients)
}

// DeleteApplication deletes an application by its id.
func (a *TokenAPI) DeleteApplication(ctx *gin.Context) {
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

// DeleteClient deletes a client by its id.
func (a *TokenAPI) DeleteClient(ctx *gin.Context) {
	withID(ctx, "id", func(id uint) {
		if client := a.DB.GetClientByID(id); client != nil && client.UserID == auth.GetUserID(ctx) {
			a.DB.DeleteClientByID(id)
		} else {
			ctx.AbortWithError(404, fmt.Errorf("client with id %d doesn't exists", id))
		}
	})
}

// UploadApplicationImage uploads an image for an application.
func (a *TokenAPI) UploadApplicationImage(ctx *gin.Context) {
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

func (a *TokenAPI) applicationExists(token string) bool {
	return a.DB.GetApplicationByToken(token) != nil
}

func (a *TokenAPI) clientExists(token string) bool {
	return a.DB.GetClientByToken(token) != nil
}

func generateNotExistingToken(generateToken func() string, tokenExists func(token string) bool) string {
	for {
		token := generateToken()
		if !tokenExists(token) {
			return token
		}
	}
}
