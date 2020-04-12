package db

import (
	"context"
	"log"

	"github.com/nmarsollier/authgo/tools/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
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

// CheckError función a llamar cuando se produce un error de db
func CheckError(err interface{}) {
	if err == topology.ErrServerSelectionTimeout {
		database = nil
	}
}

// IsUniqueKeyError retorna true si el error es de indice único
func IsUniqueKeyError(err error) bool {
	if wErr, ok := err.(mongo.WriteException); ok {
		for i := 0; i < len(wErr.WriteErrors); i++ {
			if wErr.WriteErrors[i].Code == 11000 {
				return true
			}
		}
	}
	return false
}
