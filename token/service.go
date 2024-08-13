package token

import (
	"github.com/golang/glog"
	"github.com/nmarsollier/authgo/rabbit"
	"github.com/nmarsollier/authgo/tools/apperr"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create crea un nuevo token y lo almacena en la db
func Create(userID primitive.ObjectID, ctx ...interface{}) (*Token, error) {
	token, err := insert(userID, ctx...)
	if err != nil {
		return nil, err
	}

	cacheAdd(token)

	return token, nil
}

// Find busca un token en la db
func Find(tokenID string, ctx ...interface{}) (*Token, error) {
	token, err := findByID(tokenID, ctx...)
	if err != nil {
		return nil, err
	}

	cacheAdd(token)

	return token, nil
}

// Validate dado un tokenString devuelve el Token asociado
func Validate(tokenString string, ctx ...interface{}) (*Token, error) {
	if token, err := cacheGet(tokenString); err == nil {
		return token, err
	}

	// Sino validamos el token y lo agregamos al cache
	tokenID, _, err := ExtractPayload(tokenString)
	if err != nil {
		return nil, err
	}

	// Buscamos el token en la db para validarlo
	token, err := Find(tokenID, ctx...)
	if err != nil || !token.Enabled {
		return nil, apperr.Unauthorized
	}

	return token, nil
}

// Invalidate invalida un token
func Invalidate(tokenString string, ctx ...interface{}) error {
	token, err := Validate(tokenString, ctx...)
	if err != nil {
		return apperr.Unauthorized
	}

	if err = delete(token.ID, ctx...); err != nil {
		return err
	}

	cacheRemove(token)

	go func() {
		if err = rabbit.Get(ctx...).SendLogout("bearer " + tokenString); err != nil {
			glog.Info("Rabbit logout no se pudo enviar")
		}
	}()

	return nil
}
