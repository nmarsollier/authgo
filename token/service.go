package token

import (
	"github.com/golang/glog"
	"github.com/nmarsollier/authgo/rabbit"
	"github.com/nmarsollier/authgo/tools/app_errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create crea un nuevo token y lo almacena en la db
func Create(userID primitive.ObjectID, options ...interface{}) (*Token, error) {
	token, err := insert(userID, options...)
	if err != nil {
		return nil, err
	}

	cacheAdd(token)

	return token, nil
}

// Find busca un token en la db
func Find(tokenID string, options ...interface{}) (*Token, error) {
	token, err := findByID(tokenID, options...)
	if err != nil {
		return nil, err
	}

	cacheAdd(token)

	return token, nil
}

// Validate dado un tokenString devuelve el Token asociado
func Validate(tokenString string, options ...interface{}) (*Token, error) {
	if token, err := cacheGet(tokenString); err == nil {
		return token, err
	}

	// Sino validamos el token y lo agregamos al cache
	tokenID, _, err := ExtractPayload(tokenString)
	if err != nil {
		return nil, err
	}

	// Buscamos el token en la db para validarlo
	token, err := Find(tokenID, options...)
	if err != nil || !token.Enabled {
		return nil, app_errors.Unauthorized
	}

	return token, nil
}

// Invalidate invalida un token
func Invalidate(tokenString string, options ...interface{}) error {
	token, err := Validate(tokenString, options...)
	if err != nil {
		return app_errors.Unauthorized
	}

	if err = delete(token.ID, options...); err != nil {
		return err
	}

	cacheRemove(token)

	go func() {
		if err = rabbit.Get(options...).SendLogout("bearer " + tokenString); err != nil {
			glog.Info("Rabbit logout no se pudo enviar")
		}
	}()

	return nil
}
