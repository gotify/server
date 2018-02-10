package api

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmattheis/memo/auth"
	"github.com/jmattheis/memo/model"
)

type userResponse struct {
	ID    uint   `json:"id"`
	Name  string `binding:"required" json:"name" query:"name" form:"name"`
	Pass  string `json:"pass,omitempty" form:"pass" query:"pass"`
	Admin bool   `json:"admin" form:"admin" query:"admin"`
}

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
	DB UserDatabase
}

// GetUsers returns all the users
func (a *UserAPI) GetUsers(ctx *gin.Context) {
	users := a.DB.GetUsers()

	var resp []*userResponse
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
	user := userResponse{}
	if err := ctx.Bind(&user); err == nil {
		if len(user.Pass) == 0 {
			ctx.AbortWithError(400, errors.New("password may not be empty"))
		} else {
			internal := toInternal(&user, []byte{})
			if a.DB.GetUserByName(internal.Name) == nil {
				a.DB.CreateUser(internal)
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

type userPassword struct {
	Pass string `binding:"required" json:"pass" form:"pass" query:"pass" `
}

// ChangePassword changes the password from the current user
func (a *UserAPI) ChangePassword(ctx *gin.Context) {
	pw := userPassword{}
	if err := ctx.Bind(&pw); err == nil {
		user := a.DB.GetUserByID(auth.GetUserID(ctx))
		user.Pass = auth.CreatePassword(pw.Pass)
		a.DB.UpdateUser(user)
	}
}

// UpdateUserByID updates and user by id
func (a *UserAPI) UpdateUserByID(ctx *gin.Context) {
	if id, err := toUInt(ctx.Param("id")); err == nil {
		var user *userResponse
		if err := ctx.Bind(&user); err == nil {
			if oldUser := a.DB.GetUserByID(id); oldUser != nil {
				internal := toInternal(user, oldUser.Pass)
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

func toInternal(response *userResponse, pw []byte) *model.User {
	user := &model.User{
		Name:  response.Name,
		Admin: response.Admin,
	}
	if response.Pass != "" {
		user.Pass = auth.CreatePassword(response.Pass)
	} else {
		user.Pass = pw
	}
	return user
}

func toExternal(internal *model.User) *userResponse {
	return &userResponse{
		Name:  internal.Name,
		Admin: internal.Admin,
		ID:    internal.ID,
	}
}
