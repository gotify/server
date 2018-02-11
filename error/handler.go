package error

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"net/http"
	"strings"
	"unicode"

	"gopkg.in/go-playground/validator.v8"
)

type errorWrapper struct {
	Error            string `json:"error"`
	ErrorCode        int    `json:"errorCode"`
	ErrorDescription string `json:"errorDescription"`
}

// Handler creates a gin middleware for handling errors.
func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				switch e.Type {
				case gin.ErrorTypeBind:
					errs := e.Err.(validator.ValidationErrors)
					var stringErrors []string
					for _, err := range errs {
						stringErrors = append(stringErrors, validationErrorToText(err))
					}
					writeError(c, strings.Join(stringErrors, "; "))
				default:
					writeError(c, e.Err.Error())
				}
			}
		}
	}
}

func validationErrorToText(e *validator.FieldError) string {
	runes := []rune(e.Field)
	runes[0] = unicode.ToLower(runes[0])
	fieldName := string(runes)
	switch e.Tag {
	case "required":
		return fmt.Sprintf("Field '%s' is required", fieldName)
	}
	return fmt.Sprintf("Field '%s' is not valid", fieldName)
}

func writeError(ctx *gin.Context, errString string) {
	status := http.StatusBadRequest
	if ctx.Writer.Status() != http.StatusOK {
		status = ctx.Writer.Status()
	}
	ctx.JSON(status, &errorWrapper{Error: http.StatusText(status), ErrorCode: status, ErrorDescription: errString})
}
