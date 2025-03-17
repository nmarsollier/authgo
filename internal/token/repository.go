package token

import (
	"context"

	"github.com/nmarsollier/authgo/internal/common/errs"
	"github.com/nmarsollier/authgo/internal/common/log"
	"github.com/nmarsollier/authgo/internal/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func insert(log log.LogRusEntry, userID primitive.ObjectID) (*Token, error) {
	token := newToken(userID)

	_, err := db.TokenCollection().InsertOne(context.Background(), token)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return token, nil
}

func findByID(log log.LogRusEntry, tokenID string) (*Token, error) {
	_id, err := primitive.ObjectIDFromHex(tokenID)
	if err != nil {
		log.Error(err)
		return nil, errs.Unauthorized
	}

	token := &Token{}
	filter := DbTokenIdFilter{ID: _id}

	if err = db.TokenCollection().FindOne(context.Background(), filter, token); err != nil {
		log.Error(err)
		return nil, err
	}

	return token, nil
}

func delete(tokenID primitive.ObjectID) error {
	_, err := db.TokenCollection().UpdateOne(context.Background(),
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
