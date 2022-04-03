package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(uri string) *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(fmt.Errorf("Error when connecting to DB: %v\n", err.Error()))
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		panic(fmt.Errorf("Error when pinging to DB: %v\n", err.Error()))
	}
	return client
}
