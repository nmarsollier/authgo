package mocks

import (
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/user"
	"github.com/stretchr/testify/mock"
)

func NewFakeDao() user.Dao {
	return NewCustomFakeDao("TestName", "Login", "Password", true)
}

func NewFakeErrorDao(err error) user.Dao {
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

func NewCustomFakeDao(name string, login string, passw string, userEnabled bool) user.Dao {
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
