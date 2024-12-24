package token

import (
	"context"

	"github.com/nmarsollier/authgo/engine/db"
	"github.com/nmarsollier/authgo/engine/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewCollection(log log.LogRusEntry) (db.MongoCollection, error) {
	database, err := db.CurrentDatabase()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	collection := database.Collection("tokens")

	_, err = collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{
				"userId": 1, // index in ascending order
			}, Options: nil,
		},
	)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return db.NewMongoCollection(collection), nil
}
