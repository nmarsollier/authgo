package security

import (
	"context"

	"github.com/nmarsollier/authgo/tools/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Dao es la interfaz con los método expuestos por este dao
type Dao interface {
	Create(userID primitive.ObjectID) (*Token, error)
	FindByID(tokenID string) (*Token, error)
	Delete(tokenID primitive.ObjectID) error
}

// newDao es interno, solo se puede usar en este modulo
func newDao() Dao {
	return new(daoImpl)
}

// El repositorio
type daoImpl struct {
}

// Create crea un nuevo token y lo almacena en la db
func (d daoImpl) Create(userID primitive.ObjectID) (*Token, error) {
	collection, err := getCollection()
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

// Find busca un token en la db
func (d daoImpl) FindByID(tokenID string) (*Token, error) {
	collection, err := getCollection()
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(tokenID)
	if err != nil {
		return nil, errors.Unauthorized
	}

	token := &Token{}
	filter := bson.M{"_id": _id}

	if err = collection.FindOne(context.Background(), filter).Decode(token); err != nil {
		return nil, err
	}

	return token, nil
}

// Delete como deshabilitado un token
func (d daoImpl) Delete(tokenID primitive.ObjectID) error {
	collection, err := getCollection()
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
