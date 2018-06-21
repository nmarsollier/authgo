package token

import (
	"testing"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/nmarsollier/authgo/tests/mocks"
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/errors"
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
	testDao := token.NewTestingDao(nil)
	id, err := testDao.GetID("5b2a6b7d893dc92de5a8b833")
	assert.Nil(t, err)
	assert.Equal(t, id.Hex(), "5b2a6b7d893dc92de5a8b833")

	id, err = testDao.GetID("invalid")
	assert.NotNil(t, err)
	assert.Equal(t, err, errors.ErrID)
}

func TestFindByIdInvalid(t *testing.T) {
	testDao := token.NewTestingDao(nil)
	_, err := testDao.FindByID("__invalid__")
	assert.NotNil(t, err)
	assert.Equal(t, err, errors.Unauthorized)
}

func TestFindByIdOk(t *testing.T) {
	mConn := new(mocks.FakeCollection)
	testDao := token.NewTestingDao(mConn)

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		mocks.FakeDecoder(func(v interface{}) error {
			if token, ok := v.(*token.Token); ok {
				token.ID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")
				token.UserID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")
			}
			return nil
		}),
	)

	token, err := testDao.FindByID("5b2a6b7d893dc92de5a8b833")
	assert.Nil(t, err)
	assert.Equal(t, token.ID.Hex(), "5b2a6b7d893dc92de5a8b833")
}

func TestFindByIdNotFound(t *testing.T) {
	mConn := new(mocks.FakeCollection)
	testDao := token.NewTestingDao(mConn)

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		mocks.FakeDecoder(func(v interface{}) error {
			return mongo.ErrNoDocuments
		}),
	)

	_, err := testDao.FindByID("5b2a6b7d893dc92de5a8b833")
	assert.Equal(t, err, errors.Unauthorized)
}

func TestInsertOk(t *testing.T) {
	mConn := new(mocks.FakeCollection)
	testDao := token.NewTestingDao(mConn)

	token := token.NewToken()
	token.UserID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")

	mConn.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).Return(token.ID, nil)

	token, err := testDao.Insert(token)
	assert.Nil(t, err)
	assert.NotNil(t, token.ID)
}

func TestUpdateOk(t *testing.T) {
	mConn := new(mocks.FakeCollection)
	testDao := token.NewTestingDao(mConn)

	token := token.NewToken()
	token.UserID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")

	mConn.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(1, 1, 1, nil)

	token, err := testDao.Update(token)
	assert.Nil(t, err)
	assert.NotNil(t, token.ID)
}

func TestFindByUserIdInvalid(t *testing.T) {
	testDao := token.NewTestingDao(nil)
	_, err := testDao.FindByUserID("__invalid__")
	assert.NotNil(t, err)
	assert.Equal(t, err, errors.Unauthorized)
}

func TestFindByUserIdOk(t *testing.T) {
	mConn := new(mocks.FakeCollection)
	testDao := token.NewTestingDao(mConn)

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		mocks.FakeDecoder(func(v interface{}) error {
			if token, ok := v.(*token.Token); ok {
				token.ID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")
				token.UserID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")
			}
			return nil
		}),
	)

	token, err := testDao.FindByUserID("5b2a6b7d893dc92de5a8b833")
	assert.Nil(t, err)
	assert.Equal(t, token.ID.Hex(), "5b2a6b7d893dc92de5a8b833")
}

func TestDeleteInvalid(t *testing.T) {
	testDao := token.NewTestingDao(nil)
	err := testDao.Delete("__invalid__")
	assert.NotNil(t, err)
	assert.Equal(t, err, errors.Unauthorized)
}

func TestDeleteOk(t *testing.T) {
	mConn := new(mocks.FakeCollection)
	testDao := token.NewTestingDao(mConn)

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		mocks.FakeDecoder(func(v interface{}) error {
			if token, ok := v.(*token.Token); ok {
				token.ID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")
				token.UserID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")
			}
			return nil
		}),
	)
	mConn.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(1, 1, 1, nil)

	err := testDao.Delete("5b2a6b7d893dc92de5a8b833")
	assert.Nil(t, err)
}
