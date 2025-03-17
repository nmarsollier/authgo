package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewDatabase the mongo database
func NewDatabase(
	mongoUrl string,
	name string,
) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(mongoUrl)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	database := client.Database(name)

	return database, nil
}
