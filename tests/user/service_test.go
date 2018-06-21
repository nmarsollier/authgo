package user

import (
	"testing"

	"github.com/nmarsollier/authgo/user"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/errors"
)

func TestSignUpOk(t *testing.T) {
	srv := user.NewTestingService(newCustomFakeDao("", "", "", true), newFakeTokenService())

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

	srv := user.NewTestingService(newFakeErrorDao(errResult), newFakeTokenService())

	_, err := srv.SignUp(&req)
	validation, ok := err.(validator.ValidationErrors)
	assert.Equal(t, ok, true)
	assert.Equal(t, 3, len(validation))
	assert.Equal(t, "Name", validation[0].Field())
	assert.Equal(t, "Password", validation[1].Field())
	assert.Equal(t, "Login", validation[2].Field())
}

func TestSignIn(t *testing.T) {
	srv := user.NewTestingService(newFakeDao(), newFakeTokenService())

	id, err := srv.SignIn("User", "Password")

	assert.Nil(t, err)
	assert.NotNil(t, id)
	assert.NotEmpty(t, id)
}

func TestSignInError(t *testing.T) {
	srv := user.NewTestingService(newFakeDao(), newFakeTokenService())

	_, err := srv.SignIn("User", "Password1")

	assert.Equal(t, user.ErrPassword, err)
}

func TestSignInError1(t *testing.T) {
	srv := user.NewTestingService(newCustomFakeDao("Name", "Login", "Password", false), newFakeTokenService())

	_, err := srv.SignIn("User", "Password")

	assert.Equal(t, errors.Unauthorized, err)
}

func TestChangePassword(t *testing.T) {
	srv := user.NewTestingService(newFakeDao(), newFakeTokenService())

	err := srv.ChangePassword("5b2a6b7d893dc92de5a8b833", "Password", "Password1")
	assert.Nil(t, err)

	srv = user.NewTestingService(newFakeDao(), newFakeTokenErrorService(user.ErrPassword))
	err = srv.ChangePassword("5b2a6b7d893dc92de5a8b833", "Password1", "Password1")
	assert.Equal(t, user.ErrPassword, err)

}

func newFakeDao() user.Dao {
	return newCustomFakeDao("TestName", "Login", "Password", true)
}

func newFakeErrorDao(err error) user.Dao {
	result := fakeDao{}
	result.On("Collection").Return(nil, nil)
	result.On("Insert", mock.Anything).Return(nil, err)
	result.On("Update", mock.Anything).Return(nil, err)
	result.On("FindAll").Return(nil, err)
	result.On("FindByID", mock.Anything).Return(nil, err)
	result.On("FindByLogin", mock.Anything).Return(nil, err)
	result.On("Delete", mock.Anything).Return(err)
	result.On("GetID", mock.Anything).Return(nil, err)

	return &result
}

func newCustomFakeDao(name string, login string, passw string, userEnabled bool) user.Dao {
	result := fakeDao{}

	usr := user.NewUser()
	usr.Name = name
	usr.Login = login
	usr.Enabled = userEnabled
	usr.SetPasswordText(passw)
	users := make([]*user.User, 1)
	users[0] = usr

	result.On("Collection").Return(nil, nil)
	result.On("Insert", mock.Anything).Return(usr, nil)
	result.On("Update", mock.Anything).Return(usr, nil)
	result.On("FindAll").Return(usr, nil)
	result.On("FindByID", mock.Anything).Return(usr, nil)
	result.On("FindByLogin", mock.Anything).Return(usr, nil)
	result.On("Delete", mock.Anything).Return(nil)
	result.On("GetID", mock.Anything).Return(usr.ID, nil)

	return &result
}

type fakeDao struct {
	mock.Mock
}

func (mc *fakeDao) Collection() (db.Collection, error) {
	res := mc.Called()
	t, _ := res.Get(0).(db.Collection)
	err, _ := res.Get(1).(error)
	return t, err
}
func (mc *fakeDao) Insert(usr *user.User) (*user.User, error) {
	res := mc.Called(usr)
	t, _ := res.Get(0).(*user.User)
	err, _ := res.Get(1).(error)
	return t, err
}
func (mc *fakeDao) Update(usr *user.User) (*user.User, error) {
	res := mc.Called(usr)
	t, _ := res.Get(0).(*user.User)
	err, _ := res.Get(1).(error)
	return t, err
}
func (mc *fakeDao) FindAll() ([]*user.User, error) {
	res := mc.Called()
	t, _ := res.Get(0).([]*user.User)
	err, _ := res.Get(1).(error)
	return t, err
}
func (mc *fakeDao) FindByID(userID string) (*user.User, error) {
	res := mc.Called(userID)
	t, _ := res.Get(0).(*user.User)
	err, _ := res.Get(1).(error)
	return t, err
}
func (mc *fakeDao) FindByLogin(login string) (*user.User, error) {
	res := mc.Called(login)
	t, _ := res.Get(0).(*user.User)
	err, _ := res.Get(1).(error)
	return t, err
}
func (mc *fakeDao) Delete(userID string) error {
	res := mc.Called(userID)
	err, _ := res.Get(0).(error)
	return err
}

func (mc *fakeDao) GetID(ID string) (*objectid.ObjectID, error) {
	res := mc.Called(ID)
	t, _ := res.Get(0).(*objectid.ObjectID)
	err, _ := res.Get(1).(error)
	return t, err
}

func newFakeTokenService() token.Service {
	result := fakeTokenService{}

	payload := new(token.Payload)
	payload.TokenID = "992a6b7d893dc92de5a8b811"
	payload.UserID = "112a6b7d893dc92de5a8b899"

	tokenTxt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbklEIjoiNWIyYWZkMjA2MDFlZDljNzQ0NDVhYjU3IiwidXNlcklEIjoiNWIyYTZiN2Q4OTNkYzkyZGU1YThiODMzIn0.RBcB_B5D6uL3JXRbi2xe-V9LytIOxxLSnXv0_-rFAVU"

	result.On("Create", mock.Anything).Return(tokenTxt, nil)
	result.On("Validate", mock.Anything).Return(payload, nil)
	result.On("Invalidate", mock.Anything).Return(nil)
	result.On("extractPayload", mock.Anything).Return(payload, nil)

	return &result
}

func newFakeTokenErrorService(err error) token.Service {
	result := fakeTokenService{}

	result.On("Create", mock.Anything).Return(nil, err)
	result.On("Validate", mock.Anything).Return(nil, err)
	result.On("Invalidate", mock.Anything).Return(err)
	result.On("ExtractPayload", mock.Anything).Return(nil, err)

	return &result
}

type fakeTokenService struct {
	mock.Mock
}

func (s fakeTokenService) Create(userID objectid.ObjectID) (string, error) {
	res := s.Called(userID)
	t, _ := res.Get(0).(string)
	err, _ := res.Get(1).(error)
	return t, err
}

func (s fakeTokenService) Validate(tokenString string) (*token.Payload, error) {
	res := s.Called(tokenString)
	t, _ := res.Get(0).(*token.Payload)
	err, _ := res.Get(1).(error)
	return t, err
}
func (s fakeTokenService) Invalidate(tokenString string) error {
	res := s.Called(tokenString)
	err, _ := res.Get(0).(error)
	return err
}
func (s fakeTokenService) ExtractPayload(tokenString string) (*token.Payload, error) {
	res := s.Called(tokenString)
	t, _ := res.Get(0).(*token.Payload)
	err, _ := res.Get(1).(error)
	return t, err
}
