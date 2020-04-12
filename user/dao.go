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

type daoStruct struct {
	collection mongo.Collection
}

// Dao es la interface que exponse los servicios de acceso a la DB
type Dao interface {
	Insert(user *User) (*User, error)
	Update(user *User) (*User, error)
	FindAll() ([]*User, error)
	FindByID(userID string) (*User, error)
	FindByLogin(login string) (*User, error)
}

// New dao es interno a este modulo, nadie fuera del modulo tiene acceso
func newDao() (Dao, error) {
	database, err := db.Get()
	if err != nil {
		return nil, err
	}

	collection := database.Collection("users")

	_, err = collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.M{"login": ""},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		log.Output(1, err.Error())
	}

	return daoStruct{
		collection: *collection,
	}, nil
}

func (d daoStruct) Insert(user *User) (*User, error) {
	if err := user.ValidateSchema(); err != nil {
		return nil, err
	}

	if _, err := d.collection.InsertOne(context.Background(), user); err != nil {
		return nil, err
	}

	return user, nil
}

func (d daoStruct) Update(user *User) (*User, error) {
	if err := user.ValidateSchema(); err != nil {
		return nil, err
	}

	user.Updated = time.Now()

	_, err := d.collection.UpdateOne(context.Background(),
		bson.M{"_id": user.ID},
		bson.M{
			"&set": bson.M{
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
func (d daoStruct) FindAll() ([]*User, error) {
	filter := bson.D{}
	cur, err := d.collection.Find(context.Background(), filter, nil)
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
func (d daoStruct) FindByID(userID string) (*User, error) {
	_id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.ErrID
	}

	user := &User{}
	filter := bson.M{"_id": _id}
	if err = d.collection.FindOne(context.Background(), filter).Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}

// FindByLogin lee un usuario desde la db
func (d daoStruct) FindByLogin(login string) (*User, error) {
	user := &User{}
	filter := bson.M{"login": login}
	err := d.collection.FindOne(context.Background(), filter).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrLogin
		}
		return nil, err
	}

	return user, nil
}
