package user

import (
	"context"
	"time"

	"github.com/nmarsollier/authgo/tools/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Dao es la interface que expone los servicios de acceso a la DB
type Dao interface {
	Insert(user *User) (*User, error)
	Update(user *User) (*User, error)
	FindAll() ([]*User, error)
	FindByID(userID string) (*User, error)
	FindByLogin(login string) (*User, error)
}

// New dao es interno a este modulo, nadie fuera del modulo tiene acceso
func newDao() Dao {
	return new(daoImpl)
}

type daoImpl struct {
}

func (d daoImpl) Insert(user *User) (*User, error) {
	if err := user.ValidateSchema(); err != nil {
		return nil, err
	}

	var collection, err = getCollection()
	if err != nil {
		return nil, err
	}

	if _, err := collection.InsertOne(context.Background(), user); err != nil {
		return nil, err
	}

	return user, nil
}

func (d daoImpl) Update(user *User) (*User, error) {
	if err := user.ValidateSchema(); err != nil {
		return nil, err
	}

	var collection, err = getCollection()
	if err != nil {
		return nil, err
	}

	user.Updated = time.Now()

	_, err = collection.UpdateOne(context.Background(),
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
func (d daoImpl) FindAll() ([]*User, error) {
	var collection, err = getCollection()
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
func (d daoImpl) FindByID(userID string) (*User, error) {
	var collection, err = getCollection()
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
func (d daoImpl) FindByLogin(login string) (*User, error) {
	var collection, err = getCollection()
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
