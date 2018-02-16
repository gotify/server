package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/api"
	"github.com/gotify/server/auth"
	"github.com/gotify/server/database"
	"github.com/gotify/server/error"
	"github.com/jmattheis/go-packr-swagger-ui"

	"net/http"

	"github.com/gotify/server/docs"
	"github.com/gotify/server/stream"
)

// Create creates the gin engine with all routes.
func Create(db *database.GormDatabase) (*gin.Engine, func()) {
	streamHandler := stream.New(200*time.Second, 15*time.Second)
	authentication := auth.Auth{DB: db}
	messageHandler := api.MessageAPI{Notifier: streamHandler, DB: db}
	tokenHandler := api.TokenAPI{DB: db}
	userHandler := api.UserAPI{DB: db}

	g := gin.New()

	g.Use(gin.Logger(), gin.Recovery(), error.Handler())
	g.NoRoute(error.NotFound())
	g.GET("/")
	g.GET("/swagger", docs.Serve)
	g.GET("/docs/*any", gin.WrapH(http.StripPrefix("/docs/", http.FileServer(swaggerui.GetBox()))))

	g.Use(func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/json")
	})

	// swagger:operation POST /message message createMessage
	//
	// Create a message.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// security:
	// - appTokenHeader: []
	// - appTokenQuery: []
	// parameters:
	// - name: body
	//   in: body
	//   description: the message to add
	//   required: true
	//   schema:
	//     $ref: "#/definitions/Message"
	// responses:
	//   200:
	//     description: Ok
	//     schema:
	//       type: array
	//       items:
	//         $ref: "#/definitions/Message"
	//   401:
	//     description: Unauthorized
	//     schema:
	//         $ref: "#/definitions/Error"
	g.Group("/").Use(authentication.RequireApplicationToken()).POST("/message", messageHandler.CreateMessage)

	clientAuth := g.Group("")
	{
		clientAuth.Use(authentication.RequireClient())
		app := clientAuth.Group("/application")
		{
			// swagger:operation GET /application token getApps
			//
			// Return all applications.
			//
			// ---
			// consumes:
			// - application/json
			// produces:
			// - application/json
			// security:
			// - clientTokenHeader: []
			// - clientTokenQuery: []
			// - basicAuth: []
			// responses:
			//   200:
			//     description: Ok
			//     schema:
			//       type: array
			//       items:
			//         $ref: "#/definitions/Application"
			//   401:
			//     description: Unauthorized
			//     schema:
			//         $ref: "#/definitions/Error"
			//   403:
			//     description: Forbidden
			//     schema:
			//         $ref: "#/definitions/Error"
			app.GET("", tokenHandler.GetApplications)

			// swagger:operation POST /application token createApp
			//
			// Create an application.
			//
			// ---
			// consumes:
			// - application/json
			// produces:
			// - application/json
			// security:
			// - clientTokenHeader: []
			// - clientTokenQuery: []
			// - basicAuth: []
			// parameters:
			// - name: body
			//   in: body
			//   description: the application to add
			//   required: true
			//   schema:
			//     $ref: "#/definitions/Application"
			// responses:
			//   200:
			//     description: Ok
			//     schema:
			//         $ref: "#/definitions/Application"
			//   401:
			//     description: Unauthorized
			//     schema:
			//         $ref: "#/definitions/Error"
			//   403:
			//     description: Forbidden
			//     schema:
			//         $ref: "#/definitions/Error"
			app.POST("", tokenHandler.CreateApplication)

			// swagger:operation DELETE /application/{id} token deleteApp
			//
			// Delete an application.
			//
			// ---
			// consumes:
			// - application/json
			// produces:
			// - application/json
			// parameters:
			// - name: id
			//   in: path
			//   description: the application id
			//   required: true
			//   type: string
			// security:
			// - clientTokenHeader: []
			// - clientTokenQuery: []
			// - basicAuth: []
			// responses:
			//   200:
			//     description: Ok
			//   401:
			//     description: Unauthorized
			//     schema:
			//         $ref: "#/definitions/Error"
			//   403:
			//     description: Forbidden
			//     schema:
			//         $ref: "#/definitions/Error"
			app.DELETE("/:id", tokenHandler.DeleteApplication)

			tokenMessage := app.Group("/:id/message")
			{
				// swagger:operation GET /application/{id}/message message getAppMessages
				//
				// Return all messages from a specific application.
				//
				// ---
				// produces:
				// - application/json
				// security:
				// - clientTokenHeader: []
				// - clientTokenQuery: []
				// - basicAuth: []
				// parameters:
				// - name: id
				//   in: path
				//   description: the application id
				//   required: true
				//   type: string
				// responses:
				//   200:
				//     description: Ok
				//     schema:
				//       type: array
				//       items:
				//         $ref: "#/definitions/Message"
				//   401:
				//     description: Unauthorized
				//     schema:
				//         $ref: "#/definitions/Error"
				//   403:
				//     description: Forbidden
				//     schema:
				//         $ref: "#/definitions/Error"
				tokenMessage.GET("", messageHandler.GetMessagesWithApplication)

				// swagger:operation DELETE /application/{id}/message message deleteAppMessages
				//
				// Delete all messages from a specific application.
				//
				// ---
				// produces:
				// - application/json
				// security:
				// - clientTokenHeader: []
				// - clientTokenQuery: []
				// - basicAuth: []
				// parameters:
				// - name: id
				//   in: path
				//   description: the application id
				//   required: true
				//   type: string
				// responses:
				//   200:
				//     description: Ok
				//   401:
				//     description: Unauthorized
				//     schema:
				//         $ref: "#/definitions/Error"
				//   403:
				//     description: Forbidden
				//     schema:
				//         $ref: "#/definitions/Error"
				tokenMessage.DELETE("", messageHandler.DeleteMessageWithApplication)
			}
		}

		client := clientAuth.Group("/client")
		{
			// swagger:operation GET /client token getClients
			//
			// Return all clients.
			//
			// ---
			// consumes:
			// - application/json
			// produces:
			// - application/json
			// security:
			// - clientTokenHeader: []
			// - clientTokenQuery: []
			// - basicAuth: []
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
			client.GET("", tokenHandler.GetClients)

			// swagger:operation POST /client token createClient
			//
			// Create a client.
			//
			// ---
			// consumes:
			// - application/json
			// produces:
			// - application/json
			// security:
			// - clientTokenHeader: []
			// - clientTokenQuery: []
			// - basicAuth: []
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
			//   401:
			//     description: Unauthorized
			//     schema:
			//         $ref: "#/definitions/Error"
			//   403:
			//     description: Forbidden
			//     schema:
			//         $ref: "#/definitions/Error"
			client.POST("", tokenHandler.CreateClient)

			// swagger:operation DELETE /client/{id} token deleteClient
			//
			// Delete a client.
			//
			// ---
			// consumes:
			// - application/json
			// produces:
			// - application/json
			// parameters:
			// - name: id
			//   in: path
			//   description: the client id
			//   required: true
			//   type: string
			// security:
			// - clientTokenHeader: []
			// - clientTokenQuery: []
			// - basicAuth: []
			// responses:
			//   200:
			//     description: Ok
			//   401:
			//     description: Unauthorized
			//     schema:
			//         $ref: "#/definitions/Error"
			//   403:
			//     description: Forbidden
			//     schema:
			//         $ref: "#/definitions/Error"
			client.DELETE("/:id", tokenHandler.DeleteClient)
		}

		message := clientAuth.Group("/message")
		{
			// swagger:operation GET /message message getMessages
			//
			// Return all messages.
			//
			// ---
			// produces:
			// - application/json
			// security:
			// - clientTokenHeader: []
			// - clientTokenQuery: []
			// - basicAuth: []
			// responses:
			//   200:
			//     description: Ok
			//     schema:
			//       type: array
			//       items:
			//         $ref: "#/definitions/Message"
			//   401:
			//     description: Unauthorized
			//     schema:
			//         $ref: "#/definitions/Error"
			//   403:
			//     description: Forbidden
			//     schema:
			//         $ref: "#/definitions/Error"
			message.GET("", messageHandler.GetMessages)

			// swagger:operation DELETE /message message deleteMessages
			//
			// Delete all messages.
			//
			// ---
			// produces:
			// - application/json
			// security:
			// - clientTokenHeader: []
			// - clientTokenQuery: []
			// - basicAuth: []
			// responses:
			//   200:
			//     description: Ok
			//   401:
			//     description: Unauthorized
			//     schema:
			//         $ref: "#/definitions/Error"
			//   403:
			//     description: Forbidden
			//     schema:
			//         $ref: "#/definitions/Error"
			message.DELETE("", messageHandler.DeleteMessages)

			// swagger:operation DELETE /message/{id} message deleteMessage
			//
			// Deletes a message with an id.
			//
			// ---
			// produces:
			// - application/json
			// security:
			// - clientTokenHeader: []
			// - clientTokenQuery: []
			// - basicAuth: []
			// parameters:
			// - name: id
			//   in: path
			//   description: the message id
			//   required: true
			//   type: integer
			// responses:
			//   200:
			//     description: Ok
			//   401:
			//     description: Unauthorized
			//     schema:
			//         $ref: "#/definitions/Error"
			//   403:
			//     description: Forbidden
			//     schema:
			//         $ref: "#/definitions/Error"
			message.DELETE("/:id", messageHandler.DeleteMessage)
		}

		// swagger:operation GET /stream message streamMessages
		//
		// Websocket, return newly created messages.
		//
		// ---
		// schema: ws, wss
		// produces:
		// - application/json
		// security:
		// - clientTokenHeader: []
		// - clientTokenQuery: []
		// - basicAuth: []
		// responses:
		//   200:
		//     description: Ok
		//     schema:
		//         $ref: "#/definitions/Message"
		//   401:
		//     description: Unauthorized
		//     schema:
		//         $ref: "#/definitions/Error"
		//   403:
		//     description: Forbidden
		//     schema:
		//         $ref: "#/definitions/Error"
		clientAuth.GET("/stream", streamHandler.Handle)

		// swagger:operation GET /current/user user currentUser
		//
		// Return the current user.
		//
		// ---
		// produces:
		// - application/json
		// security:
		// - clientTokenHeader: []
		// - clientTokenQuery: []
		// - basicAuth: []
		// responses:
		//   200:
		//     description: Ok
		//     schema:
		//         $ref: "#/definitions/User"
		//   401:
		//     description: Unauthorized
		//     schema:
		//         $ref: "#/definitions/Error"
		//   403:
		//     description: Forbidden
		//     schema:
		//         $ref: "#/definitions/Error"
		clientAuth.GET("current/user", userHandler.GetCurrentUser)

		// swagger:operation POST /current/user/password user updateCurrentUser
		//
		// Update the password of the current user.
		//
		// ---
		// consumes:
		// - application/json
		// produces:
		// - application/json
		// security:
		// - clientTokenHeader: []
		// - clientTokenQuery: []
		// - basicAuth: []
		// parameters:
		// - name: body
		//   in: body
		//   description: the user
		//   required: true
		//   schema:
		//     $ref: "#/definitions/User"
		// responses:
		//   200:
		//     description: Ok
		//   401:
		//     description: Unauthorized
		//     schema:
		//         $ref: "#/definitions/Error"
		//   403:
		//     description: Forbidden
		//     schema:
		//         $ref: "#/definitions/Error"
		clientAuth.POST("current/user/password", userHandler.ChangePassword)
	}

	authAdmin := g.Group("/user")
	{
		authAdmin.Use(authentication.RequireAdmin())

		// swagger:operation GET /user user getUsers
		//
		// Return all users.
		//
		// ---
		// produces:
		// - application/json
		// security:
		// - clientTokenHeader: []
		// - clientTokenQuery: []
		// - basicAuth: []
		// responses:
		//   200:
		//     description: Ok
		//     schema:
		//       type: array
		//       items:
		//         $ref: "#/definitions/User"
		//   401:
		//     description: Unauthorized
		//     schema:
		//         $ref: "#/definitions/Error"
		//   403:
		//     description: Forbidden
		//     schema:
		//         $ref: "#/definitions/Error"
		authAdmin.GET("", userHandler.GetUsers)

		// swagger:operation POST /user user createUser
		//
		// Create a user.
		//
		// ---
		// consumes:
		// - application/json
		// produces:
		// - application/json
		// security:
		// - clientTokenHeader: []
		// - clientTokenQuery: []
		// - basicAuth: []
		// parameters:
		// - name: body
		//   in: body
		//   description: the user to add
		//   required: true
		//   schema:
		//     $ref: "#/definitions/User"
		// responses:
		//   200:
		//     description: Ok
		//     schema:
		//         $ref: "#/definitions/User"
		//   401:
		//     description: Unauthorized
		//     schema:
		//         $ref: "#/definitions/Error"
		//   403:
		//     description: Forbidden
		//     schema:
		//         $ref: "#/definitions/Error"
		authAdmin.POST("", userHandler.CreateUser)

		// swagger:operation DELETE /user/{id} user deleteUser
		//
		// Deletes a user.
		//
		// ---
		// produces:
		// - application/json
		// security:
		// - clientTokenHeader: []
		// - clientTokenQuery: []
		// - basicAuth: []
		// parameters:
		// - name: id
		//   in: path
		//   description: the user id
		//   required: true
		//   type: integer
		// responses:
		//   200:
		//     description: Ok
		//   401:
		//     description: Unauthorized
		//     schema:
		//         $ref: "#/definitions/Error"
		//   403:
		//     description: Forbidden
		//     schema:
		//         $ref: "#/definitions/Error"
		authAdmin.DELETE("/:id", userHandler.DeleteUserByID)

		// swagger:operation GET /user/{id} user getUser
		//
		// Get a user.
		//
		// ---
		// consumes:
		// - application/json
		// produces:
		// - application/json
		// security:
		// - clientTokenHeader: []
		// - clientTokenQuery: []
		// - basicAuth: []
		// parameters:
		// - name: id
		//   in: path
		//   description: the user id
		//   required: true
		//   type: integer
		// responses:
		//   200:
		//     description: Ok
		//     schema:
		//         $ref: "#/definitions/User"
		//   401:
		//     description: Unauthorized
		//     schema:
		//         $ref: "#/definitions/Error"
		//   403:
		//     description: Forbidden
		//     schema:
		//         $ref: "#/definitions/Error"
		authAdmin.GET("/:id", userHandler.GetUserByID)

		// swagger:operation POST /user/{id} user updateUser
		//
		// Update a user.
		//
		// ---
		// consumes:
		// - application/json
		// produces:
		// - application/json
		// security:
		// - clientTokenHeader: []
		// - clientTokenQuery: []
		// - basicAuth: []
		// parameters:
		// - name: id
		//   in: path
		//   description: the user id
		//   required: true
		//   type: integer
		// - name: body
		//   in: body
		//   description: the updated user
		//   required: true
		//   schema:
		//     $ref: "#/definitions/User"
		// responses:
		//   200:
		//     description: Ok
		//     schema:
		//         $ref: "#/definitions/User"
		//   401:
		//     description: Unauthorized
		//     schema:
		//         $ref: "#/definitions/Error"
		//   403:
		//     description: Forbidden
		//     schema:
		//         $ref: "#/definitions/Error"
		authAdmin.POST("/:id", userHandler.UpdateUserByID)
	}
	return g, streamHandler.Close
}
