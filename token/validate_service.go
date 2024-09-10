package token

import (
	"github.com/nmarsollier/authgo/tools/errs"
)

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
		return nil, errs.Unauthorized
	}

	return token, nil
}
