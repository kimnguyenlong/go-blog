package model

import (
	"blog/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Topic struct {
	Base *Base
}

var TopicSchema = bson.M{
	"bsonType":             "object",
	"title":                "Topics",
	"required":             []string{"name"},
	"additionalProperties": false,
	"properties": bson.M{
		"_id": bson.M{
			"bsonType": "objectId",
		},
		"name": bson.M{
			"bsonType":    "string",
			"maxLength":   256,
			"description": "name is required and must be a string and less than 256 characters",
		},
		"created": bson.M{
			"bsonType":    "number",
			"description": "created_at must be a int",
		},
		"updated": bson.M{
			"bsonType":    "number",
			"description": "updated_at must be a int",
		},
	},
}

func NewTopic(db *mongo.Database) *Topic {
	return &Topic{
		Base: NewBase(db, "topics", TopicSchema),
	}
}

func (topic Topic) Save(newTopic entity.Topic) (entity.Topic, error) {
	now := time.Now().Unix()
	newTopic.Created = now
	newTopic.Updated = now
	rs, err := topic.Base.Save(newTopic)
	if err != nil {
		return entity.Topic{}, err
	}
	newTopic.ID = rs.InsertedID.(primitive.ObjectID)
	return newTopic, nil
}
