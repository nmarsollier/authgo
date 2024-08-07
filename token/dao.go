package token

import (
	"context"

	"github.com/nmarsollier/authgo/tools/app_errors"
	"github.com/nmarsollier/authgo/tools/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection db.MongoCollection

func NewProps(collection db.MongoCollection) TokenProps {
	return TokenProps{
		Collection: collection,
	}
}

type TokenProps struct {
	Collection db.MongoCollection
}

func dbCollection(props ...interface{}) (db.MongoCollection, error) {
	for _, o := range props {
		if ti, ok := o.(TokenProps); ok {
			return ti.Collection, nil
		}
	}

	if collection != nil {
		return collection, nil
	}

	database, err := db.Get()
	if err != nil {
		return nil, err
	}

	_collection := database.Collection("tokens")

	_, err = _collection.Indexes().CreateOne(
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

	collection = db.NewMongoCollection(_collection)

	return collection, nil
}

// insert crea un nuevo token y lo almacena en la db
func insert(userID primitive.ObjectID, props ...interface{}) (*Token, error) {
	collection, err := dbCollection(props...)
	if err != nil {
		return nil, err
	}

	token := newToken(userID)

	_, err = collection.InsertOne(context.Background(), token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// findByID busca un token en la db
func findByID(tokenID string, props ...interface{}) (*Token, error) {
	collection, err := dbCollection(props...)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(tokenID)
	if err != nil {
		return nil, app_errors.Unauthorized
	}

	token := &Token{}
	filter := bson.M{"_id": _id}

	if err = collection.FindOne(context.Background(), filter, token); err != nil {
		return nil, err
	}

	return token, nil
}

// delete como deshabilitado un token
func delete(tokenID primitive.ObjectID, props ...interface{}) error {
	collection, err := dbCollection(props...)
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(context.Background(),
		bson.M{"_id": tokenID},
		bson.M{"$set": bson.M{
			"enabled": false,
		}},
	)

	return err
}
