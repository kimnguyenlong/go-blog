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
		post.PATCH("/", middleware.Authenticate, postController.UpdatePost())
		post.DELETE("/", middleware.Authenticate, postController.DeletePost())
	}

	comments := post.Group("/comments")
	{
		comments.GET("/", postController.GetComments())
		comments.POST("/", middleware.Authenticate, postController.CreateComment())
	}
	comment := comments.Group("/:cid")
	{
		comment.DELETE("/", middleware.Authenticate, postController.DeleteComment())
		comment.PATCH("/", middleware.Authenticate, postController.UpdateComment())
	}
	replies := comment.Group("/replies")
	{
		replies.POST("/", middleware.Authenticate, postController.CreateReply())
		replies.GET("/", postController.GetReplies())
	}
}
