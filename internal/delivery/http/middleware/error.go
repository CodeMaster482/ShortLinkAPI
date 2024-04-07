package middleware

import (
	apperror "ShortLinkAPI/pkg/errors"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) == 0 {
			return
		}

		e := ctx.Errors[0].Unwrap()
		apiErr, ok := e.(*apperror.ApiError)

		var err error
		if ok {
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
