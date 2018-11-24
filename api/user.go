package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/auth"
	"github.com/gotify/server/auth/password"
	"github.com/gotify/server/model"
)

// The UserDatabase interface for encapsulating database access.
type UserDatabase interface {
	GetUsers() []*model.User
	GetUserByID(id uint) *model.User
	GetUserByName(name string) *model.User
	DeleteUserByID(id uint) error
	UpdateUser(user *model.User)
	CreateUser(user *model.User) error
}

// The UserAPI provides handlers for managing users.
type UserAPI struct {
	DB               UserDatabase
	PasswordStrength int
	NotifyDeleted    func(uint)
}

// GetUsers returns all the users
// swagger:operation GET /user user getUsers
//
// Return all users.
//
// ---
// produces: [application/json]
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
func (a *UserAPI) GetUsers(ctx *gin.Context) {
	users := a.DB.GetUsers()

	var resp []*model.UserExternal
	for _, user := range users {
		resp = append(resp, toExternal(user))
	}

	ctx.JSON(200, resp)
}

// GetCurrentUser returns the current user
// swagger:operation GET /current/user user currentUser
//
// Return the current user.
//
// ---
// produces: [application/json]
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
func (a *UserAPI) GetCurrentUser(ctx *gin.Context) {
	user := a.DB.GetUserByID(auth.GetUserID(ctx))
	ctx.JSON(200, toExternal(user))
}

// CreateUser creates a user
// swagger:operation POST /user user createUser
//
// Create a user.
//
// ---
// consumes: [application/json]
// produces: [application/json]
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
//     $ref: "#/definitions/UserWithPass"
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
func (a *UserAPI) CreateUser(ctx *gin.Context) {
	user := model.UserExternalWithPass{}
	if err := ctx.Bind(&user); err == nil {
		internal := a.toInternal(&user, []byte{})
		if a.DB.GetUserByName(internal.Name) == nil {
			a.DB.CreateUser(internal)
			ctx.JSON(200, toExternal(internal))
		} else {
			ctx.AbortWithError(400, errors.New("username already exists"))
		}
	}
}

// GetUserByID returns the user by id
// swagger:operation GET /user/{id} user getUser
//
// Get a user.
//
// ---
// consumes: [application/json]
// produces: [application/json]
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
func (a *UserAPI) GetUserByID(ctx *gin.Context) {
	withID(ctx, "id", func(id uint) {
		if user := a.DB.GetUserByID(uint(id)); user != nil {
			ctx.JSON(200, toExternal(user))
		} else {
			ctx.AbortWithError(404, errors.New("user does not exist"))
		}
	})
}

// DeleteUserByID deletes the user by id
// swagger:operation DELETE /user/{id} user deleteUser
//
// Deletes a user.
//
// ---
// produces: [application/json]
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
func (a *UserAPI) DeleteUserByID(ctx *gin.Context) {
	withID(ctx, "id", func(id uint) {
		if user := a.DB.GetUserByID(id); user != nil {
			a.NotifyDeleted(id)
			a.DB.DeleteUserByID(id)
		} else {
			ctx.AbortWithError(404, errors.New("user does not exist"))
		}
	})
}

// ChangePassword changes the password from the current user
// swagger:operation POST /current/user/password user updateCurrentUser
//
// Update the password of the current user.
//
// ---
// consumes: [application/json]
// produces: [application/json]
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
//     $ref: "#/definitions/UserPass"
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
func (a *UserAPI) ChangePassword(ctx *gin.Context) {
	pw := model.UserExternalPass{}
	if err := ctx.Bind(&pw); err == nil {
		user := a.DB.GetUserByID(auth.GetUserID(ctx))
		user.Pass = password.CreatePassword(pw.Pass, a.PasswordStrength)
		a.DB.UpdateUser(user)
	}
}

// UpdateUserByID updates and user by id
// swagger:operation POST /user/{id} user updateUser
//
// Update a user.
//
// ---
// consumes: [application/json]
// produces: [application/json]
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
//     $ref: "#/definitions/UserWithPass"
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
func (a *UserAPI) UpdateUserByID(ctx *gin.Context) {
	withID(ctx, "id", func(id uint) {
		var user *model.UserExternalWithPass
		if err := ctx.Bind(&user); err == nil {
			if oldUser := a.DB.GetUserByID(id); oldUser != nil {
				internal := a.toInternal(user, oldUser.Pass)
				internal.ID = id
				a.DB.UpdateUser(internal)
				ctx.JSON(200, toExternal(internal))
			} else {
				ctx.AbortWithError(404, errors.New("user does not exist"))
			}
		}
	})
}

func (a *UserAPI) toInternal(response *model.UserExternalWithPass, pw []byte) *model.User {
	user := &model.User{
		Name:  response.Name,
		Admin: response.Admin,
	}
	if response.Pass != "" {
		user.Pass = password.CreatePassword(response.Pass, a.PasswordStrength)
	} else {
		user.Pass = pw
	}
	return user
}

func toExternal(internal *model.User) *model.UserExternal {
	return &model.UserExternal{
		Name:  internal.Name,
		Admin: internal.Admin,
		ID:    internal.ID,
	}
}
