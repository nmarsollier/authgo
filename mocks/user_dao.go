package mocks

import (
	"github.com/nmarsollier/authgo/user"
	"github.com/stretchr/testify/mock"
)

func MockedUserDao() user.Dao {
	return MockedUserDaoCustom("TestName", "Login", "Password", true)
}

func MockedUserDaoError(err error) user.Dao {
	result := fakeDao{}

	result.On("Insert", mock.Anything).Return(nil, err)
	result.On("Update", mock.Anything).Return(nil, err)
	result.On("FindAll").Return(nil, err)
	result.On("FindByID", mock.Anything).Return(nil, err)
	result.On("FindByLogin", mock.Anything).Return(nil, err)

	return &result
}

func MockedUserDaoCustom(name string, login string, passw string, userEnabled bool) user.Dao {
	result := fakeDao{}

	usr := user.NewUser()
	usr.Name = name
	usr.Login = login
	usr.Enabled = userEnabled
	usr.SetPasswordText(passw)
	users := make([]*user.User, 1)
	users[0] = usr

	result.On("Insert", mock.Anything).Return(usr, nil)
	result.On("Update", mock.Anything).Return(usr, nil)
	result.On("FindAll").Return(usr, nil)
	result.On("FindByID", mock.Anything).Return(usr, nil)
	result.On("FindByLogin", mock.Anything).Return(usr, nil)

	return &result
}

type fakeDao struct {
	mock.Mock
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
