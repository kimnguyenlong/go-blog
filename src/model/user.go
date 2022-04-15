package model

import (
	"blog/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Base *Base
}

var UserSchema = bson.M{
	"bsonType":             "object",
	"title":                "Users",
	"required":             []string{"username", "email", "password"},
	"additionalProperties": false,
	"properties": bson.M{
		"_id": bson.M{
			"bsonType": "objectId",
		},
		"username": bson.M{
			"bsonType":    "string",
			"maxLength":   256,
			"description": "username is required and must be a string and less than 256 characters",
		},
		"email": bson.M{
			"bsonType":    "string",
			"pattern":     "[a-z0-9]+@[a-z0-9]+",
			"maxLength":   256,
			"description": "email is required and must be a string and matches with pattern [a-z0-9]@[a-z0-9] and less than 256 characters",
		},
		"password": bson.M{
			"bsonType":    "string",
			"maxLength":   256,
			"description": "password must be less than 256 characters",
		},
		"following": bson.M{
			"bsonType": "array",
			"items": bson.M{
				"bsonType": "objectId",
			},
		},
		"followers": bson.M{
			"bsonType": "array",
			"items": bson.M{
				"bsonType": "objectId",
			},
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

var user *User = nil

func NewUser(db *mongo.Database) *User {
	if user != nil {
		return user
	}
	user = &User{
		Base: NewBase(db, "users", UserSchema),
	}
	return user
}

func (user User) Save(newUser entity.User) (entity.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return entity.User{}, err
	}
	now := time.Now().Unix()
	newUser.Password = string(hashedPassword)
	newUser.Created = now
	newUser.Updated = now
	rs, err := user.Base.Save(newUser)
	if err != nil {
		return entity.User{}, err
	}
	newUser.ID = rs.InsertedID.(primitive.ObjectID)
	return newUser, nil
}

func (user User) Follow(uid1 primitive.ObjectID, uid2 primitive.ObjectID, opr string) (entity.User, error) {
	user1Filter := bson.D{{Key: "_id", Value: uid1}}
	user1Update := bson.D{{
		Key: opr,
		Value: bson.D{{
			Key:   "following",
			Value: uid2,
		}},
	}}
	user2Filter := bson.D{{Key: "_id", Value: uid2}}
	user2Update := bson.D{{
		Key: opr,
		Value: bson.D{{
			Key:   "followers",
			Value: uid1,
		}},
	}}
	models := []mongo.WriteModel{
		mongo.NewUpdateOneModel().SetFilter(user1Filter).SetUpdate(user1Update),
		mongo.NewUpdateOneModel().SetFilter(user2Filter).SetUpdate(user2Update),
	}
	opts := options.BulkWrite().SetOrdered(false)

	_, err := user.Base.BulkWrite(models, opts)

	if err != nil {
		return entity.User{}, err
	}

	var updatedUser entity.User

	err = user.Base.FindOne(user1Filter).Decode(&updatedUser)

	return updatedUser, err
}
