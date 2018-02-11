package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmattheis/memo/api"
	"github.com/jmattheis/memo/auth"
	"github.com/jmattheis/memo/database"
	"github.com/jmattheis/memo/error"
	"github.com/jmattheis/memo/stream"
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

	g.GET("/")

	g.Group("/").Use(authentication.RequireApplicationToken()).POST("/message", messageHandler.CreateMessage)

	clientAuth := g.Group("")
	{
		clientAuth.Use(authentication.RequireClient())
		app := clientAuth.Group("/application")
		{
			app.GET("", tokenHandler.GetApplications)
			app.POST("", tokenHandler.CreateApplication)
			app.DELETE("/:id", tokenHandler.DeleteApplication)
			tokenMessage := app.Group("/:id/message")
			{
				tokenMessage.GET("", messageHandler.GetMessagesWithApplication)
				tokenMessage.DELETE("", messageHandler.DeleteMessageWithApplication)
			}
		}

		client := clientAuth.Group("/client")
		{
			client.GET("", tokenHandler.GetClients)
			client.POST("", tokenHandler.CreateClient)
			client.DELETE("/:id", tokenHandler.DeleteClient)
		}

		message := clientAuth.Group("/message")
		{
			message.GET("", messageHandler.GetMessages)
			message.DELETE("", messageHandler.DeleteMessages)
			message.DELETE("/:id", messageHandler.DeleteMessage)
		}

		clientAuth.GET("/stream", streamHandler.Handle)

		clientAuth.GET("current/user", userHandler.GetCurrentUser)
		clientAuth.POST("current/user/password", userHandler.ChangePassword)
	}

	authAdmin := g.Group("/user")
	{
		authAdmin.Use(authentication.RequireAdmin())

		authAdmin.GET("", userHandler.GetUsers)
		authAdmin.POST("", userHandler.CreateUser)

		authAdmin.DELETE("/:id", userHandler.DeleteUserByID)
		authAdmin.GET("/:id", userHandler.GetUserByID)
		authAdmin.POST("/:id", userHandler.UpdateUserByID)
	}
	return g, streamHandler.Close
}
