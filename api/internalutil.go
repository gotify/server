package api

import (
	"errors"
	"math/bits"
	"strconv"

	"github.com/gin-gonic/gin"
)

func withID(ctx *gin.Context, name string, f func(id uint)) {
	if id, err := strconv.ParseUint(ctx.Param(name), 10, bits.UintSize); err == nil {
		f(uint(id))
	} else {
		ctx.AbortWithError(400, errors.New("invalid id"))
	}
}
