package token

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create crea un nuevo token y lo almacena en la db
func Create(userID primitive.ObjectID, deps ...interface{}) (*Token, error) {
	token, err := insert(userID, deps...)
	if err != nil {
		return nil, err
	}

	cacheAdd(token)

	return token, nil
}
