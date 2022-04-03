package route

import (
	"blog/controller"
	"blog/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func ConfigRoutePost(router *gin.Engine, dbCon *mongo.Client) {
	postController := controller.NewPostController(dbCon.Database("blog"))
	posts := router.Group("/api/posts")
	{
		posts.POST("/", middleware.Authenticate, postController.CreatePost())
		posts.GET("/", postController.GetPosts())
	}

	post := posts.Group("/:id")
	{
		post.GET("/", postController.GetSinglePost())
	}

	comments := post.Group("/comments")
	{
		comments.GET("/", postController.GetComments())
		comments.POST("/", middleware.Authenticate, postController.CreateComment())
	}
}
