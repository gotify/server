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
	GetApplicationByID(id string) *model.Application
	GetApplicationsByUser(userID uint) []*model.Application
	DeleteApplicationByID(id string) error

	CreateClient(client *model.Client) error
	GetClientByID(id string) *model.Client
	GetClientsByUser(userID uint) []*model.Client
	DeleteClientByID(id string) error
}

// The TokenAPI provides handlers for managing clients and applications.
type TokenAPI struct {
	DB TokenDatabase
}

// CreateApplication creates an application and returns the access token.
func (a *TokenAPI) CreateApplication(ctx *gin.Context) {
	app := model.Application{}
	if err := ctx.Bind(&app); err == nil {
		app.ID = generateNotExistingToken(auth.GenerateApplicationToken, a.applicationExists)
		app.UserID = auth.GetUserID(ctx)
		a.DB.CreateApplication(&app)
		ctx.JSON(200, app)
	}
}

// CreateClient creates a client and returns the access token.
func (a *TokenAPI) CreateClient(ctx *gin.Context) {
	client := model.Client{}
	if err := ctx.Bind(&client); err == nil {
		client.ID = generateNotExistingToken(auth.GenerateClientToken, a.clientExists)
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
	appID := ctx.Param("id")
	if app := a.DB.GetApplicationByID(appID); app != nil && app.UserID == auth.GetUserID(ctx) {
		a.DB.DeleteApplicationByID(appID)
	} else {
		ctx.AbortWithError(404, fmt.Errorf("app with id %s doesn't exists", appID))
	}
}

// DeleteClient deletes a client by its id.
func (a *TokenAPI) DeleteClient(ctx *gin.Context) {
	clientID := ctx.Param("id")
	if client := a.DB.GetClientByID(clientID); client != nil && client.UserID == auth.GetUserID(ctx) {
		a.DB.DeleteClientByID(clientID)
	} else {
		ctx.AbortWithError(404, fmt.Errorf("client with id %s doesn't exists", clientID))
	}
}

func (a *TokenAPI) applicationExists(appID string) bool {
	return a.DB.GetApplicationByID(appID) != nil
}

func (a *TokenAPI) clientExists(clientID string) bool {
	return a.DB.GetClientByID(clientID) != nil
}

func generateNotExistingToken(generateToken func() string, tokenExists func(token string) bool) string {
	for {
		token := generateToken()
		if !tokenExists(token) {
			return token
		}
	}
}
