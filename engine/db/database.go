package db

import (
	"context"

	"github.com/nmarsollier/authgo/engine/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var database *mongo.Database

// CurrentDatabase the mongo database
func CurrentDatabase() (*mongo.Database, error) {
	if database == nil {
		clientOptions := options.Client().ApplyURI(env.Get().MongoURL)

		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			return nil, err
		}

		database = client.Database("auth")
	}
	return database, nil
}
