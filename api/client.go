package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/auth"
	"github.com/gotify/server/model"
)

// The ClientDatabase interface for encapsulating database access.
type ClientDatabase interface {
	CreateClient(client *model.Client) error
	GetClientByToken(token string) *model.Client
	GetClientByID(id uint) *model.Client
	GetClientsByUser(userID uint) []*model.Client
	DeleteClientByID(id uint) error
}

// The ClientAPI provides handlers for managing clients and applications.
type ClientAPI struct {
	DB            ClientDatabase
	ImageDir      string
	NotifyDeleted func(uint, string)
}

// CreateClient creates a client and returns the access token.
// swagger:operation POST /client client createClient
//
// Create a client.
//
// ---
// consumes: [application/json]
// produces: [application/json]
// security: [clientTokenHeader: [], clientTokenQuery: [], basicAuth: []]
// parameters:
// - name: body
//   in: body
//   description: the client to add
//   required: true
//   schema:
//     $ref: "#/definitions/Client"
// responses:
//   200:
//     description: Ok
//     schema:
//         $ref: "#/definitions/Client"
//   400:
//     description: Bad Request
//     schema:
//         $ref: "#/definitions/Error"
//   401:
//     description: Unauthorized
//     schema:
//         $ref: "#/definitions/Error"
//   403:
//     description: Forbidden
//     schema:
//         $ref: "#/definitions/Error"
func (a *ClientAPI) CreateClient(ctx *gin.Context) {
	client := model.Client{}
	if err := ctx.Bind(&client); err == nil {
		client.Token = generateNotExistingToken(auth.GenerateClientToken, a.clientExists)
		client.UserID = auth.GetUserID(ctx)
		a.DB.CreateClient(&client)
		ctx.JSON(200, client)
	}
}

// GetClients returns all clients a user has.
// swagger:operation GET /client client getClients
//
// Return all clients.
//
// ---
// consumes: [application/json]
// produces: [application/json]
// security: [clientTokenHeader: [], clientTokenQuery: [], basicAuth: []]
// responses:
//   200:
//     description: Ok
//     schema:
//       type: array
//       items:
//         $ref: "#/definitions/Client"
//   401:
//     description: Unauthorized
//     schema:
//         $ref: "#/definitions/Error"
//   403:
//     description: Forbidden
//     schema:
//         $ref: "#/definitions/Error"
func (a *ClientAPI) GetClients(ctx *gin.Context) {
	userID := auth.GetUserID(ctx)
	clients := a.DB.GetClientsByUser(userID)
	ctx.JSON(200, clients)
}

// DeleteClient deletes a client by its id.
// swagger:operation DELETE /client/{id} client deleteClient
//
// Delete a client.
//
// ---
// consumes: [application/json]
// produces: [application/json]
// parameters:
// - name: id
//   in: path
//   description: the client id
//   required: true
//   type: integer
// security: [clientTokenHeader: [], clientTokenQuery: [], basicAuth: []]
// responses:
//   200:
//     description: Ok
//   400:
//     description: Bad Request
//     schema:
//         $ref: "#/definitions/Error"
//   401:
//     description: Unauthorized
//     schema:
//         $ref: "#/definitions/Error"
//   403:
//     description: Forbidden
//     schema:
//         $ref: "#/definitions/Error"
//   404:
//     description: Not Found
//     schema:
//         $ref: "#/definitions/Error"
func (a *ClientAPI) DeleteClient(ctx *gin.Context) {
	withID(ctx, "id", func(id uint) {
		if client := a.DB.GetClientByID(id); client != nil && client.UserID == auth.GetUserID(ctx) {
			a.NotifyDeleted(client.UserID, client.Token)
			a.DB.DeleteClientByID(id)
		} else {
			ctx.AbortWithError(404, fmt.Errorf("client with id %d doesn't exists", id))
		}
	})
}

func (a *ClientAPI) clientExists(token string) bool {
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
