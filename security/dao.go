package security

import (
	"context"

	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Dao es la interfaz con los método expuestos por este dao
type Dao interface {
	Create(userID primitive.ObjectID) (*Token, error)
	FindByID(tokenID string) (*Token, error)
	Delete(tokenID primitive.ObjectID) error
}

// MockedDao con fines de testing para mockear db.collection
func MockedDao(coll mongo.Collection) Dao {
	return daoStruct{
		collection: coll,
	}
}

// newDao es interno, solo se puede usar en este modulo
func newDao() (Dao, error) {
	database, err := db.Get()
	if err != nil {
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
		return nil, err
	}

	return daoStruct{
		collection: *collection,
	}, nil
}

// El repositorio
type daoStruct struct {
	collection mongo.Collection
}

// Create crea un nuevo token y lo almacena en la db
func (d daoStruct) Create(userID primitive.ObjectID) (*Token, error) {
	token := newToken(userID)

	_, err := d.collection.InsertOne(context.Background(), token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// Find busca un token en la db
func (d daoStruct) FindByID(tokenID string) (*Token, error) {
	_id, err := primitive.ObjectIDFromHex(tokenID)
	if err != nil {
		return nil, errors.Unauthorized
	}

	token := &Token{}
	filter := bson.M{"_id": _id}

	if err = d.collection.FindOne(context.Background(), filter).Decode(token); err != nil {
		return nil, err
	}

	return token, nil
}

// Delete como deshabilitado un token
func (d daoStruct) Delete(tokenID primitive.ObjectID) error {
	_, err := d.collection.UpdateOne(context.Background(),
		bson.M{"_id": tokenID},
		bson.M{"$set": bson.M{
			"enabled": false,
		}},
	)

	return err
}
