package test

import (
	"testing"

	"github.com/mongodb/mongo-go-driver/bson/objectid"

	"github.com/nmarsollier/authgo/mocks"
	"github.com/nmarsollier/authgo/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	mConn := new(mocks.Collection)
	service := security.MockedService(mConn)

	ID, _ := objectid.FromHex("5b2a6b7d893dc92de5a8b833")

	mConn.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).Return(ID, nil)

	token, err := service.Create(ID)

	assert.Nil(t, err)
	assert.NotNil(t, token)
}

func TestFind(t *testing.T) {
	mConn := new(mocks.Collection)
	service := security.MockedService(mConn)

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		mocks.Decoder(func(v interface{}) error {
			return nil
		}),
	)

	token, err := service.Find("5b2a6b7d893dc92de5a8b833")
	assert.Nil(t, err)
	assert.NotNil(t, token)
}

func TestValidate(t *testing.T) {
	mConn := new(mocks.Collection)
	service := security.MockedService(mConn)

	ID, _ := objectid.FromHex("5b2a6b7d893dc92de5a8b833")
	mConn.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).Return(ID, nil)
	token, err := service.Create(ID)
	tokenString, _ := token.Encode()

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		mocks.Decoder(func(v interface{}) error {
			return nil
		}),
	)

	tkn, err := service.Validate(tokenString)
	assert.Nil(t, err)
	assert.Equal(t, tkn.ID.Hex(), token.ID.Hex())
}

func TestValidateCache(t *testing.T) {
	mConn := new(mocks.Collection)
	service := security.MockedService(mConn)

	ID, _ := objectid.FromHex("5b2a6b7d893dc92de5a8b833")
	userID, _ := objectid.FromHex("11116b7d893dc92de5a8b833")
	mConn.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).Return(ID, nil)
	token, err := service.Create(ID)
	token.UserID = userID
	tokenString, _ := token.Encode()

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		mocks.Decoder(func(v interface{}) error {
			if token, ok := v.(*security.Token); ok {
				token.ID, _ = objectid.FromHex("112a6b7d893dc92de5a8b811")
				token.Enabled = true
			}
			return nil
		}),
	)

	token, err = service.Validate(tokenString)
	assert.Nil(t, err)
	assert.Equal(t, token.ID.Hex(), "112a6b7d893dc92de5a8b811")
}
