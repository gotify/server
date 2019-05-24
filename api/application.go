package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/auth"
	"github.com/gotify/server/model"
	"github.com/h2non/filetype"
)

// The ApplicationDatabase interface for encapsulating database access.
type ApplicationDatabase interface {
	CreateApplication(application *model.Application) error
	GetApplicationByToken(token string) (*model.Application, error)
	GetApplicationByID(id uint) (*model.Application, error)
	GetApplicationsByUser(userID uint) ([]*model.Application, error)
	DeleteApplicationByID(id uint) error
	UpdateApplication(application *model.Application) error
}

// The ApplicationAPI provides handlers for managing applications.
type ApplicationAPI struct {
	DB       ApplicationDatabase
	ImageDir string
}

// CreateApplication creates an application and returns the access token.
// swagger:operation POST /application application createApp
//
// Create an application.
//
// ---
// consumes: [application/json]
// produces: [application/json]
// security: [clientTokenHeader: [], clientTokenQuery: [], basicAuth: []]
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
func (a *ApplicationAPI) CreateApplication(ctx *gin.Context) {
	app := model.Application{}
	if err := ctx.Bind(&app); err == nil {
		app.Token = auth.GenerateNotExistingToken(generateApplicationToken, a.applicationExists)
		app.UserID = auth.GetUserID(ctx)
		app.Internal = false
		if success := successOrAbort(ctx, 500, a.DB.CreateApplication(&app)); !success {
			return
		}
		ctx.JSON(200, withResolvedImage(&app))
	}
}

// GetApplications returns all applications a user has.
// swagger:operation GET /application application getApps
//
// Return all applications.
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
//         $ref: "#/definitions/Application"
//   401:
//     description: Unauthorized
//     schema:
//         $ref: "#/definitions/Error"
//   403:
//     description: Forbidden
//     schema:
//         $ref: "#/definitions/Error"
func (a *ApplicationAPI) GetApplications(ctx *gin.Context) {
	userID := auth.GetUserID(ctx)
	apps, err := a.DB.GetApplicationsByUser(userID)
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}
	for _, app := range apps {
		withResolvedImage(app)
	}
	ctx.JSON(200, apps)
}

// DeleteApplication deletes an application by its id.
// swagger:operation DELETE /application/{id} application deleteApp
//
// Delete an application.
//
// ---
// consumes: [application/json]
// produces: [application/json]
// parameters:
// - name: id
//   in: path
//   description: the application id
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
func (a *ApplicationAPI) DeleteApplication(ctx *gin.Context) {
	withID(ctx, "id", func(id uint) {
		app, err := a.DB.GetApplicationByID(id)
		if success := successOrAbort(ctx, 500, err); !success {
			return
		}
		if app != nil && app.UserID == auth.GetUserID(ctx) {
			if app.Internal {
				ctx.AbortWithError(400, errors.New("cannot delete internal application"))
				return
			}
			if success := successOrAbort(ctx, 500, a.DB.DeleteApplicationByID(id)); !success {
				return
			}
			if app.Image != "" {
				os.Remove(a.ImageDir + app.Image)
			}
		} else {
			ctx.AbortWithError(404, fmt.Errorf("app with id %d doesn't exists", id))
		}
	})
}

// UpdateApplication updates an application info by its id.
// swagger:operation PUT /application/{id} application updateApplication
//
// Update an application.
//
// ---
// consumes: [application/json]
// produces: [application/json]
// security: [clientTokenHeader: [], clientTokenQuery: [], basicAuth: []]
// parameters:
// - name: body
//   in: body
//   description: the application to update
//   required: true
//   schema:
//     $ref: "#/definitions/Application"
// - name: id
//   in: path
//   description: the application id
//   required: true
//   type: integer
// responses:
//   200:
//     description: Ok
//     schema:
//         $ref: "#/definitions/Application"
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
func (a *ApplicationAPI) UpdateApplication(ctx *gin.Context) {
	withID(ctx, "id", func(id uint) {
		app, err := a.DB.GetApplicationByID(id)
		if success := successOrAbort(ctx, 500, err); !success {
			return
		}
		if app != nil && app.UserID == auth.GetUserID(ctx) {
			newValues := &model.Application{}
			if err := ctx.Bind(newValues); err == nil {
				app.Description = newValues.Description
				app.Name = newValues.Name

				if success := successOrAbort(ctx, 500, a.DB.UpdateApplication(app)); !success {
					return
				}
				ctx.JSON(200, withResolvedImage(app))

			}
		} else {
			ctx.AbortWithError(404, fmt.Errorf("app with id %d doesn't exists", id))
		}
	})
}

// UploadApplicationImage uploads an image for an application.
// swagger:operation POST /application/{id}/image application uploadAppImage
//
// Upload an image for an application.
//
// ---
// consumes:
// - multipart/form-data
// produces: [application/json]
// security: [clientTokenHeader: [], clientTokenQuery: [], basicAuth: []]
// parameters:
// - name: file
//   in: formData
//   description: the application image
//   required: true
//   type: file
// - name: id
//   in: path
//   description: the application id
//   required: true
//   type: integer
// responses:
//   200:
//     description: Ok
//     schema:
//         $ref: "#/definitions/Application"
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
//   500:
//     description: Server Error
//     schema:
//         $ref: "#/definitions/Error"
func (a *ApplicationAPI) UploadApplicationImage(ctx *gin.Context) {
	withID(ctx, "id", func(id uint) {
		app, err := a.DB.GetApplicationByID(id)
		if success := successOrAbort(ctx, 500, err); !success {
			return
		}
		if app != nil && app.UserID == auth.GetUserID(ctx) {
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

			name := generateNonExistingImageName(a.ImageDir, func() string {
				return generateImageName() + ext
			})

			err = ctx.SaveUploadedFile(file, a.ImageDir+name)
			if err != nil {
				ctx.AbortWithError(500, err)
				return
			}

			if app.Image != "" {
				os.Remove(a.ImageDir + app.Image)
			}

			app.Image = name
			if success := successOrAbort(ctx, 500, a.DB.UpdateApplication(app)); !success {
				return
			}
			ctx.JSON(200, withResolvedImage(app))
		} else {
			ctx.AbortWithError(404, fmt.Errorf("app with id %d doesn't exists", id))
		}
	})
}

func withResolvedImage(app *model.Application) *model.Application {
	if app.Image == "" {
		app.Image = "static/defaultapp.png"
	} else {
		app.Image = "image/" + app.Image
	}
	return app
}

func (a *ApplicationAPI) applicationExists(token string) bool {
	app, _ := a.DB.GetApplicationByToken(token)
	return app != nil
}

func exist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func generateNonExistingImageName(imgDir string, gen func() string) string {
	for {
		name := gen()
		if !exist(imgDir + name) {
			return name
		}
	}
}
