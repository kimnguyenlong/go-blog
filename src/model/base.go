package model

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Base struct {
	collection *mongo.Collection
}

func NewBase(db *mongo.Database, name string, schema bson.M) *Base {
	validator := bson.M{
		"$jsonSchema": schema,
	}

	// update schema if the collection exists
	updateValidatorCmd := bson.D{
		{Key: "collMod", Value: name},
		{Key: "validator", Value: validator},
	}
	err := db.RunCommand(context.Background(), updateValidatorCmd).Err()
	if err == nil {
		return &Base{
			collection: db.Collection(name),
		}
	}

	// create a new collection if the collection doesn't exist
	opts := options.CreateCollection().SetValidator(validator)
	err = db.CreateCollection(context.Background(), name, opts)
	if err != nil {
		fmt.Printf("Error when creating collection: %v\n", err.Error())
	}
	return &Base{
		collection: db.Collection(name),
	}
}

func (base Base) Save(doc interface{}) (*mongo.InsertOneResult, error) {
	return base.collection.InsertOne(context.Background(), doc)
}

func (base Base) FindOne(filter bson.D, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return base.collection.FindOne(context.Background(), filter, opts...)
}

func (base Base) Find(filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return base.collection.Find(context.Background(), filter, opts...)
}

func (base Base) DeleteMany(filter bson.D, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return base.collection.DeleteMany(context.Background(), filter, opts...)
}

func (base Base) DeleteOne(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return base.collection.DeleteOne(context.Background(), filter, opts...)
}

func (base Base) FindOneAndUpdateOne(filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	return base.collection.FindOneAndUpdate(context.Background(), filter, update, opts...)
}

func (base Base) UpdateOne(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return base.collection.UpdateOne(context.Background(), filter, update, opts...)
}

func (base Base) BulkWrite(models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	return base.collection.BulkWrite(context.Background(), models, opts...)
}
