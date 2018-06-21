package token

import (
	"context"
	"log"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/errors"
)

// Uso en tests solamente

type daoImpl struct {
	dbCollection db.Collection
}

// Dao es la interfaz con los método expuestos por este dao
type Dao interface {
	Collection() (db.Collection, error)
	Insert(token *Token) (*Token, error)
	Update(token *Token) (*Token, error)
	FindByID(tokenID string) (*Token, error)
	FindByUserID(tokenID string) (*Token, error)
	Delete(tokenID string) error
	GetID(ID string) (*objectid.ObjectID, error)
}

func newDao() Dao {
	return daoImpl{}
}

// NewTestingDao nuevo dao con fines de testing
func NewTestingDao(coll db.Collection) Dao {
	return daoImpl{
		dbCollection: coll,
	}
}

func (d daoImpl) Collection() (db.Collection, error) {
	if d.dbCollection != nil {
		return d.dbCollection, nil
	}

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
		log.Output(1, err.Error())
	}

	d.dbCollection = db.WrapCollection(collection)
	return d.dbCollection, nil
}

// Save agrega un token a la base de datos
func (d daoImpl) Insert(token *Token) (*Token, error) {
	collection, err := d.Collection()
	if err != nil {
		return nil, err
	}

	_, err = collection.InsertOne(context.Background(), token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// Save agrega un token a la base de datos
func (d daoImpl) Update(token *Token) (*Token, error) {
	collection, err := d.Collection()
	if err != nil {
		return nil, err
	}

	doc, err := bson.NewDocumentEncoder().EncodeDocument(token)
	if err != nil {
		return nil, err
	}

	_, err = collection.UpdateOne(context.Background(),
		bson.NewDocument(bson.EC.ObjectID("_id", token.ID)),
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				doc.LookupElement("enabled"),
			),
		))

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (d daoImpl) FindByID(tokenID string) (*Token, error) {
	_id, err := d.GetID(tokenID)
	if err != nil {
		return nil, errors.Unauthorized
	}

	collection, err := d.Collection()
	if err != nil {
		return nil, err
	}

	token := &Token{}
	filter := bson.NewDocument(bson.EC.ObjectID("_id", *_id))

	if err = collection.FindOne(context.Background(), filter).Decode(token); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Unauthorized
		}
		return nil, err
	}

	return token, nil
}

func (d daoImpl) FindByUserID(tokenID string) (*Token, error) {
	_id, err := d.GetID(tokenID)
	if err != nil {
		return nil, errors.Unauthorized
	}

	collection, err := d.Collection()
	if err != nil {
		return nil, err
	}

	token := &Token{}

	filter := bson.NewDocument(
		bson.EC.String("userId", _id.Hex()),
		bson.EC.Boolean("enabled", true),
	)

	if err = collection.FindOne(context.Background(), filter).Decode(token); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Unauthorized
		}
		return nil, err
	}

	return token, nil
}

func (d daoImpl) Delete(tokenID string) error {
	token, err := d.FindByID(tokenID)
	if err != nil {
		return err
	}

	token.Enabled = false
	_, err = d.Update(token)

	return err
}

func (d daoImpl) GetID(ID string) (*objectid.ObjectID, error) {
	_id, err := objectid.FromHex(ID)
	if err != nil {
		return nil, errors.ErrID
	}
	return &_id, nil
}
