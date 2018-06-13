package token

import (
	"context"
	"fmt"
	"strings"

	validator "gopkg.in/go-playground/validator.v8"

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
		fmt.Print(err.Error())
	}

	return db.Collection("tokens"), nil
}

// Save agrega un token a la base de datos
func insert(token *Token) (*Token, error) {
	if err := validateSchema(token); err != nil {
		return nil, err
	}

	collection, err := collection()
	if err != nil {
		db.HandleError(err)
		return nil, err
	}

	res, err := collection.InsertOne(context.Background(), token)
	if err != nil {
		db.HandleError(err)
		return nil, err
	}

	token.setID(res.InsertedID.(objectid.ObjectID))

	return token, nil
}

// Save agrega un token a la base de datos
func update(token *Token) (*Token, error) {
	if err := validateSchema(token); err != nil {
		return nil, err
	}

	collection, err := collection()
	if err != nil {
		db.HandleError(err)
		return nil, err
	}

	_id, err := objectid.FromHex(token.id())
	if err != nil {
		return nil, err
	}

	doc, err := bson.NewDocumentEncoder().EncodeDocument(token)
	if err != nil {
		db.HandleError(err)
		return nil, err
	}

	_, err = collection.UpdateOne(context.Background(),
		bson.NewDocument(bson.EC.ObjectID("_id", _id)),
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

// Save agrega un token a la base de datos
func save(token *Token) (*Token, error) {
	if len(token.id()) > 0 {
		return update(token)
	}
	return insert(token)
}

func validateSchema(token *Token) error {
	token.UserID = strings.TrimSpace(token.UserID)

	result := make(validator.ValidationErrors)

	if len(token.id()) > 0 {
		if _, err := objectid.FromHex(token.id()); err != nil {
			result["id"] = &validator.FieldError{
				Field: "id",
				Tag:   "Invalid",
			}
		}
	}
	if len(token.UserID) == 0 {
		result["userId"] = &validator.FieldError{
			Field: "userId",
			Tag:   "Requerido",
		}
	} else {
		if _, err := objectid.FromHex(token.UserID); err != nil {
			result["userId"] = &validator.FieldError{
				Field: "userId",
				Tag:   "Invalid",
			}
		}
	}

	if len(result) > 0 {
		return result
	}

	return nil
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

	result := bson.NewDocument()
	filter := bson.NewDocument(bson.EC.ObjectID("_id", *_id))
	err = collection.FindOne(context.Background(), filter).Decode(result)
	if err != nil {
		db.HandleError(err)
		if err == mongo.ErrNoDocuments {
			return nil, errors.Unauthorized
		}
		return nil, err
	}

	token := newTokenFromBson(*result)

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

	result := bson.NewDocument()

	filter := bson.NewDocument(
		bson.EC.String("userId", _id.Hex()),
		bson.EC.Boolean("enabled", true),
	)
	err = collection.FindOne(context.Background(), filter).Decode(result)
	if err != nil {
		db.HandleError(err)
		if err == mongo.ErrNoDocuments {
			return nil, errors.Unauthorized
		}
		return nil, err
	}

	token := newTokenFromBson(*result)

	return token, nil
}

func delete(tokenID string) error {
	token, err := findByID(tokenID)
	if err != nil {
		db.HandleError(err)
		return err
	}

	token.Enabled = false
	_, err = save(token)

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
