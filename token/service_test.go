package token

import (
	"testing"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/test"
)

func TestCreate(t *testing.T) {
	collectionTest = fakeServiceCollection()

	tokenID, _ := objectid.FromHex("5b2a6b7d893dc92de5a8b833")
	token, err := Create(tokenID)
	assert.NotNil(t, token)
	assert.Nil(t, err)

	payload, err := extractPayload(token)
	assert.Nil(t, err)
	assert.NotNil(t, payload.TokenID)
	assert.Equal(t, payload.UserID, "5b2a6b7d893dc92de5a8b833")
}

func TestValidate(t *testing.T) {
	collectionTest = fakeServiceCollection()

	token := "__invalid__"

	payload, err := Validate(token)
	assert.NotNil(t, err)
	assert.Nil(t, payload)

	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbklEIjoiNWIyYWZkMjA2MDFlZDljNzQ0NDVhYjU3IiwidXNlcklEIjoiNWIyYTZiN2Q4OTNkYzkyZGU1YThiODMzIn0.RBcB_B5D6uL3JXRbi2xe-V9LytIOxxLSnXv0_-rFAVU"

	payload, err = Validate(token)
	assert.Nil(t, err)
	assert.NotNil(t, payload.TokenID)
	assert.Equal(t, payload.UserID, "5b2a6b7d893dc92de5a8b833")
}

func TestInvalidate(t *testing.T) {
	collectionTest = fakeServiceCollection()

	token := "__invalid__"

	err := Invalidate(token)
	assert.NotNil(t, err)

	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbklEIjoiNWIyYWZkMjA2MDFlZDljNzQ0NDVhYjU3IiwidXNlcklEIjoiNWIyYTZiN2Q4OTNkYzkyZGU1YThiODMzIn0.RBcB_B5D6uL3JXRbi2xe-V9LytIOxxLSnXv0_-rFAVU"

	err = Invalidate(token)
	assert.Nil(t, err)
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
