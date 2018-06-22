package security

import (
	"context"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/errors"
)

// Dao es la interfaz con los método expuestos por este dao
type Dao interface {
	Create(userID objectid.ObjectID) (*Token, error)
	FindByID(tokenID string) (*Token, error)
	Delete(tokenID objectid.ObjectID) error
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
			Keys: bson.NewDocument(
				bson.EC.String("userId", ""),
			),
			Options: bson.NewDocument(),
		},
	)
	if err != nil {
		return nil, err
	}

	return daoStruct{
		collection: db.WrapCollection(collection),
	}, nil
}

// MockedDao con fines de testing para mockear db.collection
func MockedDao(coll db.Collection) Dao {
	return daoStruct{
		collection: coll,
	}
}

// El repositorio
type daoStruct struct {
	collection db.Collection
}

// Create crea un nuevo token y lo almacena en la db
func (d daoStruct) Create(userID objectid.ObjectID) (*Token, error) {
	token := newToken(userID)

	_, err := d.collection.InsertOne(context.Background(), token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// Find busca un token en la db
func (d daoStruct) FindByID(tokenID string) (*Token, error) {
	_id, err := objectid.FromHex(tokenID)
	if err != nil {
		return nil, errors.Unauthorized
	}

	token := &Token{}
	filter := bson.NewDocument(bson.EC.ObjectID("_id", _id))

	if err = d.collection.FindOne(context.Background(), filter).Decode(token); err != nil {
		return nil, err
	}

	return token, nil
}

// Delete como deshabilitado un token
func (d daoStruct) Delete(tokenID objectid.ObjectID) error {
	_, err := d.collection.UpdateOne(context.Background(),
		bson.NewDocument(bson.EC.ObjectID("_id", tokenID)),
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				bson.EC.Boolean("enabled", false),
			),
		))

	return err
}
