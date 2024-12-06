package token

import (
	"github.com/nmarsollier/authgo/tools/log"
)

// insert crea un nuevo token y lo almacena en la db
func insert(userID string, deps ...interface{}) (token *Token, err error) {
	collection, err := GetTokenDao(deps...)
	if err != nil {
		log.Get(deps...).Error(err)
		return
	}

	token = newToken(userID)

	err = collection.Insert(token)
	if err != nil {
		log.Get(deps...).Error(err)
	}

	return
}

// findByID busca un token en la db
func findByID(tokenID string, deps ...interface{}) (token *Token, err error) {
	collection, err := GetTokenDao(deps...)
	if err != nil {
		log.Get(deps...).Error(err)
		return
	}

	if token, err = collection.FindById(tokenID); err != nil {
		log.Get(deps...).Error(err)
	}

	return
}

// delete como deshabilitado un token
func delete(tokenID string, deps ...interface{}) (err error) {
	collection, err := GetTokenDao(deps...)
	if err != nil {
		log.Get(deps...).Error(err)
		return
	}

	err = collection.Delete(tokenID)
	return
}

type DbDeleteTokenBody struct {
	Enabled bool `bson:"enabled"`
}

type DbDeleteTokenDocument struct {
	Set DbDeleteTokenBody `bson:"$set"`
}
