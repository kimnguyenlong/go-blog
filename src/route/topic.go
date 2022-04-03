package route

import (
	"blog/controller"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func ConfigRouteTopic(router *gin.Engine, dbCon *mongo.Client) {
	topicController := controller.NewTopicController(dbCon.Database("blog"))
	topics := router.Group("/api/topics")
	{
		topics.POST("/", topicController.CreateTopic())
		topics.GET("/", topicController.GetTopics())
	}
}
