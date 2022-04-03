package route

import (
	"blog/controller"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func ConfigRouteAuth(router *gin.Engine, dbCon *mongo.Client) {
	authController := controller.NewAuthController(dbCon.Database("blog"))
	auth := router.Group("/api/auth")
	{
		auth.POST("/register/", authController.Register())
		auth.POST("/login/", authController.Login())
	}
}
