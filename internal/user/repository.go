package user

import (
	"context"
	"time"

	"github.com/nmarsollier/authgo/internal/common/log"
	"github.com/nmarsollier/authgo/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func insert(
	log log.LogRusEntry,
	user *User,
) (*User, error) {
	if err := user.validateSchema(); err != nil {
		log.Error(err)
		return nil, err
	}

	if _, err := db.UserCollection().InsertOne(context.Background(), user); err != nil {
		log.Error(err)
		return nil, err
	}

	return user, nil
}

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

func update(
	log log.LogRusEntry,
	user *User,
) (*User, error) {

	if err := user.validateSchema(); err != nil {
		log.Error(err)
		return nil, err
	}

	user.Updated = time.Now()

	_, err := db.UserCollection().UpdateOne(context.Background(),
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
		nil,
	)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return user, nil
}

func findAll(
	log log.LogRusEntry,
) ([]*User, error) {
	filter := bson.D{}
	cur, err := db.UserCollection().Find(context.Background(), filter)
	if err != nil {
		log.Error(err)
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

func findByID(
	log log.LogRusEntry,
	userID string,
) (*User, error) {
	_id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Error(err)
		return nil, ErrID
	}

	user := &User{}
	filter := DbUserIdFilter{ID: _id}
	if err = db.UserCollection().FindOne(context.Background(), filter, user); err != nil {
		log.Error(err)
		return nil, err
	}

	return user, nil
}

type DbUserIdFilter struct {
	ID primitive.ObjectID `bson:"_id"`
}

type DbUserLoginFilter struct {
	Login string `bson:"login"`
}

func findByLogin(
	log log.LogRusEntry,
	login string,
) (*User, error) {
	user := &User{}
	filter := DbUserLoginFilter{Login: login}
	err := db.UserCollection().FindOne(context.Background(), filter, user)
	if err != nil {
		log.Error(err)
		if err == mongo.ErrNoDocuments {
			return nil, ErrLogin
		}
		return nil, err
	}

	return user, nil
}
