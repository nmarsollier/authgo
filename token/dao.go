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

type dao interface {
	collection() (db.Collection, error)
	insert(token *Token) (*Token, error)
	update(token *Token) (*Token, error)
	findByID(tokenID string) (*Token, error)
	findByUserID(tokenID string) (*Token, error)
	delete(tokenID string) error
	getID(ID string) (*objectid.ObjectID, error)
}

func newDao() dao {
	return daoImpl{}
}

func newTestingDao(coll db.Collection) dao {
	return daoImpl{
		dbCollection: coll,
	}
}

func (d daoImpl) collection() (db.Collection, error) {
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
func (d daoImpl) insert(token *Token) (*Token, error) {
	collection, err := d.collection()
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
func (d daoImpl) update(token *Token) (*Token, error) {
	collection, err := d.collection()
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

func (d daoImpl) findByID(tokenID string) (*Token, error) {
	_id, err := d.getID(tokenID)
	if err != nil {
		return nil, errors.Unauthorized
	}

	collection, err := d.collection()
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

func (d daoImpl) findByUserID(tokenID string) (*Token, error) {
	_id, err := d.getID(tokenID)
	if err != nil {
		return nil, errors.Unauthorized
	}

	collection, err := d.collection()
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

func (d daoImpl) delete(tokenID string) error {
	token, err := d.findByID(tokenID)
	if err != nil {
		return err
	}

	token.Enabled = false
	_, err = d.update(token)

	return err
}

func (d daoImpl) getID(ID string) (*objectid.ObjectID, error) {
	_id, err := objectid.FromHex(ID)
	if err != nil {
		return nil, errors.ErrID
	}
	return &_id, nil
}
