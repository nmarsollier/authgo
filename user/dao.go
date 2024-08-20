package user

import (
	"context"
	"time"

	"github.com/nmarsollier/authgo/log"
	"github.com/nmarsollier/authgo/tools/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection db.MongoCollection

type DbUserUpdateDocumentBody struct {
	Name        string    `bson:"name" validate:"required,min=1,max=100"`
	Password    string    `bson:"password" validate:"required"`
	Permissions []string  `bson:"permissions"`
	Enabled     bool      `bson:"enabled"`
	Updated     time.Time `bson:"updated"`
}

type DbUserUpdateDocument struct {
	Set DbUserUpdateDocumentBody `bson:"$set"`
}

type DbUserIdFilter struct {
	ID primitive.ObjectID `bson:"_id"`
}
type DbUserLoginFilter struct {
	Login string `bson:"login"`
}

func dbCollection(ctx ...interface{}) (db.MongoCollection, error) {
	for _, p := range ctx {
		if coll, ok := p.(db.MongoCollection); ok {
			return coll, nil
		}
	}

	if collection != nil {
		return collection, nil
	}

	database, err := db.Get(ctx...)
	if err != nil {
		log.Get(ctx...).Error(err)
		return nil, err
	}

	col := database.Collection("users")

	_, err = col.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    DbUserLoginFilter{Login: ""},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		log.Get(ctx...).Error(err)
	}

	collection = db.NewMongoCollection(col)
	return collection, nil
}

func insert(user *User, ctx ...interface{}) (*User, error) {
	if err := user.ValidateSchema(); err != nil {
		log.Get(ctx...).Error(err)
		return nil, err
	}

	var collection, err = dbCollection(ctx...)
	if err != nil {
		log.Get(ctx...).Error(err)
		return nil, err
	}

	if _, err := collection.InsertOne(context.Background(), user); err != nil {
		log.Get(ctx...).Error(err)
		return nil, err
	}

	return user, nil
}

func update(user *User, ctx ...interface{}) (*User, error) {
	if err := user.ValidateSchema(); err != nil {
		log.Get(ctx...).Error(err)
		return nil, err
	}

	var collection, err = dbCollection(ctx...)
	if err != nil {
		log.Get(ctx...).Error(err)
		return nil, err
	}

	user.Updated = time.Now()

	_, err = collection.UpdateOne(context.Background(),
		DbUserIdFilter{ID: user.ID},
		DbUserUpdateDocument{
			Set: DbUserUpdateDocumentBody{
				Password:    user.Password,
				Name:        user.Name,
				Enabled:     user.Enabled,
				Updated:     user.Updated,
				Permissions: user.Permissions,
			},
		},
	)

	if err != nil {
		log.Get(ctx...).Error(err)
		return nil, err
	}

	return user, nil
}

// FindAll devuelve todos los usuarios
func findAll(ctx ...interface{}) ([]*User, error) {
	var collection, err = dbCollection(ctx...)
	if err != nil {
		log.Get(ctx...).Error(err)
		return nil, err
	}

	filter := bson.D{}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Get(ctx...).Error(err)
		return nil, err
	}
	defer cur.Close(context.Background())

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
func findByID(userID string, ctx ...interface{}) (*User, error) {
	var collection, err = dbCollection(ctx...)
	if err != nil {
		log.Get(ctx...).Error(err)
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Get(ctx...).Error(err)
		return nil, ErrID
	}

	user := &User{}
	filter := DbUserIdFilter{ID: _id}
	if err = collection.FindOne(context.Background(), filter, user); err != nil {
		log.Get(ctx...).Error(err)
		return nil, err
	}

	return user, nil
}

// FindByLogin lee un usuario desde la db
func findByLogin(login string, ctx ...interface{}) (*User, error) {
	var collection, err = dbCollection(ctx...)
	if err != nil {
		log.Get(ctx...).Error(err)
		return nil, err
	}

	user := &User{}
	filter := DbUserLoginFilter{Login: login}
	err = collection.FindOne(context.Background(), filter, user)
	if err != nil {
		log.Get(ctx...).Error(err)
		if err == mongo.ErrNoDocuments {
			return nil, ErrLogin
		}
		return nil, err
	}

	return user, nil
}
