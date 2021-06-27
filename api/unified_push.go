package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gotify/server/v2/auth"
	"github.com/gotify/server/v2/model"
)

// The UnifiedPushAPI provides handlers request for UnifiedPush
type UnifiedPushAPI struct {
	DB       MessageDatabase
	Notifier Notifier
}

// CreateMessage creates a message, authentication via application-token is required.
// swagger:operation POST /UP message createMessage
//
// Create a message.
//
// __NOTE__: This API ONLY accepts an application token as authentication.
// ---
// consumes: [application/json]
// produces: [application/json]
// security: [appTokenHeader: [], appTokenQuery: []]
// parameters:
// - name: body
//   in: body
//   description: the message of push notification
//   required: true
//   schema:
//     $ref: "#/definitions/string"
// responses:
//   200:
//     description: Ok
//     schema:
//       $ref: "#/definitions/Message"
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
func (a *UnifiedPushAPI) CreateMessage(ctx *gin.Context) {
	message := model.MessageExternal{}
	if err := ctx.Bind(&message.Message); err != nil {
		return
	}

	application, err := a.DB.GetApplicationByToken(auth.GetTokenID(ctx))
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}
	message.ApplicationID = application.ID
	message.Title = application.Name
	message.Date = timeNow()
	message.ID = 0
	msgInternal := toInternalMessage(&message)
	if success := successOrAbort(ctx, 500, a.DB.CreateMessage(msgInternal)); !success {
		return
	}
	a.Notifier.Notify(auth.GetUserID(ctx), toExternalMessage(msgInternal))
	ctx.JSON(200, toExternalMessage(msgInternal))
}
