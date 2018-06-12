package user

import (
	"context"
	"fmt"
	"strings"

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

	res, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		db.HandleError(err)
		return nil, err
	}

	user.SetID(res.InsertedID.(objectid.ObjectID))

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

	_id, err := objectid.FromHex(user.ID())
	if err != nil {
		return nil, err
	}

	_, err = collection.UpdateOne(context.Background(),
		bson.NewDocument(bson.EC.ObjectID("_id", _id)),
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				bson.EC.String("password", user.Password),
				bson.EC.String("name", user.Name),
				bson.EC.Boolean("enabled", user.Enabled),
			),
		))

	if err != nil {
		db.HandleError(err)
		return nil, err
	}

	return user, nil
}

func save(user *User) (*User, error) {
	if len(user.ID()) > 0 {
		return update(user)
	}
	return insert(user)
}

func validateSchema(user *User) error {
	user.Login = strings.TrimSpace(user.Login)
	user.Name = strings.TrimSpace(user.Name)
	user.Password = strings.TrimSpace(user.Password)

	result := make(validator.ValidationErrors)

	if len(user.ID()) > 0 {
		if _, err := objectid.FromHex(user.ID()); err != nil {
			result["id"] = &validator.FieldError{
				Field: "id",
				Tag:   "Invalid",
			}
		}
	}
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

	result := bson.NewDocument()
	filter := bson.NewDocument(bson.EC.ObjectID("_id", _id))
	err = collection.FindOne(context.Background(), filter).Decode(result)
	if err != nil {
		db.HandleError(err)
		if err == mongo.ErrNoDocuments {
			return nil, ErrID
		}
		return nil, err
	}

	user := newUserFromBson(*result)

	return user, nil
}

// FindByLogin lee un usuario desde la db
func findByLogin(login string) (*User, error) {
	collection, collectionError := collection()
	if collectionError != nil {
		db.HandleError(collectionError)
		return nil, collectionError
	}

	result := bson.NewDocument()
	filter := bson.NewDocument(bson.EC.String("login", login))
	err := collection.FindOne(context.Background(), filter).Decode(result)
	if err != nil {
		db.HandleError(err)
		if err == mongo.ErrNoDocuments {
			return nil, ErrLogin
		}
		return nil, err
	}

	user := newUserFromBson(*result)

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
