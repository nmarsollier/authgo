package db

import (
	"context"
	"log"

	"github.com/nmarsollier/authgo/tools/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var database *mongo.Database

// Get the mongo database
func Get() (*mongo.Database, error) {
	if database == nil {
		clientOptions := options.Client().ApplyURI(env.Get().MongoURL)

		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		database = client.Database("auth2")
	}
	return database, nil
}
