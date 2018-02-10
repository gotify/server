package api

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmattheis/memo/auth"
	"github.com/jmattheis/memo/model"
)

// The MessageDatabase interface for encapsulating database access.
type MessageDatabase interface {
	GetMessagesByApplication(tokenID string) []*model.Message
	GetApplicationByID(id string) *model.Application
	GetMessagesByUser(userID uint) []*model.Message
	DeleteMessageByID(id uint) error
	GetMessageByID(id uint) *model.Message
	DeleteMessagesByUser(userID uint) error
	DeleteMessagesByApplication(applicationID string) error
	CreateMessage(message *model.Message) error
}

// Notifier notifies when a new message was created.
type Notifier interface {
	Notify(userID uint, message *model.Message)
}

// The MessageAPI provides handlers for managing messages.
type MessageAPI struct {
	DB MessageDatabase
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
	appID := ctx.Param("appid")
	messages := a.DB.GetMessagesByApplication(appID)
	ctx.JSON(200, messages)
}

// DeleteMessages delete all messages from a user.
func (a *MessageAPI) DeleteMessages(ctx *gin.Context) {
	userID := auth.GetUserID(ctx)
	a.DB.DeleteMessagesByUser(userID)
}

// DeleteMessageWithApplication deletes all messages from a specific application.
func (a *MessageAPI) DeleteMessageWithApplication(ctx *gin.Context) {
	appID := ctx.Param("appid")
	if application := a.DB.GetApplicationByID(appID); application != nil && application.UserID == auth.GetUserID(ctx) {
		a.DB.DeleteMessagesByApplication(appID)
	} else {
		ctx.AbortWithError(404, errors.New("application does not exists"))
	}
}

// DeleteMessage deletes a message with an id.
func (a *MessageAPI) DeleteMessage(ctx *gin.Context) {
	id := ctx.Param("id")
	if parsedUInt, err := strconv.ParseUint(id, 10, 32); err == nil {
		if msg := a.DB.GetMessageByID(uint(parsedUInt)); msg != nil && a.DB.GetApplicationByID(msg.ApplicationID).UserID == auth.GetUserID(ctx) {
			a.DB.DeleteMessageByID(uint(parsedUInt))
		} else {
			ctx.AbortWithError(404, errors.New("message does not exists"))
		}
	} else {
		ctx.AbortWithError(400, errors.New("message does not exist"))
	}
}

// CreateMessage creates a message, authentication via application-token is required.
func (a *MessageAPI) CreateMessage(ctx *gin.Context) {
	message := model.Message{}
	if err := ctx.Bind(&message); err == nil {
		message.ApplicationID = auth.GetTokenID(ctx)
		message.Date = time.Now()
		a.DB.CreateMessage(&message)
		a.Notifier.Notify(auth.GetUserID(ctx), &message)
		ctx.JSON(200, message)
	}
}
