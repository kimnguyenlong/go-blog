package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Comment struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ParentID primitive.ObjectID `bson:"parent_id,omitempty" json:"parent_id"`
	UserID   primitive.ObjectID `bson:"user_id,omitempty" json:"user_id"`
	PostID   primitive.ObjectID `bson:"post_id,omitempty" json:"post_id"`
	Content  string             `bson:"content,omitempty" json:"content"`
	Created  int64              `bson:"created,omitempty" json:"created"`
	Updated  int64              `bson:"updated,omitempty" json:"updated"`
}
