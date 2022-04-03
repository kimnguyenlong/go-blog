package controller

import (
	customerror "blog/custom-error"
	"blog/entity"
	"blog/model"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TopicController interface {
	CreateTopic() gin.HandlerFunc
	GetTopics() gin.HandlerFunc
}

type topicController struct {
	topicModel *model.Topic
}

func NewTopicController(db *mongo.Database) TopicController {
	return &topicController{
		topicModel: model.NewTopic(db),
	}
}

func (topicController topicController) CreateTopic() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var topic entity.Topic
		err := ctx.BindJSON(&topic)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		newTopic, err := topicController.topicModel.Save(topic)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"error":   false,
			"message": "Create new topic successfully",
			"data":    newTopic,
		})
	}
}

func (topicController topicController) GetTopics() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cursor, err := topicController.topicModel.Base.Find(bson.D{})
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		var topics []entity.Topic
		err = cursor.All(context.Background(), &topics)
		if err != nil {
			panic(customerror.NewAPIError(err.Error(), http.StatusBadRequest))
		}
		ctx.JSON(http.StatusOK, gin.H{
			"error": false,
			"data":  topics,
		})
	}
}
