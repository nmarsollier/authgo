package user

import (
	"context"
	"log"
	"time"

	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define mongo Collection
var collection *mongo.Collection

func dbCollection() (*mongo.Collection, error) {
	if collection != nil {
		return collection, nil
	}

	database, err := db.Get()
	if err != nil {
		return nil, err
	}

	col := database.Collection("users")

	_, err = col.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.M{"login": ""},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		log.Output(1, err.Error())
	}

	collection = col
	return collection, nil
}

func insert(user *User) (*User, error) {
	if err := user.ValidateSchema(); err != nil {
		return nil, err
	}

	var collection, err = dbCollection()
	if err != nil {
		return nil, err
	}

	if _, err := collection.InsertOne(context.Background(), user); err != nil {
		return nil, err
	}

	return user, nil
}

func update(user *User) (*User, error) {
	if err := user.ValidateSchema(); err != nil {
		return nil, err
	}

	var collection, err = dbCollection()
	if err != nil {
		return nil, err
	}

	user.Updated = time.Now()

	_, err = collection.UpdateOne(context.Background(),
		bson.M{"_id": user.ID},
		bson.M{
			"$set": bson.M{
				"password":    user.Password,
				"name":        user.Name,
				"enabled":     user.Enabled,
				"updated":     user.Updated,
				"permissions": user.Permissions,
			},
		},
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindAll devuelve todos los usuarios
func findAll() ([]*User, error) {
	var collection, err = dbCollection()
	if err != nil {
		return nil, err
	}

	filter := bson.D{}
	cur, err := collection.Find(context.Background(), filter, nil)
	defer cur.Close(context.Background())

	if err != nil {
		return nil, err
	}

	users := []*User{}
	for cur.Next(context.Background()) {
		user := &User{}
		if err := cur.Decode(user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// FindByID lee un usuario desde la db
func findByID(userID string) (*User, error) {
	var collection, err = dbCollection()
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.ErrID
	}

	user := &User{}
	filter := bson.M{"_id": _id}
	if err = collection.FindOne(context.Background(), filter).Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}

// FindByLogin lee un usuario desde la db
func findByLogin(login string) (*User, error) {
	var collection, err = dbCollection()
	if err != nil {
		return nil, err
	}

	user := &User{}
	filter := bson.M{"login": login}
	err = collection.FindOne(context.Background(), filter).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrLogin
		}
		return nil, err
	}

	return user, nil
}
