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

func collection() (*mongo.Collection, error) {
	db, err := db.Get()
	if err != nil {
		return nil, err
	}

	collection := db.Collection("tokens")

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

	return db.Collection("tokens"), nil
}

// Save agrega un token a la base de datos
func insert(token *Token) (*Token, error) {
	collection, err := collection()
	if err != nil {
		db.HandleError(err)
		return nil, err
	}

	_, err = collection.InsertOne(context.Background(), token)
	if err != nil {
		db.HandleError(err)
		return nil, err
	}

	return token, nil
}

// Save agrega un token a la base de datos
func update(token *Token) (*Token, error) {
	collection, err := collection()
	if err != nil {
		db.HandleError(err)
		return nil, err
	}

	doc, err := bson.NewDocumentEncoder().EncodeDocument(token)
	if err != nil {
		db.HandleError(err)
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
		db.HandleError(err)
		return nil, err
	}

	return token, nil
}

func findByID(tokenID string) (*Token, error) {
	_id, err := getID(tokenID)
	if err != nil {
		return nil, errors.Unauthorized
	}

	collection, err := collection()
	if err != nil {
		db.HandleError(err)
		return nil, err
	}

	token := &Token{}
	filter := bson.NewDocument(bson.EC.ObjectID("_id", *_id))
	err = collection.FindOne(context.Background(), filter).Decode(token)
	if err != nil {
		db.HandleError(err)
		if err == mongo.ErrNoDocuments {
			return nil, errors.Unauthorized
		}
		return nil, err
	}

	return token, nil
}

func findByUserID(tokenID string) (*Token, error) {
	_id, err := getID(tokenID)
	if err != nil {
		return nil, errors.Unauthorized
	}

	collection, err := collection()
	if err != nil {
		db.HandleError(err)
		return nil, err
	}

	token := &Token{}

	filter := bson.NewDocument(
		bson.EC.String("userId", _id.Hex()),
		bson.EC.Boolean("enabled", true),
	)
	err = collection.FindOne(context.Background(), filter).Decode(token)
	if err != nil {
		db.HandleError(err)
		if err == mongo.ErrNoDocuments {
			return nil, errors.Unauthorized
		}
		return nil, err
	}

	return token, nil
}

func delete(tokenID string) error {
	token, err := findByID(tokenID)
	if err != nil {
		db.HandleError(err)
		return err
	}

	token.Enabled = false
	_, err = update(token)

	if err != nil {
		db.HandleError(err)
		return err
	}

	return nil
}

func getID(ID string) (*objectid.ObjectID, error) {
	_id, err := objectid.FromHex(ID)
	if err != nil {
		return nil, ErrID
	}
	return &_id, nil
}
