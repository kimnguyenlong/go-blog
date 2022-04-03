package middleware

import (
	customerror "blog/custom-error"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Authenticate(ctx *gin.Context) {
	var header struct {
		Authorization string `header:"Authorization" binding:"required"`
	}
	err := ctx.BindHeader(&header)
	if err != nil {
		panic(customerror.NewAPIError("Please provide Authorization header", http.StatusUnauthorized))
	}
	tokenStrings := strings.Split(header.Authorization, " ")
	if len(tokenStrings) != 2 || tokenStrings[0] != "Bearer" {
		panic(customerror.NewAPIError("Please provide a Bearer token", http.StatusUnauthorized))
	}
	token, err := jwt.Parse(tokenStrings[1], func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Method.Alg())
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		panic(customerror.NewAPIError("Invalid token", http.StatusUnauthorized))
	}
	if claims, ok := token.Claims.(jwt.MapClaims); token.Valid && ok {
		exp := claims["exp"].(float64)
		if int64(exp) < time.Now().Unix() {
			panic(customerror.NewAPIError("Token is expired", http.StatusUnauthorized))
		}
		uid := claims["uid"].(string)
		email := claims["email"].(string)
		ctx.Set("uid", uid)
		ctx.Set("email", email)
	} else {
		panic(customerror.NewAPIError("Invalid token", http.StatusUnauthorized))
	}
}
