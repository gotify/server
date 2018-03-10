package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/auth"
	"github.com/gotify/server/model"
)

// The TokenDatabase interface for encapsulating database access.
type TokenDatabase interface {
	CreateApplication(application *model.Application) error
	GetApplicationByToken(token string) *model.Application
	GetApplicationByID(id uint) *model.Application
	GetApplicationsByUser(userID uint) []*model.Application
	DeleteApplicationByID(id uint) error

	CreateClient(client *model.Client) error
	GetClientByToken(token string) *model.Client
	GetClientByID(id uint) *model.Client
	GetClientsByUser(userID uint) []*model.Client
	DeleteClientByID(id uint) error
}

// The TokenAPI provides handlers for managing clients and applications.
type TokenAPI struct {
	DB TokenDatabase
}

// CreateApplication creates an application and returns the access token.
func (a *TokenAPI) CreateApplication(ctx *gin.Context) {
	app := model.Application{}
	if err := ctx.Bind(&app); err == nil {
		app.Token = generateNotExistingToken(auth.GenerateApplicationToken, a.applicationExists)
		app.UserID = auth.GetUserID(ctx)
		a.DB.CreateApplication(&app)
		ctx.JSON(200, app)
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
	ctx.JSON(200, apps)
}

// GetClients returns all clients a user has.
func (a *TokenAPI) GetClients(ctx *gin.Context) {
	userID := auth.GetUserID(ctx)
	apps := a.DB.GetClientsByUser(userID)
	ctx.JSON(200, apps)
}

// DeleteApplication deletes an application by its id.
func (a *TokenAPI) DeleteApplication(ctx *gin.Context) {
	withID(ctx, "id", func(id uint) {
		if app := a.DB.GetApplicationByID(id); app != nil && app.UserID == auth.GetUserID(ctx) {
			a.DB.DeleteApplicationByID(id)
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
