package model

import (
	"blog/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var CommentSchema = bson.M{
	"bsonType":             "object",
	"title":                "Comments",
	"required":             []string{"user_id", "post_id", "content"},
	"additionalProperties": false,
	"properties": bson.M{
		"_id": bson.M{
			"bsonType": "objectId",
		},
		"parent_id": bson.M{
			"bsonType": "objectId",
		},
		"user_id": bson.M{
			"bsonType":    "objectId",
			"description": "user_id is required",
		},
		"post_id": bson.M{
			"bsonType":    "objectId",
			"description": "post_id is required",
		},
		"content": bson.M{
			"bsonType":    "string",
			"description": "content is required",
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

type Comment struct {
	Base *Base
}

func NewComment(db *mongo.Database) *Comment {
	return &Comment{
		Base: NewBase(db, "comments", CommentSchema),
	}
}

func (comment Comment) Save(newComment entity.Comment) (entity.Comment, error) {
	now := time.Now().Unix()
	newComment.Created = now
	newComment.Updated = now
	rs, err := comment.Base.Save(newComment)
	if err != nil {
		return entity.Comment{}, err
	}
	newComment.ID = rs.InsertedID.(primitive.ObjectID)
	return newComment, nil
}
