package token

import (
	"github.com/nmarsollier/authgo/internal/common/cache"
	"github.com/nmarsollier/authgo/internal/common/errs"
	"github.com/nmarsollier/authgo/internal/common/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var tokenCache = cache.NewCache[Token]()

func Create(log log.LogRusEntry, userID string) (*Token, error) {
	_id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Error(err)
		return nil, errs.Unauthorized
	}

	token, err := insert(log, _id)
	if err != nil {
		return nil, err
	}

	tokenString, err := Encode(token)
	if err != nil {
		return nil, err
	}

	tokenCache.Add(tokenString, token)

	return token, nil
}

func Validate(log log.LogRusEntry, tokenString string) (*Token, error) {
	if token, err := tokenCache.Get(tokenString); err == nil {
		return token, err
	}

	tokenID, _, err := extractPayload(tokenString)
	if err != nil {
		return nil, err
	}

	token, err := findByID(log, tokenID)
	if err != nil || !token.Enabled {
		return nil, errs.Unauthorized
	}

	return token, nil
}

func Invalidate(log log.LogRusEntry, tokenString string) error {
	tokenID, _, err := extractPayload(tokenString)
	if err != nil {
		return errs.Unauthorized
	}

	_id, err := primitive.ObjectIDFromHex(tokenID)
	if err != nil {
		log.Error(err)
		return errs.Unauthorized
	}

	if err = delete(_id); err != nil {
		return err
	}

	tokenCache.Remove(tokenString)

	return nil
}

// Find busca un token en la db
func Find(log log.LogRusEntry, tokenID string) (*Token, error) {
	token, err := findByID(log, tokenID)
	if err != nil {
		return nil, err
	}

	tokenString, err := Encode(token)
	if err != nil {
		return nil, err
	}

	tokenCache.Add(tokenString, token)

	return token, nil
}
