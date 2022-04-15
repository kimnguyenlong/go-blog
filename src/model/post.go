package model

import (
	"blog/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var PostSchema = bson.M{
	"bsonType":             "object",
	"title":                "Topics",
	"required":             []string{"user_id", "topics", "title", "content"},
	"additionalProperties": false,
	"properties": bson.M{
		"_id": bson.M{
			"bsonType": "objectId",
		},
		"user_id": bson.M{
			"bsonType":    "objectId",
			"description": "user_id is required",
		},
		"topics": bson.M{
			"bsonType": "array",
			"items": bson.M{
				"bsonType": "string",
			},
			"description": "topics is required",
		},
		"title": bson.M{
			"bsonType":    "string",
			"maxLength":   256,
			"description": "title is required and must be less than 256 characters",
		},
		"description": bson.M{
			"bsonType":    "string",
			"maxLength":   1000,
			"description": "description must be less than 1000 characters",
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

type Post struct {
	Base *Base
}

var post *Post = nil

func NewPost(db *mongo.Database) *Post {
	if post != nil {
		return post
	}
	post = &Post{
		Base: NewBase(db, "posts", PostSchema),
	}
	return post
}

func (post Post) Save(newPost entity.Post) (entity.Post, error) {
	now := time.Now().Unix()
	newPost.Created = now
	newPost.Updated = now
	rs, err := post.Base.Save(newPost)
	if err != nil {
		return entity.Post{}, err
	}
	newPost.ID = rs.InsertedID.(primitive.ObjectID)
	return newPost, nil
}
