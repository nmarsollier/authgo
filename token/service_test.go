package token

import (
	"testing"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nmarsollier/authgo/test"
	"github.com/nmarsollier/authgo/tools/db"
)

func TestCreateValidate(t *testing.T) {
	collectionTest = fakeServiceCollection()

	tokenID, _ := objectid.FromHex("5b2a6b7d893dc92de5a8b833")
	token, err := Create(tokenID)
	assert.NotNil(t, token)
	assert.Nil(t, err)

	payload, err := Validate(token)
	assert.Nil(t, err)
	assert.NotNil(t, payload.TokenID)
	assert.Equal(t, payload.UserID, "5b2a6b7d893dc92de5a8b833")
}

func fakeServiceCollection() db.Collection {
	mConn := new(test.FakeCollection)

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		test.FakeDecoder(func(v interface{}) error {
			if token, ok := v.(*Token); ok {
				token.ID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")
				token.UserID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")
				token.Enabled = true
			}
			return nil
		}),
	)

	tokenID, _ := objectid.FromHex("5b2a6b7d893dc92de5a8b833")

	mConn.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).Return(tokenID, nil)

	mConn.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(1, 1, 1, nil)

	return mConn
}
