package api

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/auth"
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
}

// GetUsers returns all the users
func (a *UserAPI) GetUsers(ctx *gin.Context) {
	users := a.DB.GetUsers()

	var resp []*model.UserExternal
	for _, user := range users {
		resp = append(resp, toExternal(user))
	}

	ctx.JSON(200, resp)
}

// GetCurrentUser returns the current user
func (a *UserAPI) GetCurrentUser(ctx *gin.Context) {
	user := a.DB.GetUserByID(auth.GetUserID(ctx))
	ctx.JSON(200, toExternal(user))
}

// CreateUser creates a user
func (a *UserAPI) CreateUser(ctx *gin.Context) {
	user := model.UserExternalWithPass{}
	if err := ctx.Bind(&user); err == nil {
		if len(user.Pass) == 0 {
			ctx.AbortWithError(400, errors.New("password may not be empty"))
		} else {
			internal := a.toInternal(&user, []byte{})
			if a.DB.GetUserByName(internal.Name) == nil {
				a.DB.CreateUser(internal)
				ctx.JSON(200, toExternal(internal))
			} else {
				ctx.AbortWithError(400, errors.New("username already exists"))
			}
		}
	}
}

// GetUserByID returns the user by id
func (a *UserAPI) GetUserByID(ctx *gin.Context) {
	if id, err := toUInt(ctx.Param("id")); err == nil {
		if user := a.DB.GetUserByID(uint(id)); user != nil {
			ctx.JSON(200, toExternal(user))
		} else {
			ctx.AbortWithError(404, errors.New("user does not exist"))
		}
	} else {
		ctx.AbortWithError(400, errors.New("invalid id"))
	}
}

// DeleteUserByID deletes the user by id
func (a *UserAPI) DeleteUserByID(ctx *gin.Context) {
	if id, err := toUInt(ctx.Param("id")); err == nil {
		if user := a.DB.GetUserByID(id); user != nil {
			a.DB.DeleteUserByID(id)
		} else {
			ctx.AbortWithError(404, errors.New("user does not exist"))
		}
	} else {
		ctx.AbortWithError(400, errors.New("invalid id"))
	}
}

// ChangePassword changes the password from the current user
func (a *UserAPI) ChangePassword(ctx *gin.Context) {
	pw := model.UserExternalPass{}
	if err := ctx.Bind(&pw); err == nil {
		user := a.DB.GetUserByID(auth.GetUserID(ctx))
		user.Pass = auth.CreatePassword(pw.Pass, a.PasswordStrength)
		a.DB.UpdateUser(user)
	}
}

// UpdateUserByID updates and user by id
func (a *UserAPI) UpdateUserByID(ctx *gin.Context) {
	if id, err := toUInt(ctx.Param("id")); err == nil {
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
	} else {
		ctx.AbortWithError(400, errors.New("invalid id"))
	}
}

func toUInt(id string) (uint, error) {
	parsed, err := strconv.ParseUint(id, 10, 32)
	return uint(parsed), err
}

func (a *UserAPI) toInternal(response *model.UserExternalWithPass, pw []byte) *model.User {
	user := &model.User{
		Name:  response.Name,
		Admin: response.Admin,
	}
	if response.Pass != "" {
		user.Pass = auth.CreatePassword(response.Pass, a.PasswordStrength)
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
