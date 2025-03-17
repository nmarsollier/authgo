package db

import (
	"github.com/nmarsollier/authgo/internal/common/db"
	"github.com/nmarsollier/authgo/internal/common/log"
	"github.com/nmarsollier/authgo/internal/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
)

var currDatabase *mongo.Database

func database() *mongo.Database {
	if currDatabase != nil {
		return currDatabase
	}

	database, err := db.NewDatabase(env.Get().MongoURL, "auth")
	if err != nil {
		log := log.Get(env.Get().FluentURL, env.Get().ServerName)
		log.Fatal(err)
		return nil
	}

	return database
}

func isDbTimeoutError(err error) {
	if err == topology.ErrServerSelectionTimeout {
		currDatabase = nil
		tokenCollection = nil
		userCollection = nil
	}
}
