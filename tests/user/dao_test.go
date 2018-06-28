package user

import (
	"testing"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/nmarsollier/authgo/tests/mocks"
	"github.com/nmarsollier/authgo/tools/errors"
	"github.com/nmarsollier/authgo/user"
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

func TestFindByIdInvalid(t *testing.T) {
	testDao := user.MockedDao(nil)

	_, err := testDao.FindByID("__invalid__")
	assert.NotNil(t, err)
	assert.Equal(t, err, errors.ErrID)
}

func TestFindByIdOk(t *testing.T) {
	mConn := new(mocks.Collection)
	testDao := user.MockedDao(mConn)

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		mocks.Decoder(func(v interface{}) error {
			if user, ok := v.(*user.User); ok {
				user.ID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")
			}
			return nil
		}),
	)

	user, err := testDao.FindByID("5b2a6b7d893dc92de5a8b833")
	assert.Nil(t, err)
	assert.Equal(t, user.ID.Hex(), "5b2a6b7d893dc92de5a8b833")
}

func TestFindByIdNotFound(t *testing.T) {
	mConn := new(mocks.Collection)
	testDao := user.MockedDao(mConn)

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		mocks.Decoder(func(v interface{}) error {
			return mongo.ErrNoDocuments
		}),
	)

	_, err := testDao.FindByID("5b2a6b7d893dc92de5a8b833")
	assert.Equal(t, err, mongo.ErrNoDocuments)
}

func TestInsertOk(t *testing.T) {
	mConn := new(mocks.Collection)
	testDao := user.MockedDao(mConn)

	usr := user.NewUser()
	usr.Name = "Name"
	usr.Login = "Login"
	usr.SetPasswordText("Login")

	mConn.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).Return(usr.ID, nil)

	usr, err := testDao.Insert(usr)
	assert.Nil(t, err)
	assert.NotNil(t, usr.ID)
}

func TestInsertError(t *testing.T) {
	mConn := new(mocks.Collection)
	testDao := user.MockedDao(mConn)

	usr := user.NewUser()

	mConn.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).Return(usr.ID, nil)

	_, err := testDao.Insert(usr)
	validation, ok := err.(validator.ValidationErrors)
	assert.Equal(t, ok, true)
	assert.Equal(t, 3, len(validation))
	assert.Equal(t, "Name", validation[0].Field())
	assert.Equal(t, "Login", validation[1].Field())
	assert.Equal(t, "Password", validation[2].Field())
}

func TestUpdateOk(t *testing.T) {
	mConn := new(mocks.Collection)
	testDao := user.MockedDao(mConn)

	usr := user.NewUser()
	usr.Name = "Name"
	usr.Login = "Login"
	usr.SetPasswordText("Login")

	mConn.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(1, 1, 1, nil)

	usr, err := testDao.Update(usr)
	assert.Nil(t, err)
	assert.NotNil(t, usr.ID)
}

func TestUpdateError(t *testing.T) {
	mConn := new(mocks.Collection)
	testDao := user.MockedDao(mConn)

	usr := user.NewUser()

	mConn.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(1, 1, 1, nil)

	_, err := testDao.Update(usr)
	validation, ok := err.(validator.ValidationErrors)
	assert.Equal(t, ok, true)
	assert.Equal(t, 3, len(validation))
	assert.Equal(t, "Name", validation[0].Field())
	assert.Equal(t, "Login", validation[1].Field())
	assert.Equal(t, "Password", validation[2].Field())
}

func TestFindByLoginOk(t *testing.T) {
	mConn := new(mocks.Collection)
	testDao := user.MockedDao(mConn)

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		mocks.Decoder(func(v interface{}) error {
			if user, ok := v.(*user.User); ok {
				user.ID, _ = objectid.FromHex("5b2a6b7d893dc92de5a8b833")
			}
			return nil
		}),
	)

	user, err := testDao.FindByLogin("5b2a6b7d893dc92de5a8b833")
	assert.Nil(t, err)
	assert.Equal(t, user.ID.Hex(), "5b2a6b7d893dc92de5a8b833")
}
