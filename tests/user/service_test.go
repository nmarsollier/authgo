package user

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/nmarsollier/authgo/tests/mocks"
	"github.com/nmarsollier/authgo/user"

	"github.com/stretchr/testify/assert"

	"github.com/nmarsollier/authgo/tools/errors"
)

func TestSignUpOk(t *testing.T) {
	user.DaoInstance = mocks.MockedUserDaoCustom("", "", "", true)
	user.SecServiceInstance = mocks.MockedSecurityService()
	srv := user.NewService()

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

	user.DaoInstance = mocks.MockedUserDaoError(errResult)
	user.SecServiceInstance = mocks.MockedSecurityService()
	srv := user.NewService()

	_, err := srv.SignUp(&req)
	validation, ok := err.(validator.ValidationErrors)
	assert.Equal(t, ok, true)
	assert.Equal(t, 3, len(validation))
	assert.Equal(t, "Name", validation[0].Field())
	assert.Equal(t, "Password", validation[1].Field())
	assert.Equal(t, "Login", validation[2].Field())
}

func TestSignIn(t *testing.T) {
	user.DaoInstance = mocks.MockedUserDao()
	user.SecServiceInstance = mocks.MockedSecurityService()
	srv := user.NewService()

	id, err := srv.SignIn("User", "Password")

	assert.Nil(t, err)
	assert.NotNil(t, id)
	assert.NotEmpty(t, id)
}

func TestSignInError(t *testing.T) {
	user.DaoInstance = mocks.MockedUserDao()
	user.SecServiceInstance = mocks.MockedSecurityService()
	srv := user.NewService()

	_, err := srv.SignIn("User", "Password1")

	assert.Equal(t, user.ErrPassword, err)
}

func TestSignInError1(t *testing.T) {
	user.DaoInstance = mocks.MockedUserDaoCustom("Name", "Login", "Password", false)
	user.SecServiceInstance = mocks.MockedSecurityService()
	srv := user.NewService()

	_, err := srv.SignIn("User", "Password")

	assert.Equal(t, errors.Unauthorized, err)
}

func TestChangePassword(t *testing.T) {
	user.DaoInstance = mocks.MockedUserDao()
	user.SecServiceInstance = mocks.MockedSecurityService()
	srv := user.NewService()

	err := srv.ChangePassword("5b2a6b7d893dc92de5a8b833", "Password", "Password1")
	assert.Nil(t, err)

	user.DaoInstance = mocks.MockedUserDaoError(user.ErrPassword)

	err = srv.ChangePassword("5b2a6b7d893dc92de5a8b833", "Password1", "Password1")
	assert.Equal(t, user.ErrPassword, err)
}
