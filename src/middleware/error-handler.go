package middleware

import (
	customerror "blog/custom-error"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			var message string = "Something went wrong, please try again"
			var code int = http.StatusInternalServerError
			switch err.(type) {
			case customerror.APIError:
				apiErr := err.(customerror.APIError)
				message = apiErr.Message
				code = apiErr.Code
			default:
				fmt.Printf("Internal server error: %v\n", err)
			}
			ctx.JSON(code, gin.H{
				"error":   true,
				"message": message,
			})
			ctx.Abort()
		}
	}()
	ctx.Next()
}
