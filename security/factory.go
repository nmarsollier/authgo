package security

import (
	"context"

	"github.com/nmarsollier/authgo/tools/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collectionInstance *mongo.Collection

// newDao es interno, solo se puede usar en este modulo
func getCollection() (*mongo.Collection, error) {
	if collectionInstance != nil {
		return collectionInstance, nil
	}

	database, err := db.Get()
	if err != nil {
		return nil, err
	}

	collectionInstance = database.Collection("tokens")

	_, err = collectionInstance.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{
				"userId": 1, // index in ascending order
			}, Options: nil,
		},
	)
	if err != nil {
		return nil, err
	}

	return collectionInstance, nil
}

// public DaoInstance allow us to mock daos
var DaoInstance Dao

func getDao() Dao {
	if DaoInstance != nil {
		return DaoInstance
	}

	DaoInstance = newDao()
	return DaoInstance
}
