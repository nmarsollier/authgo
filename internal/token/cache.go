package token

import (
	"time"

	"github.com/nmarsollier/authgo/internal/engine/errs"
	gocache "github.com/patrickmn/go-cache"
)

type TokenCache interface {
	Add(token *Token) error
	Get(tokenString string) (*Token, error)
	Remove(token string)
}

func NewTokenCache() TokenCache {
	cache := gocache.New(60*time.Minute, 10*time.Minute)

	return &tokenCache{
		cache: cache,
	}
}

type tokenCache struct {
	cache *gocache.Cache
}

// Add genera un nuevo token al cache
func (c *tokenCache) Add(token *Token) error {
	tokenString, err := Encode(token)
	if err != nil {
		return err
	}
	c.cache.Set(tokenString, token, gocache.DefaultExpiration)
	return nil
}

func (c *tokenCache) Get(tokenString string) (*Token, error) {
	// Si esta en cache, retornamos el cache
	if found, ok := c.cache.Get(tokenString); ok {
		if token, ok := found.(*Token); ok {
			return token, nil
		}
	}

	return nil, errs.NotFound
}

// Remove elimia un token del cache
func (c *tokenCache) Remove(tokenString string) {
	c.cache.Delete(tokenString)
}
