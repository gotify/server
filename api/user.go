package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/v3/auth"
	"github.com/gotify/server/v3/auth/password"
	"github.com/gotify/server/v3/database"
	"github.com/gotify/server/v3/model"
)

var errCannotDeleteLastAdmin = errors.New("cannot delete last admin")

// UserChangeNotifier notifies listeners for user changes.
type UserChangeNotifier struct {
	userDeletedCallbacks []func(uid uint) error
	userAddedCallbacks   []func(uid uint) error
}

// OnUserDeleted is called on user deletion.
func (c *UserChangeNotifier) OnUserDeleted(cb func(uid uint) error) {
	c.userDeletedCallbacks = append(c.userDeletedCallbacks, cb)
}

// OnUserAdded is called on user creation.
func (c *UserChangeNotifier) OnUserAdded(cb func(uid uint) error) {
	c.userAddedCallbacks = append(c.userAddedCallbacks, cb)
}

func (c *UserChangeNotifier) fireUserDeleted(uid uint) error {
	for _, cb := range c.userDeletedCallbacks {
		if err := cb(uid); err != nil {
			return err
		}
	}
	return nil
}

func (c *UserChangeNotifier) fireUserAdded(uid uint) error {
	for _, cb := range c.userAddedCallbacks {
		if err := cb(uid); err != nil {
			return err
		}
	}
	return nil
}

// The UserAPI provides handlers for managing users.
type UserAPI struct {
	DB                 *database.GormDatabase
	PasswordStrength   int
	UserChangeNotifier *UserChangeNotifier
	Registration       bool
}

// GetUsers returns all the users
// swagger:operation GET /user user getUsers
//
// Return all users.
//
// Requires elevated authentication.
//
//	---
//	produces: [application/json]
//	security: [clientTokenAuthorizationHeader: [], clientTokenHeader: [], clientTokenQuery: [], basicAuth: []]
//	responses:
//	  200:
//	    description: Ok
//	    schema:
//	      type: array
//	      items:
//	        $ref: "#/definitions/User"
//	  401:
//	    description: Unauthorized
//	    schema:
//	        $ref: "#/definitions/Error"
//	  403:
//	    description: Forbidden
//	    schema:
//	        $ref: "#/definitions/Error"
func (a *UserAPI) GetUsers(ctx *gin.Context) {
	users, err := a.DB.GetUsers()
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}
	var resp []*model.UserExternal
	for _, user := range users {
		resp = append(resp, toExternalUser(user))
	}

	ctx.JSON(200, resp)
}

// GetCurrentUser returns the current user
// swagger:operation GET /current/user user currentUser
//
// Return the current user.
//
// Requires elevated authentication.
//
//	---
//	produces: [application/json]
//	security: [clientTokenAuthorizationHeader: [], clientTokenHeader: [], clientTokenQuery: [], basicAuth: []]
//	responses:
//	  200:
//	    description: Ok
//	    schema:
//	        $ref: "#/definitions/User"
//	  401:
//	    description: Unauthorized
//	    schema:
//	        $ref: "#/definitions/Error"
//	  403:
//	    description: Forbidden
//	    schema:
//	        $ref: "#/definitions/Error"
func (a *UserAPI) GetCurrentUser(ctx *gin.Context) {
	user, err := a.DB.GetUserByID(auth.GetUserID(ctx))
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}
	result := &model.CurrentUserExternal{
		ID:        user.ID,
		Name:      user.Name,
		Admin:     user.Admin,
		CreatedAt: user.CreatedAt,
	}
	client := auth.GetClient(ctx)
	if client != nil {
		result.ClientID = client.ID
		if client.ElevatedUntil != nil && time.Now().Before(*client.ElevatedUntil) {
			result.ElevatedUntil = client.ElevatedUntil
		}
	}
	ctx.JSON(200, result)
}

// CreateUser create a user.
// swagger:operation POST /user user createUser
//
// Create a user.
//
// With enabled registration: non admin users can be created without authentication.
// With disabled registrations: users can only be created by admin users.
//
// Requires elevated authentication.
//
//	---
//	consumes: [application/json]
//	produces: [application/json]
//	security: [clientTokenAuthorizationHeader: [], clientTokenHeader: [], clientTokenQuery: [], basicAuth: []]
//	parameters:
//	- name: body
//	  in: body
//	  description: the user to add
//	  required: true
//	  schema:
//	    $ref: "#/definitions/CreateUserExternal"
//	responses:
//	  200:
//	    description: Ok
//	    schema:
//	        $ref: "#/definitions/User"
//	  400:
//	    description: Bad Request
//	    schema:
//	        $ref: "#/definitions/Error"
//	  401:
//	    description: Unauthorized
//	    schema:
//	        $ref: "#/definitions/Error"
//	  403:
//	    description: Forbidden
//	    schema:
//	        $ref: "#/definitions/Error"
func (a *UserAPI) CreateUser(ctx *gin.Context) {
	user := model.CreateUserExternal{}
	if err := ctx.Bind(&user); err == nil {
		if err := password.ValidateNewPassword(user.Pass); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		pw, err := password.CreatePassword(user.Pass, a.PasswordStrength)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to prepare password: %s", err))
			return
		}
		internal := &model.User{
			Name:  user.Name,
			Admin: user.Admin,
			Pass:  pw,
		}
		existingUser, err := a.DB.GetUserByName(internal.Name)
		if success := successOrAbort(ctx, 500, err); !success {
			return
		}

		var requestedBy *model.User
		uid := auth.TryGetUserID(ctx)
		if uid != nil {
			requestedBy, err = a.DB.GetUserByID(*uid)
			if err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not get user: %s", err))
				return
			}
		}

		if requestedBy == nil || !requestedBy.Admin {
			status := http.StatusUnauthorized
			if requestedBy != nil {
				status = http.StatusForbidden
			}
			if !a.Registration {
				ctx.AbortWithError(status, errors.New("you are not allowed to access this api"))
				return
			}
			if internal.Admin {
				ctx.AbortWithError(status, errors.New("you are not allowed to create an admin user"))
				return
			}
		}

		if existingUser == nil {
			if success := successOrAbort(ctx, 500, a.DB.CreateUser(internal)); !success {
				return
			}
			if err := a.UserChangeNotifier.fireUserAdded(internal.ID); err != nil {
				ctx.AbortWithError(500, err)
				return
			}
			ctx.JSON(200, toExternalUser(internal))
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
// Requires elevated authentication.
//
//	---
//	consumes: [application/json]
//	produces: [application/json]
//	security: [clientTokenAuthorizationHeader: [], clientTokenHeader: [], clientTokenQuery: [], basicAuth: []]
//	parameters:
//	- name: id
//	  in: path
//	  description: the user id
//	  required: true
//	  type: integer
//	  format: int64
//	responses:
//	  200:
//	    description: Ok
//	    schema:
//	        $ref: "#/definitions/User"
//	  400:
//	    description: Bad Request
//	    schema:
//	        $ref: "#/definitions/Error"
//	  401:
//	    description: Unauthorized
//	    schema:
//	        $ref: "#/definitions/Error"
//	  403:
//	    description: Forbidden
//	    schema:
//	        $ref: "#/definitions/Error"
//	  404:
//	    description: Not Found
//	    schema:
//	        $ref: "#/definitions/Error"
func (a *UserAPI) GetUserByID(ctx *gin.Context) {
	withID(ctx, "id", func(id uint) {
		user, err := a.DB.GetUserByID(id)
		if success := successOrAbort(ctx, 500, err); !success {
			return
		}
		if user != nil {
			ctx.JSON(200, toExternalUser(user))
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
// Requires elevated authentication.
//
//	---
//	produces: [application/json]
//	security: [clientTokenAuthorizationHeader: [], clientTokenHeader: [], clientTokenQuery: [], basicAuth: []]
//	parameters:
//	- name: id
//	  in: path
//	  description: the user id
//	  required: true
//	  type: integer
//	  format: int64
//	responses:
//	  200:
//	    description: Ok
//	  400:
//	    description: Bad Request
//	    schema:
//	        $ref: "#/definitions/Error"
//	  401:
//	    description: Unauthorized
//	    schema:
//	        $ref: "#/definitions/Error"
//	  403:
//	    description: Forbidden
//	    schema:
//	        $ref: "#/definitions/Error"
//	  404:
//	    description: Not Found
//	    schema:
//	        $ref: "#/definitions/Error"
func (a *UserAPI) DeleteUserByID(ctx *gin.Context) {
	withID(ctx, "id", func(id uint) {
		user, err := a.DB.GetUserByID(id)
		if success := successOrAbort(ctx, 500, err); !success {
			return
		}
		if user != nil {
			for range 3 {
				err = a.DB.Txn(func(txdb *database.GormDatabase) error {
					if err := txdb.DeleteUserByID(id); err != nil {
						return err
					}
					anotherAdmin, err := txdb.GetUsers(&model.User{Admin: true})
					if err != nil {
						return err
					}
					if user.Admin && len(anotherAdmin) == 0 {
						ctx.AbortWithError(400, errCannotDeleteLastAdmin)
						return errCannotDeleteLastAdmin
					}
					return a.UserChangeNotifier.fireUserDeleted(id)
				})
				if err == nil || ctx.IsAborted() {
					return
				}
			}
			successOrAbort(ctx, 500, err)
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
// Requires elevated authentication.
//
//	---
//	consumes: [application/json]
//	produces: [application/json]
//	security: [clientTokenAuthorizationHeader: [], clientTokenHeader: [], clientTokenQuery: [], basicAuth: []]
//	parameters:
//	- name: body
//	  in: body
//	  description: the user
//	  required: true
//	  schema:
//	    $ref: "#/definitions/UserPass"
//	responses:
//	  200:
//	    description: Ok
//	  400:
//	    description: Bad Request
//	    schema:
//	        $ref: "#/definitions/Error"
//	  401:
//	    description: Unauthorized
//	    schema:
//	        $ref: "#/definitions/Error"
//	  403:
//	    description: Forbidden
//	    schema:
//	        $ref: "#/definitions/Error"
func (a *UserAPI) ChangePassword(ctx *gin.Context) {
	pw := model.UserExternalPass{}
	if err := ctx.Bind(&pw); err == nil {
		if err := password.ValidateNewPassword(pw.Pass); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		user, err := a.DB.GetUserByID(auth.GetUserID(ctx))
		if success := successOrAbort(ctx, 500, err); !success {
			return
		}
		pw, err := password.CreatePassword(pw.Pass, a.PasswordStrength)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to prepare password: %s", err))
			return
		}
		user.Pass = pw
		successOrAbort(ctx, 500, a.DB.UpdateUser(user))
	}
}

// UpdateUserByID updates and user by id
// swagger:operation POST /user/{id} user updateUser
//
// Update a user.
//
// Requires elevated authentication.
//
//	---
//	consumes: [application/json]
//	produces: [application/json]
//	security: [clientTokenAuthorizationHeader: [], clientTokenHeader: [], clientTokenQuery: [], basicAuth: []]
//	parameters:
//	- name: id
//	  in: path
//	  description: the user id
//	  required: true
//	  type: integer
//	  format: int64
//	- name: body
//	  in: body
//	  description: the updated user
//	  required: true
//	  schema:
//	    $ref: "#/definitions/UpdateUserExternal"
//	responses:
//	  200:
//	    description: Ok
//	    schema:
//	        $ref: "#/definitions/User"
//	  400:
//	    description: Bad Request
//	    schema:
//	        $ref: "#/definitions/Error"
//	  401:
//	    description: Unauthorized
//	    schema:
//	        $ref: "#/definitions/Error"
//	  403:
//	    description: Forbidden
//	    schema:
//	        $ref: "#/definitions/Error"
//	  404:
//	    description: Not Found
//	    schema:
//	        $ref: "#/definitions/Error"
func (a *UserAPI) UpdateUserByID(ctx *gin.Context) {
	withID(ctx, "id", func(id uint) {
		var updatedUser *model.UpdateUserExternal
		if err := ctx.Bind(&updatedUser); err == nil {
			dbUser, err := a.DB.GetUserByID(id)
			if success := successOrAbort(ctx, 500, err); !success {
				return
			}
			if dbUser != nil {
				dbUserWasAdmin := dbUser.Admin
				dbUser.Name = updatedUser.Name
				dbUser.Admin = updatedUser.Admin

				if updatedUser.Pass != "" {
					if err := password.ValidateNewPassword(updatedUser.Pass); err != nil {
						ctx.AbortWithError(http.StatusBadRequest, err)
						return
					}
					pw, err := password.CreatePassword(updatedUser.Pass, a.PasswordStrength)
					if err != nil {
						ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to prepare password: %s", err))
						return
					}
					dbUser.Pass = pw
				}

				for range 3 {
					err = a.DB.Txn(func(txdb *database.GormDatabase) error {
						if err := txdb.UpdateUser(dbUser); err != nil {
							return err
						}

						anotherAdmin, err := txdb.GetUsers(&model.User{Admin: true})
						if err != nil {
							return err
						}
						if !updatedUser.Admin && dbUserWasAdmin && len(anotherAdmin) == 0 {
							ctx.AbortWithError(400, errCannotDeleteLastAdmin)
							return errCannotDeleteLastAdmin
						}

						return nil
					})

					if ctx.IsAborted() {
						return
					}

					if err == nil {
						ctx.JSON(200, toExternalUser(dbUser))
						return
					}
				}
				ctx.AbortWithError(500, err)
			} else {
				ctx.AbortWithError(404, errors.New("user does not exist"))
			}
		}
	})
}

func toExternalUser(internal *model.User) *model.UserExternal {
	return &model.UserExternal{
		Name:      internal.Name,
		Admin:     internal.Admin,
		ID:        internal.ID,
		CreatedAt: internal.CreatedAt,
	}
}
