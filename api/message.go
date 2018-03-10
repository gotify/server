package api

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/auth"
	"github.com/gotify/server/model"
)

// The MessageDatabase interface for encapsulating database access.
type MessageDatabase interface {
	GetMessagesByApplication(id uint) []*model.Message
	GetApplicationByID(id uint) *model.Application
	GetMessagesByUser(userID uint) []*model.Message
	DeleteMessageByID(id uint) error
	GetMessageByID(id uint) *model.Message
	DeleteMessagesByUser(userID uint) error
	DeleteMessagesByApplication(applicationID uint) error
	CreateMessage(message *model.Message) error
	GetApplicationByToken(token string) *model.Application
}

// Notifier notifies when a new message was created.
type Notifier interface {
	Notify(userID uint, message *model.Message)
}

// The MessageAPI provides handlers for managing messages.
type MessageAPI struct {
	DB       MessageDatabase
	Notifier Notifier
}

// GetMessages returns all messages from a user.
func (a *MessageAPI) GetMessages(ctx *gin.Context) {
	userID := auth.GetUserID(ctx)
	messages := a.DB.GetMessagesByUser(userID)
	ctx.JSON(200, messages)
}

// GetMessagesWithApplication returns all messages from a specific application.
func (a *MessageAPI) GetMessagesWithApplication(ctx *gin.Context) {
	withID(ctx, "appid", func(id uint) {
		if app := a.DB.GetApplicationByID(id); app != nil && app.UserID == auth.GetUserID(ctx) {
			messages := a.DB.GetMessagesByApplication(id)
			ctx.JSON(200, messages)
		} else {
			ctx.AbortWithError(404, errors.New("application does not exist"))
		}
	})
}

// DeleteMessages delete all messages from a user.
func (a *MessageAPI) DeleteMessages(ctx *gin.Context) {
	userID := auth.GetUserID(ctx)
	a.DB.DeleteMessagesByUser(userID)
}

// DeleteMessageWithApplication deletes all messages from a specific application.
func (a *MessageAPI) DeleteMessageWithApplication(ctx *gin.Context) {
	withID(ctx, "appid", func(id uint) {
		if application := a.DB.GetApplicationByID(id); application != nil && application.UserID == auth.GetUserID(ctx) {
			a.DB.DeleteMessagesByApplication(id)
		} else {
			ctx.AbortWithError(404, errors.New("application does not exists"))
		}
	});
}

// DeleteMessage deletes a message with an id.
func (a *MessageAPI) DeleteMessage(ctx *gin.Context) {
	withID(ctx, "id", func(id uint) {
		if msg := a.DB.GetMessageByID(id); msg != nil && a.DB.GetApplicationByID(msg.ApplicationID).UserID == auth.GetUserID(ctx) {
			a.DB.DeleteMessageByID(id)
		} else {
			ctx.AbortWithError(404, errors.New("message does not exists"))
		}
	})
}

// CreateMessage creates a message, authentication via application-token is required.
func (a *MessageAPI) CreateMessage(ctx *gin.Context) {
	message := model.Message{}
	if err := ctx.Bind(&message); err == nil {
		message.ApplicationID = a.DB.GetApplicationByToken(auth.GetTokenID(ctx)).ID
		message.Date = time.Now()
		a.DB.CreateMessage(&message)
		a.Notifier.Notify(auth.GetUserID(ctx), &message)
		ctx.JSON(200, message)
	}
}
