package token

import (
	"context"

	"github.com/nmarsollier/commongo/db"
	"github.com/nmarsollier/commongo/errs"
	"github.com/nmarsollier/commongo/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenRepository interface {
	Insert(userID primitive.ObjectID) (*Token, error)
	FindByID(tokenID string) (*Token, error)
	Delete(tokenID primitive.ObjectID) error
}

func NewTokenRepository(
	log log.LogRusEntry,
	collection db.Collection,
) (TokenRepository, error) {
	return &tokenRepository{
		log:        log,
		collection: collection,
	}, nil
}

type tokenRepository struct {
	log        log.LogRusEntry
	collection db.Collection
}

func (r *tokenRepository) Insert(userID primitive.ObjectID) (*Token, error) {
	token := newToken(userID)

	_, err := r.collection.InsertOne(context.Background(), token)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	return token, nil
}

func (r *tokenRepository) FindByID(tokenID string) (*Token, error) {
	_id, err := primitive.ObjectIDFromHex(tokenID)
	if err != nil {
		r.log.Error(err)
		return nil, errs.Unauthorized
	}

	token := &Token{}
	filter := DbTokenIdFilter{ID: _id}

	if err = r.collection.FindOne(context.Background(), filter, token); err != nil {
		r.log.Error(err)
		return nil, err
	}

	return token, nil
}

func (r *tokenRepository) Delete(tokenID primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(context.Background(),
		DbTokenIdFilter{ID: tokenID},
		DbDeleteTokenDocument{Set: DbDeleteTokenBody{Enabled: false}},
		nil,
	)

	return err
}

type DbDeleteTokenBody struct {
	Enabled bool `bson:"enabled"`
}

type DbDeleteTokenDocument struct {
	Set DbDeleteTokenBody `bson:"$set"`
}

type DbTokenIdFilter struct {
	ID primitive.ObjectID `bson:"_id"`
}
