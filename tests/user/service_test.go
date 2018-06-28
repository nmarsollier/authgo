package user

import (
	"testing"

	"github.com/nmarsollier/authgo/tests/mocks"
	"github.com/nmarsollier/authgo/user"

	"github.com/stretchr/testify/assert"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/nmarsollier/authgo/tools/errors"
)

func TestSignUpOk(t *testing.T) {
	srv := user.MockedService(mocks.MockedUserDaoCustom("", "", "", true), mocks.MockedSecurityService())

	req := user.SignUpRequest{
		Name:     "Test",
		Login:    "Login",
		Password: "Pass",
	}

	id, err := srv.SignUp(&req)

	assert.Nil(t, err)
	assert.NotNil(t, id)
	assert.NotEmpty(t, id)
}

func TestSignUpError(t *testing.T) {
	req := user.SignUpRequest{}
	validate := validator.New()
	validate.SetTagName("binding")
	errResult := validate.Struct(req)

	srv := user.MockedService(mocks.MockedUserDaoError(errResult), mocks.MockedSecurityService())

	_, err := srv.SignUp(&req)
	validation, ok := err.(validator.ValidationErrors)
	assert.Equal(t, ok, true)
	assert.Equal(t, 3, len(validation))
	assert.Equal(t, "Name", validation[0].Field())
	assert.Equal(t, "Password", validation[1].Field())
	assert.Equal(t, "Login", validation[2].Field())
}

func TestSignIn(t *testing.T) {
	srv := user.MockedService(mocks.MockedUserDao(), mocks.MockedSecurityService())

	id, err := srv.SignIn("User", "Password")

	assert.Nil(t, err)
	assert.NotNil(t, id)
	assert.NotEmpty(t, id)
}

func TestSignInError(t *testing.T) {
	srv := user.MockedService(mocks.MockedUserDao(), mocks.MockedSecurityService())

	_, err := srv.SignIn("User", "Password1")

	assert.Equal(t, user.ErrPassword, err)
}

func TestSignInError1(t *testing.T) {
	srv := user.MockedService(mocks.MockedUserDaoCustom("Name", "Login", "Password", false), mocks.MockedSecurityService())

	_, err := srv.SignIn("User", "Password")

	assert.Equal(t, errors.Unauthorized, err)
}

func TestChangePassword(t *testing.T) {
	srv := user.MockedService(mocks.MockedUserDao(), mocks.MockedSecurityService())

	err := srv.ChangePassword("5b2a6b7d893dc92de5a8b833", "Password", "Password1")
	assert.Nil(t, err)

	srv = user.MockedService(mocks.MockedUserDao(), mocks.MockedSecurityServiceError(user.ErrPassword))
	err = srv.ChangePassword("5b2a6b7d893dc92de5a8b833", "Password1", "Password1")
	assert.Equal(t, user.ErrPassword, err)

}
