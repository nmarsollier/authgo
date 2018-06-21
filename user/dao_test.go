package user

import (
	"testing"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	validator "gopkg.in/go-playground/validator.v9"

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
	assert.Equal(t, err, errors.ErrID)
}

func TestFindByIdOk(t *testing.T) {
	mConn := new(test.FakeCollection)
	testDao := newTestingDao(mConn)

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		test.FakeDecoder(func(v interface{}) error {
			if user, ok := v.(*User); ok {
				user.ID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")
			}
			return nil
		}),
	)

	user, err := testDao.findByID("5b2a6b7d893dc92de5a8b833")
	assert.Nil(t, err)
	assert.Equal(t, user.ID.Hex(), "5b2a6b7d893dc92de5a8b833")
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
	assert.Equal(t, err, mongo.ErrNoDocuments)
}

func TestInsertOk(t *testing.T) {
	mConn := new(test.FakeCollection)
	testDao := newTestingDao(mConn)

	user := newUser()
	user.Name = "Name"
	user.Login = "Login"
	user.setPasswordText("Login")

	mConn.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).Return(user.ID, nil)

	user, err := testDao.insert(user)
	assert.Nil(t, err)
	assert.NotNil(t, user.ID)
}

func TestInsertError(t *testing.T) {
	mConn := new(test.FakeCollection)
	testDao := newTestingDao(mConn)

	user := newUser()

	mConn.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).Return(user.ID, nil)

	_, err := testDao.insert(user)
	validation, ok := err.(validator.ValidationErrors)
	assert.Equal(t, ok, true)
	assert.Equal(t, 3, len(validation))
	assert.Equal(t, "Name", validation[0].Field())
	assert.Equal(t, "Login", validation[1].Field())
	assert.Equal(t, "Password", validation[2].Field())
}

func TestUpdateOk(t *testing.T) {
	mConn := new(test.FakeCollection)
	testDao := newTestingDao(mConn)

	user := newUser()
	user.Name = "Name"
	user.Login = "Login"
	user.setPasswordText("Login")

	mConn.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(1, 1, 1, nil)

	user, err := testDao.update(user)
	assert.Nil(t, err)
	assert.NotNil(t, user.ID)
}

func TestUpdateError(t *testing.T) {
	mConn := new(test.FakeCollection)
	testDao := newTestingDao(mConn)

	user := newUser()

	mConn.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(1, 1, 1, nil)

	_, err := testDao.update(user)
	validation, ok := err.(validator.ValidationErrors)
	assert.Equal(t, ok, true)
	assert.Equal(t, 3, len(validation))
	assert.Equal(t, "Name", validation[0].Field())
	assert.Equal(t, "Login", validation[1].Field())
	assert.Equal(t, "Password", validation[2].Field())
}

func TestFindByLoginOk(t *testing.T) {
	mConn := new(test.FakeCollection)
	testDao := newTestingDao(mConn)

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		test.FakeDecoder(func(v interface{}) error {
			if user, ok := v.(*User); ok {
				user.ID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")
			}
			return nil
		}),
	)

	user, err := testDao.findByLogin("5b2a6b7d893dc92de5a8b833")
	assert.Nil(t, err)
	assert.Equal(t, user.ID.Hex(), "5b2a6b7d893dc92de5a8b833")
}

func TestDeleteOk(t *testing.T) {
	mConn := new(test.FakeCollection)
	testDao := newTestingDao(mConn)

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		test.FakeDecoder(func(v interface{}) error {
			if user, ok := v.(*User); ok {
				user.ID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")
			}
			return nil
		}),
	)
	mConn.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(1, 1, 1, nil)

	err := testDao.delete("5b2a6b7d893dc92de5a8b833")
	assert.Nil(t, err)
}
