package user

import (
	"context"
	"time"

	"github.com/nmarsollier/authgo/engine/db"
	"github.com/nmarsollier/authgo/engine/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Insert(usr *User) (*User, error)
	Update(usr *User) (*User, error)
	FindAll() ([]*User, error)
	FindByID(userID string) (*User, error)
	FindByLogin(login string) (*User, error)
}

func NewUserRepository(
	log log.LogRusEntry,
	collection db.MongoCollection,
) (UserRepository, error) {
	return &userRepository{
		log:        log,
		collection: collection,
	}, nil
}

type userRepository struct {
	log        log.LogRusEntry
	collection db.MongoCollection
}

func (r *userRepository) Insert(user *User) (*User, error) {
	if err := user.validateSchema(); err != nil {
		r.log.Error(err)
		return nil, err
	}

	if _, err := r.collection.InsertOne(context.Background(), user); err != nil {
		r.log.Error(err)
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Update(user *User) (*User, error) {
	if err := user.validateSchema(); err != nil {
		r.log.Error(err)
		return nil, err
	}

	user.Updated = time.Now()

	_, err := r.collection.UpdateOne(context.Background(),
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
		r.log.Error(err)
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

func (r *userRepository) FindAll() ([]*User, error) {
	filter := bson.D{}
	cur, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		r.log.Error(err)
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

func (r *userRepository) FindByID(userID string) (*User, error) {
	_id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		r.log.Error(err)
		return nil, ErrID
	}

	user := &User{}
	filter := DbUserIdFilter{ID: _id}
	if err = r.collection.FindOne(context.Background(), filter, user); err != nil {
		r.log.Error(err)
		return nil, err
	}

	return user, nil
}

type DbUserIdFilter struct {
	ID primitive.ObjectID `bson:"_id"`
}

func (r *userRepository) FindByLogin(login string) (*User, error) {
	user := &User{}
	filter := DbUserLoginFilter{Login: login}
	err := r.collection.FindOne(context.Background(), filter, user)
	if err != nil {
		r.log.Error(err)
		if err == mongo.ErrNoDocuments {
			return nil, ErrLogin
		}
		return nil, err
	}

	return user, nil
}

type DbUserLoginFilter struct {
	Login string `bson:"login"`
}
