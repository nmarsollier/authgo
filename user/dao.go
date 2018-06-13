package user

import (
	"context"
	"fmt"
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
		db.HandleError(err)
		fmt.Print(err.Error())
	}

	return collection, nil
}

func insert(user *User) (*User, error) {
	if err := validateSchema(user); err != nil {
		return nil, err
	}

	collection, err := collection()
	if err != nil {
		db.HandleError(err)
		return nil, err
	}

	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		db.HandleError(err)
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
		db.HandleError(err)
		return nil, err
	}

	user.Updated = time.Now()

	doc, err := bson.NewDocumentEncoder().EncodeDocument(user)
	if err != nil {
		db.HandleError(err)
		return nil, err
	}

	_, err = collection.UpdateOne(context.Background(),
		bson.NewDocument(bson.EC.ObjectID("_id", user.ID)),
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				doc.LookupElement("password"),
				doc.LookupElement("name"),
				doc.LookupElement("enabled"),
				doc.LookupElement("updated"),
			),
		))

	if err != nil {
		db.HandleError(err)
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
		db.HandleError(err)
		return nil, err
	}

	filter := bson.NewDocument(bson.EC.ObjectID("_id", _id))

	user := &User{}
	err = collection.FindOne(context.Background(), filter).Decode(user)
	if err != nil {
		db.HandleError(err)
		if err == mongo.ErrNoDocuments {
			return nil, ErrID
		}
		return nil, err
	}

	return user, nil
}

// FindByLogin lee un usuario desde la db
func findByLogin(login string) (*User, error) {
	collection, collectionError := collection()
	if collectionError != nil {
		db.HandleError(collectionError)
		return nil, collectionError
	}

	user := &User{}
	filter := bson.NewDocument(bson.EC.String("login", login))
	err := collection.FindOne(context.Background(), filter).Decode(user)
	if err != nil {
		db.HandleError(err)
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
		db.HandleError(err)
		return err
	}

	user := newUser()
	user.Enabled = false
	filter := bson.NewDocument(bson.EC.ObjectID("_id", *_id))

	_, err = collection.UpdateOne(context.Background(), filter, user)
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
