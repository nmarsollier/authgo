package mocks

import (
	"github.com/nmarsollier/authgo/security"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MockedTokenDaoError(err error) security.Dao {
	result := fakeTokenDao{}

	result.On("Create", mock.Anything).Return(nil, err)
	result.On("FindByID", mock.Anything).Return(nil, err)
	result.On("Delete").Return(err)

	return &result
}

func MockedTokenDao() security.Dao {
	result := fakeTokenDao{}

	tokenID, _ := primitive.ObjectIDFromHex("111111111111111111111111")
	userID, _ := primitive.ObjectIDFromHex("999999999999999999999999")
	token := &security.Token{
		ID:      tokenID,
		UserID:  userID,
		Enabled: true,
	}

	result.On("Create", mock.Anything).Return(token, nil)
	result.On("FindByID", mock.Anything).Return(token, nil)
	result.On("Delete").Return(nil)

	return &result
}

type fakeTokenDao struct {
	mock.Mock
}

func (mc *fakeTokenDao) Create(userID primitive.ObjectID) (*security.Token, error) {
	res := mc.Called(userID)
	t, _ := res.Get(0).(*security.Token)
	err, _ := res.Get(1).(error)
	return t, err
}
func (mc *fakeTokenDao) FindByID(tokenID string) (*security.Token, error) {
	res := mc.Called(tokenID)
	t, _ := res.Get(0).(*security.Token)
	err, _ := res.Get(1).(error)
	return t, err
}

func (mc *fakeTokenDao) Delete(tokenID primitive.ObjectID) error {
	res := mc.Called()
	err, _ := res.Get(0).(error)
	return err
}
