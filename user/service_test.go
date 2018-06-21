package user

import (
	"testing"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/errors"
	"github.com/nmarsollier/authgo/tools/test"
)

func TestSignUpOk(t *testing.T) {
	FakeServiceCollection(true)

	req := SignUpRequest{
		Name:     "Test",
		Login:    "Login",
		Password: "Pass",
	}

	id, err := SignUp(&req)

	assert.Nil(t, err)
	assert.NotNil(t, id)
	assert.NotEmpty(t, id)
}

func TestSignUpError(t *testing.T) {
	FakeServiceCollection(true)

	req := SignUpRequest{}

	_, err := SignUp(&req)

	validation, ok := err.(validator.ValidationErrors)
	assert.Equal(t, ok, true)
	assert.Equal(t, 2, len(validation))
	assert.Equal(t, "Name", validation[0].Field())
	assert.Equal(t, "Login", validation[1].Field())
}

func TestSignIn(t *testing.T) {
	FakeServiceCollection(true)

	id, err := SignIn("User", "Password")

	assert.Nil(t, err)
	assert.NotNil(t, id)
	assert.NotEmpty(t, id)
}

func TestSignInError(t *testing.T) {
	FakeServiceCollection(true)

	_, err := SignIn("User", "Password1")

	assert.Equal(t, err, ErrPassword)
}

func TestSignInError1(t *testing.T) {
	FakeServiceCollection(false)

	_, err := SignIn("User", "Password")

	assert.Equal(t, err, errors.Unauthorized)
}

func TestChangePassword(t *testing.T) {
	FakeServiceCollection(false)

	err := ChangePassword("5b2a6b7d893dc92de5a8b833", "Password", "Password1")
	assert.Nil(t, err)

	err = ChangePassword("5b2a6b7d893dc92de5a8b833", "Password1", "Password1")
	assert.Equal(t, err, ErrPassword)

}

func FakeServiceCollection(enabledUser bool) {
	FakeTokenServiceCollection()
	mConn := new(test.FakeCollection)
	CollectionTest = mConn

	objectID, _ := objectid.FromHex("5b2a6b7d893dc92de5a8b833")

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		test.FakeDecoder(func(v interface{}) error {
			if user, ok := v.(*User); ok {
				user.ID = objectID
				user.Name = "User"
				user.Login = "Login"
				user.Enabled = enabledUser
				user.setPasswordText("Password")
			}
			return nil
		}),
	)
	mConn.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(1, 1, 1, nil)

	mConn.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).Return(objectID, nil)
}

func FakeTokenServiceCollection() {
	mConn := new(test.FakeCollection)
	token.CollectionTest = mConn

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		test.FakeDecoder(func(v interface{}) error {
			if token, ok := v.(*token.Token); ok {
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
}
