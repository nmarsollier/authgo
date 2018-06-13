package user

import (
	"context"
	"log"
	"strings"
	"time"

	validator "gopkg.in/go-playground/validator.v8"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/nmarsollier/authgo/tools/db"
)

// UsersCollection obtiene la colección de Usuarios
func collection() (*mongo.Collection, error) {
	database, err := db.Get()
	if err != nil {
		return nil, err
	}

	collection := database.Collection("users")

	_, err = collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.NewDocument(
				bson.EC.String("login", ""),
			),
			Options: bson.NewDocument(
				bson.EC.Boolean("unique", true),
			),
		},
	)
	if err != nil {
		log.Output(1, err.Error())
	}

	return collection, nil
}

func insert(user *User) (*User, error) {
	if err := validateSchema(user); err != nil {
		return nil, err
	}

	collection, err := collection()
	if err != nil {
		return nil, err
	}

	if _, err = collection.InsertOne(context.Background(), user); err != nil {
		return nil, err
	}

	return user, nil
}

func update(user *User) (*User, error) {
	if err := validateSchema(user); err != nil {
		return nil, err
	}

	collection, err := collection()
	if err != nil {
		return nil, err
	}

	user.Updated = time.Now()

	doc, err := bson.NewDocumentEncoder().EncodeDocument(user)
	if err != nil {
		return nil, err
	}

	_, err = collection.UpdateOne(context.Background(),
		bson.NewDocument(doc.LookupElement("_id")),
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				doc.LookupElement("password"),
				doc.LookupElement("name"),
				doc.LookupElement("enabled"),
				doc.LookupElement("updated"),
			),
		))

	if err != nil {
		return nil, err
	}

	return user, nil
}

func validateSchema(user *User) error {
	user.Login = strings.TrimSpace(user.Login)
	user.Name = strings.TrimSpace(user.Name)
	user.Password = strings.TrimSpace(user.Password)

	result := make(validator.ValidationErrors)

	if len(user.Name) == 0 {
		result["name"] = &validator.FieldError{
			Field: "name",
			Tag:   "Requerido",
		}
	}
	if len(user.Password) == 0 {
		result["password"] = &validator.FieldError{
			Field: "password",
			Tag:   "Requerido",
		}
	}
	if len(user.Login) == 0 {
		result["login"] = &validator.FieldError{
			Field: "login",
			Tag:   "Requerido",
		}
	}

	if len(result) > 0 {
		return result
	} else {
		return nil
	}
}

// FindByID lee un usuario desde la db
func findByID(userID string) (*User, error) {
	_id, err := objectid.FromHex(userID)
	if err != nil {
		return nil, ErrID
	}

	collection, err := collection()
	if err != nil {
		return nil, err
	}

	filter := bson.NewDocument(bson.EC.ObjectID("_id", _id))

	user := &User{}
	err = collection.FindOne(context.Background(), filter).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindByLogin lee un usuario desde la db
func findByLogin(login string) (*User, error) {
	collection, collectionError := collection()
	if collectionError != nil {
		return nil, collectionError
	}

	user := &User{}
	filter := bson.NewDocument(bson.EC.String("login", login))
	err := collection.FindOne(context.Background(), filter).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrLogin
		}
		return nil, err
	}

	return user, nil
}

// Delete marca un usuario como borrado en la base de datos
func delete(userID string) error {
	_id, err := getID(userID)
	if err != nil {
		return err
	}

	collection, err := collection()
	if err != nil {
		return err
	}

	user := newUser()
	user.ID = *_id
	user.Enabled = false
	user.Updated = time.Now()

	doc, err := bson.NewDocumentEncoder().EncodeDocument(user)
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(context.Background(),
		bson.NewDocument(doc.LookupElement("_id")),
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				doc.LookupElement("enabled"),
				doc.LookupElement("updated"),
			),
		))

	if err != nil {
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
