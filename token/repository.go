package token

import (
	"context"

	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/errs"
	"github.com/nmarsollier/authgo/tools/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection db.MongoCollection

func dbCollection(ctx ...interface{}) (db.MongoCollection, error) {
	for _, o := range ctx {
		if tc, ok := o.(db.MongoCollection); ok {
			return tc, nil
		}
	}

	if collection != nil {
		return collection, nil
	}

	database, err := db.Get(ctx...)
	if err != nil {
		log.Get(ctx...).Error(err)
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
		log.Get(ctx...).Error(err)
		return nil, err
	}

	collection = db.NewMongoCollection(_collection)

	return collection, nil
}

// insert crea un nuevo token y lo almacena en la db
func insert(userID primitive.ObjectID, ctx ...interface{}) (*Token, error) {
	collection, err := dbCollection(ctx...)
	if err != nil {
		log.Get(ctx...).Error(err)
		return nil, err
	}

	token := newToken(userID)

	_, err = collection.InsertOne(context.Background(), token)
	if err != nil {
		log.Get(ctx...).Error(err)
		return nil, err
	}

	return token, nil
}

// findByID busca un token en la db
func findByID(tokenID string, ctx ...interface{}) (*Token, error) {
	collection, err := dbCollection(ctx...)
	if err != nil {
		log.Get(ctx...).Error(err)
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(tokenID)
	if err != nil {
		log.Get(ctx...).Error(err)
		return nil, errs.Unauthorized
	}

	token := &Token{}
	filter := DbTokenIdFilter{ID: _id}

	if err = collection.FindOne(context.Background(), filter, token); err != nil {
		log.Get(ctx...).Error(err)
		return nil, err
	}

	return token, nil
}

// delete como deshabilitado un token
func delete(tokenID primitive.ObjectID, ctx ...interface{}) error {
	collection, err := dbCollection(ctx...)
	if err != nil {
		log.Get(ctx...).Error(err)
		return err
	}

	_, err = collection.UpdateOne(context.Background(),
		DbTokenIdFilter{ID: tokenID},
		DbDeleteTokenDocument{Set: DbDeleteTokenBody{Enabled: false}},
	)

	return err
}

type DbDeleteTokenBody struct {
	Enabled bool `bson:"enabled"`
}

type DbDeleteTokenDocument struct {
	Set DbDeleteTokenBody `bson:"$set"`
}

type DbTokenIdFilter struct {
	ID primitive.ObjectID `bson:"_id"`
}
