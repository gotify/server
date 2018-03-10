package api

import (
	"github.com/gin-gonic/gin"
	"errors"
	"strconv"
)

func withID(ctx *gin.Context, name string, f func(id uint)) {
	if id, err := strconv.ParseUint(ctx.Param(name), 10, 32); err == nil {
		f(uint(id));
	} else {
		ctx.AbortWithError(400, errors.New("invalid id"))
	}
}
