package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCache(t *testing.T) {
	userID, _ := primitive.ObjectIDFromHex("111111111111111111111111")
	token := newToken(userID)
	token1 := newToken(userID)
	token2 := newToken(userID)

	assert.Nil(t, cacheAdd(token))
	assert.Nil(t, cacheAdd(token1))
	assert.Nil(t, cacheAdd(token2))

	tokenString, _ := token.Encode()
	found, err := cacheGet(tokenString)
	assert.Nil(t, err)
	assert.Equal(t, token.ID.Hex(), found.ID.Hex())

	tokenString, _ = token1.Encode()
	found, err = cacheGet(tokenString)
	assert.Nil(t, err)
	assert.Equal(t, token1.ID.Hex(), found.ID.Hex())

	tokenString, _ = token2.Encode()
	found, err = cacheGet(tokenString)
	assert.Nil(t, err)
	assert.Equal(t, token2.ID.Hex(), found.ID.Hex())

	cacheRemove(token)
	tokenString, err = token.Encode()
	found, err = cacheGet(tokenString)
	assert.NotNil(t, err)
	assert.Nil(t, found)
}
