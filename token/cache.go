package token

import (
	"time"

	"github.com/nmarsollier/authgo/tools/errors"
	gocache "github.com/patrickmn/go-cache"
)

var cache = gocache.New(60*time.Minute, 10*time.Minute)

// Add genera un nuevo token al cache
func cacheAdd(token *Token) error {
	tokenString, err := Encode(token)
	if err != nil {
		return err
	}
	cache.Set(tokenString, token, gocache.DefaultExpiration)
	return nil
}

func cacheGet(tokenString string) (*Token, error) {
	// Si esta en cache, retornamos el cache
	if found, ok := cache.Get(tokenString); ok {
		if token, ok := found.(*Token); ok {
			return token, nil
		}
	}

	return nil, errors.NotFound
}

// Remove elimia un token del cache
func cacheRemove(token *Token) {
	if tokenString, err := Encode(token); err == nil {
		cache.Delete(tokenString)
	}
}
