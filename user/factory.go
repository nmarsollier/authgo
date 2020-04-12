package user

import (
	"context"
	"log"

	"github.com/nmarsollier/authgo/security"
	"github.com/nmarsollier/authgo/tools/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define mongo Collection
var collectionInstance *mongo.Collection

func getCollection() (*mongo.Collection, error) {
	if collectionInstance != nil {
		return collectionInstance, nil
	}

	database, err := db.Get()
	if err != nil {
		return nil, err
	}

	col := database.Collection("users")

	_, err = col.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.M{"login": ""},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		log.Output(1, err.Error())
	}

	collectionInstance = col
	return collectionInstance, nil
}

// DaoInstance allows to mock dao
var DaoInstance Dao

func getDao() Dao {
	if DaoInstance != nil {
		return DaoInstance
	}

	DaoInstance = newDao()
	return DaoInstance
}

// SecServiceInstance allows mock service
var SecServiceInstance security.Service

func getService() security.Service {
	if SecServiceInstance != nil {
		return SecServiceInstance
	}

	SecServiceInstance = security.NewService()
	return SecServiceInstance
}
