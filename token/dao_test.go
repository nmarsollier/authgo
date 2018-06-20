package token

import (
	"testing"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/nmarsollier/authgo/test"
	"github.com/nmarsollier/authgo/tools/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var MongoUnique = mongo.WriteErrors{
	mongo.WriteError{
		Index:   1,
		Code:    11000,
		Message: "Index",
	},
}
var MongoError = mongo.WriteErrors{
	mongo.WriteError{
		Index:   1,
		Code:    11001,
		Message: "Other",
	},
}

func TestGetId(t *testing.T) {
	id, err := getID("5b2a6b7d893dc92de5a8b833")
	assert.Nil(t, err)
	assert.Equal(t, id.Hex(), "5b2a6b7d893dc92de5a8b833")

	id, err = getID("invalid")
	assert.NotNil(t, err)
	assert.Equal(t, err, errors.ErrID)
}

func TestFindByIdInvalid(t *testing.T) {
	_, err := findByID("__invalid__")
	assert.NotNil(t, err)
	assert.Equal(t, err, errors.Unauthorized)
}

func TestFindByIdOk(t *testing.T) {
	mConn := new(test.FakeCollection)
	collectionTest = mConn

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		test.FakeDecoder(func(v interface{}) error {
			if token, ok := v.(*Token); ok {
				token.ID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")
				token.UserID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")
			}
			return nil
		}),
	)

	token, err := findByID("5b2a6b7d893dc92de5a8b833")
	assert.Nil(t, err)
	assert.Equal(t, token.ID.Hex(), "5b2a6b7d893dc92de5a8b833")
}

func TestFindByIdNotFound(t *testing.T) {
	mConn := new(test.FakeCollection)
	collectionTest = mConn

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		test.FakeDecoder(func(v interface{}) error {
			return mongo.ErrNoDocuments
		}),
	)

	_, err := findByID("5b2a6b7d893dc92de5a8b833")
	assert.Equal(t, err, errors.Unauthorized)
}

func TestInsertOk(t *testing.T) {
	mConn := new(test.FakeCollection)
	collectionTest = mConn

	token := newToken()
	token.UserID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")

	mConn.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).Return(token.ID, nil)

	token, err := insert(token)
	assert.Nil(t, err)
	assert.NotNil(t, token.ID)
}

func TestInsertUniqueError(t *testing.T) {
	mConn := new(test.FakeCollection)
	collectionTest = mConn

	token := newToken()

	mConn.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, MongoUnique)

	token, err := insert(token)
	assert.Equal(t, err, MongoUnique)
}

func TestInsertOtherError(t *testing.T) {
	mConn := new(test.FakeCollection)
	collectionTest = mConn

	token := newToken()

	mConn.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, MongoError)

	token, err := insert(token)
	assert.Equal(t, err, MongoError)
}

func TestUpdateOk(t *testing.T) {
	mConn := new(test.FakeCollection)
	collectionTest = mConn

	token := newToken()
	token.UserID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")

	mConn.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(1, 1, 1, nil)

	token, err := update(token)
	assert.Nil(t, err)
	assert.NotNil(t, token.ID)
}

func TestUpdateUniqueError(t *testing.T) {
	mConn := new(test.FakeCollection)
	collectionTest = mConn

	token := newToken()

	mConn.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(1, 1, 1, MongoUnique)

	token, err := update(token)
	assert.Equal(t, err, MongoUnique)
}

func TestUpdateOtherError(t *testing.T) {
	mConn := new(test.FakeCollection)
	collectionTest = mConn

	token := newToken()

	mConn.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(1, 1, 1, MongoError)

	token, err := update(token)
	assert.Equal(t, err, MongoError)
}

func TestFindByUserIdInvalid(t *testing.T) {
	_, err := findByUserID("__invalid__")
	assert.NotNil(t, err)
	assert.Equal(t, err, errors.Unauthorized)
}

func TestFindByUserIdOk(t *testing.T) {
	mConn := new(test.FakeCollection)
	collectionTest = mConn

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		test.FakeDecoder(func(v interface{}) error {
			if token, ok := v.(*Token); ok {
				token.ID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")
				token.UserID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")
			}
			return nil
		}),
	)

	token, err := findByUserID("5b2a6b7d893dc92de5a8b833")
	assert.Nil(t, err)
	assert.Equal(t, token.ID.Hex(), "5b2a6b7d893dc92de5a8b833")
}

func TestFindByUserIdNotFound(t *testing.T) {
	mConn := new(test.FakeCollection)
	collectionTest = mConn

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		test.FakeDecoder(func(v interface{}) error {
			return mongo.ErrNoDocuments
		}),
	)

	_, err := findByUserID("5b2a6b7d893dc92de5a8b833")
	assert.Equal(t, err, errors.Unauthorized)
}

func TestDeleteInvalid(t *testing.T) {
	err := delete("__invalid__")
	assert.NotNil(t, err)
	assert.Equal(t, err, errors.Unauthorized)
}

func TestDeleteOk(t *testing.T) {
	mConn := new(test.FakeCollection)
	collectionTest = mConn

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		test.FakeDecoder(func(v interface{}) error {
			if token, ok := v.(*Token); ok {
				token.ID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")
				token.UserID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")
			}
			return nil
		}),
	)
	mConn.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(1, 1, 1, nil)

	err := delete("5b2a6b7d893dc92de5a8b833")
	assert.Nil(t, err)
}
