package db

import (
	"context"

	"github.com/nmarsollier/authgo/tools/env"
	"github.com/nmarsollier/authgo/tools/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var database *mongo.Database

// Get the mongo database
func Get(ctx ...interface{}) (*mongo.Database, error) {
	if database == nil {
		clientOptions := options.Client().ApplyURI(env.Get().MongoURL)

		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Get(ctx...).Fatal(err)
			return nil, err
		}

		database = client.Database("auth")
	}
	return database, nil
}
