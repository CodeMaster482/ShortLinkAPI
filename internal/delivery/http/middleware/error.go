package middleware

import (
	"errors"

	apperror "github.com/CodeMaster482/ShortLinkAPI/pkg/errors"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) == 0 {
			return
		}

		e := ctx.Errors[0].Unwrap()

		var err error

		var apiErr *apperror.APIError
		if errors.As(e, &apiErr) {
			err = apiErr.Unwrap()
		} else {
			err = apperror.ErrInternalServer
		}

		ctx.JSON(apperror.Errors[err].Code, gin.H{
			"status":  apperror.Errors[err].Code,
			"message": apperror.Errors[err].Message,
		})
	}

	return fn
}
