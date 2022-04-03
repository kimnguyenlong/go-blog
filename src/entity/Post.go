package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id,omitempty" json:"user_id"`
	Topics      []string           `bson:"topics,omitempty" json:"topics"`
	Title       string             `bson:"title,omitempty" json:"title"`
	Description string             `bson:"description,omitempty" json:"description"`
	Content     string             `bson:"content,omitempty" json:"content"`
	Created     int64              `bson:"created,omitempty" json:"created"`
	Updated     int64              `bson:"updated,omitempty" json:"updated"`
}
