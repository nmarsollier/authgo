package db

import (
	"context"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/core/topology"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/nmarsollier/authgo/tools/env"
)

var database *mongo.Database

// Get the mongo database
func Get() (*mongo.Database, error) {
	if database == nil {
		client, err := mongo.NewClientWithOptions(
			env.Get().MongoURL,
			mongo.ClientOpt.ServerSelectionTimeout(time.Second),
		)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		err = client.Connect(context.TODO())
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		database = client.Database("auth2")
	}
	return database, nil
}

// HandleError función a llamar cuando se produce un error de db
func HandleError(err error) {
	if err == topology.ErrServerSelectionTimeout {
		database = nil
	}
}
