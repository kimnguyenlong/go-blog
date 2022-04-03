package route

import (
	"blog/controller"
	"blog/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func ConfigRouteUser(router *gin.Engine, dbCon *mongo.Client) {
	userController := controller.NewUserController(dbCon.Database("blog"))
	users := router.Group("/api/users")
	{
		users.GET("/", userController.GetUsers())
		users.GET("/:id/", userController.GetSingleUser())
		users.PATCH("/:id/", middleware.Authenticate, userController.Update())
	}
}
