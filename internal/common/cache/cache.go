package cache

import (
	"time"

	"github.com/nmarsollier/authgo/internal/common/errs"
	gocache "github.com/patrickmn/go-cache"
)

type Cache[T any] interface {
	Add(key string, value *T) error
	Get(key string) (*T, error)
	Remove(key string)
}

// NewCache creates a new instance of Cache with a specified type.
// It initializes the cache with a default expiration time of 60 minutes
// and a cleanup interval of 10 minutes.
//
// T: The type of the items to be stored in the cache.
//
// Returns:
//
//	A new instance of Cache with the specified type.
func NewCache[T any]() Cache[T] {
	return &theCache[T]{
		cache: gocache.New(60*time.Minute, 10*time.Minute),
	}
}

type theCache[T any] struct {
	cache *gocache.Cache
}

func (c *theCache[T]) Add(key string, value *T) error {
	c.cache.Set(key, value, gocache.DefaultExpiration)
	return nil
}

func (c *theCache[T]) Get(key string) (*T, error) {
	if found, ok := c.cache.Get(key); ok {
		if token, ok := found.(*T); ok {
			return token, nil
		}
	}

	return nil, errs.NotFound
}

func (c *theCache[T]) Remove(key string) {
	c.cache.Delete(key)
}
