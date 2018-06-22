package test

import (
	"testing"

	"github.com/mongodb/mongo-go-driver/bson/objectid"

	"github.com/nmarsollier/authgo/mocks"
	"github.com/nmarsollier/authgo/security"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	service := security.MockedService(mocks.MockedTokenDao())

	ID, _ := objectid.FromHex("5b2a6b7d893dc92de5a8b833")
	token, err := service.Create(ID)

	assert.Nil(t, err)
	assert.NotNil(t, token)
}

func TestFind(t *testing.T) {
	service := security.MockedService(mocks.MockedTokenDao())

	token, err := service.Find("5b2a6b7d893dc92de5a8b833")
	assert.Nil(t, err)
	assert.NotNil(t, token)
}

func TestValidate(t *testing.T) {
	service := security.MockedService(mocks.MockedTokenDao())

	ID, _ := objectid.FromHex("5b2a6b7d893dc92de5a8b833")
	token, err := service.Create(ID)
	tokenString, _ := token.Encode()

	tkn, err := service.Validate(tokenString)
	assert.Nil(t, err)
	assert.Equal(t, tkn.ID.Hex(), token.ID.Hex())
}

func TestValidateCache(t *testing.T) {
	service := security.MockedService(mocks.MockedTokenDao())

	ID, _ := objectid.FromHex("111111111111111111111111")
	userID, _ := objectid.FromHex("11116b7d893dc92de5a8b833")
	token, err := service.Create(ID)
	token.UserID = userID
	tokenString, _ := token.Encode()

	token, err = service.Validate(tokenString)
	assert.Nil(t, err)
	assert.Equal(t, token.ID.Hex(), "111111111111111111111111")
}
