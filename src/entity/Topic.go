package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Topic struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name    string             `bson:"name,omitempty" json:"name"`
	Created int64              `bson:"created,omitempty" json:"created"`
	Updated int64              `bson:"updated,omitempty" json:"updated"`
}
