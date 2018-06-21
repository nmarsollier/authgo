package token

import (
	"testing"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/nmarsollier/authgo/tools/errors"
	"github.com/nmarsollier/authgo/tools/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testMongoUnique = mongo.WriteErrors{
	mongo.WriteError{
		Index:   1,
		Code:    11000,
		Message: "Index",
	},
}
var testMongoError = mongo.WriteErrors{
	mongo.WriteError{
		Index:   1,
		Code:    11001,
		Message: "Other",
	},
}

func TestGetId(t *testing.T) {
	testDao := newTestingDao(nil)
	id, err := testDao.getID("5b2a6b7d893dc92de5a8b833")
	assert.Nil(t, err)
	assert.Equal(t, id.Hex(), "5b2a6b7d893dc92de5a8b833")

	id, err = testDao.getID("invalid")
	assert.NotNil(t, err)
	assert.Equal(t, err, errors.ErrID)
}

func TestFindByIdInvalid(t *testing.T) {
	testDao := newTestingDao(nil)
	_, err := testDao.findByID("__invalid__")
	assert.NotNil(t, err)
	assert.Equal(t, err, errors.Unauthorized)
}

func TestFindByIdOk(t *testing.T) {
	mConn := new(test.FakeCollection)
	testDao := newTestingDao(mConn)

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		test.FakeDecoder(func(v interface{}) error {
			if token, ok := v.(*Token); ok {
				token.ID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")
				token.UserID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")
			}
			return nil
		}),
	)

	token, err := testDao.findByID("5b2a6b7d893dc92de5a8b833")
	assert.Nil(t, err)
	assert.Equal(t, token.ID.Hex(), "5b2a6b7d893dc92de5a8b833")
}

func TestFindByIdNotFound(t *testing.T) {
	mConn := new(test.FakeCollection)
	testDao := newTestingDao(mConn)

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		test.FakeDecoder(func(v interface{}) error {
			return mongo.ErrNoDocuments
		}),
	)

	_, err := testDao.findByID("5b2a6b7d893dc92de5a8b833")
	assert.Equal(t, err, errors.Unauthorized)
}

func TestInsertOk(t *testing.T) {
	mConn := new(test.FakeCollection)
	testDao := newTestingDao(mConn)

	token := newToken()
	token.UserID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")

	mConn.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).Return(token.ID, nil)

	token, err := testDao.insert(token)
	assert.Nil(t, err)
	assert.NotNil(t, token.ID)
}

func TestUpdateOk(t *testing.T) {
	mConn := new(test.FakeCollection)
	testDao := newTestingDao(mConn)

	token := newToken()
	token.UserID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")

	mConn.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(1, 1, 1, nil)

	token, err := testDao.update(token)
	assert.Nil(t, err)
	assert.NotNil(t, token.ID)
}

func TestFindByUserIdInvalid(t *testing.T) {
	testDao := newTestingDao(nil)
	_, err := testDao.findByUserID("__invalid__")
	assert.NotNil(t, err)
	assert.Equal(t, err, errors.Unauthorized)
}

func TestFindByUserIdOk(t *testing.T) {
	mConn := new(test.FakeCollection)
	testDao := newTestingDao(mConn)

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		test.FakeDecoder(func(v interface{}) error {
			if token, ok := v.(*Token); ok {
				token.ID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")
				token.UserID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")
			}
			return nil
		}),
	)

	token, err := testDao.findByUserID("5b2a6b7d893dc92de5a8b833")
	assert.Nil(t, err)
	assert.Equal(t, token.ID.Hex(), "5b2a6b7d893dc92de5a8b833")
}

func TestDeleteInvalid(t *testing.T) {
	testDao := newTestingDao(nil)
	err := testDao.delete("__invalid__")
	assert.NotNil(t, err)
	assert.Equal(t, err, errors.Unauthorized)
}

func TestDeleteOk(t *testing.T) {
	mConn := new(test.FakeCollection)
	testDao := newTestingDao(mConn)

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

	err := testDao.delete("5b2a6b7d893dc92de5a8b833")
	assert.Nil(t, err)
}
