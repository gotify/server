package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/api"
	"github.com/gotify/server/auth"
	"github.com/gotify/server/database"
	"github.com/gotify/server/error"
	"github.com/gotify/server/ui"
	"github.com/jmattheis/go-packr-swagger-ui"

	"net/http"

	"github.com/gotify/location"
	"github.com/gotify/server/api/stream"
	"github.com/gotify/server/config"
	"github.com/gotify/server/docs"
	"github.com/gotify/server/mode"
	"github.com/gotify/server/model"
)

// Create creates the gin engine with all routes.
func Create(db *database.GormDatabase, vInfo *model.VersionInfo, conf *config.Configuration) (*gin.Engine, func()) {
	streamHandler := stream.New(200*time.Second, 15*time.Second)
	authentication := auth.Auth{DB: db}
	messageHandler := api.MessageAPI{Notifier: streamHandler, DB: db}
	clientHandler := api.ClientAPI{
		DB:            db,
		ImageDir:      conf.UploadedImagesDir,
		NotifyDeleted: streamHandler.NotifyDeletedClient,
	}
	applicationHandler := api.ApplicationAPI{
		DB:       db,
		ImageDir: conf.UploadedImagesDir,
	}
	userHandler := api.UserAPI{DB: db, PasswordStrength: conf.PassStrength, NotifyDeleted: streamHandler.NotifyDeletedUser}

	g := gin.New()

	g.Use(gin.Logger(), gin.Recovery(), error.Handler(), location.Default())
	g.NoRoute(error.NotFound())

	ui.Register(g)

	g.GET("/swagger", docs.Serve)
	g.Static("/image", conf.UploadedImagesDir)
	g.GET("/docs/*any", gin.WrapH(http.StripPrefix("/docs/", http.FileServer(swaggerui.GetBox()))))

	g.Use(func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/json")
		if mode.IsDev() {
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.Header("Access-Control-Allow-Methods", "GET,POST,DELETE,OPTIONS,PUT")
			ctx.Header("Access-Control-Allow-Headers", "X-Gotify-Key,Authorization,Content-Type,Upgrade,Origin,Connection,Accept-Encoding,Accept-Language,Host")
		}
	})

	g.OPTIONS("/*any")

	// swagger:operation GET /version version getVersion
	//
	// Get version information.
	//
	// ---
	// produces:
	// - application/json
	// responses:
	//   200:
	//     description: Ok
	//     schema:
	//         $ref: "#/definitions/VersionInfo"
	g.GET("version", func(ctx *gin.Context) {
		ctx.JSON(200, vInfo)
	})

	g.Group("/").Use(authentication.RequireApplicationToken()).POST("/message", messageHandler.CreateMessage)

	clientAuth := g.Group("")
	{
		clientAuth.Use(authentication.RequireClient())
		app := clientAuth.Group("/application")
		{

			app.GET("", applicationHandler.GetApplications)

			app.POST("", applicationHandler.CreateApplication)

			app.POST("/:id/image", applicationHandler.UploadApplicationImage)


			app.PUT("/:id", applicationHandler.UpdateApplication)

			app.DELETE("/:id", applicationHandler.DeleteApplication)

			tokenMessage := app.Group("/:id/message")
			{

				tokenMessage.GET("", messageHandler.GetMessagesWithApplication)

				tokenMessage.DELETE("", messageHandler.DeleteMessageWithApplication)
			}
		}

		client := clientAuth.Group("/client")
		{

			client.GET("", clientHandler.GetClients)

			client.POST("", clientHandler.CreateClient)

			client.DELETE("/:id", clientHandler.DeleteClient)
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
