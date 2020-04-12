package mocks

import (
	"github.com/nmarsollier/authgo/security"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MockedSecurityService() security.Service {
	result := fakeTokenService{}

	tokenID, _ := primitive.ObjectIDFromHex("111111111111111111111111")
	userID, _ := primitive.ObjectIDFromHex("999999999999999999999999")
	token := &security.Token{
		ID:      tokenID,
		UserID:  userID,
		Enabled: true,
	}

	result.On("Create", mock.Anything).Return(token, nil)
	result.On("Find", mock.Anything).Return(token, nil)
	result.On("Validate", mock.Anything).Return(token, nil)
	result.On("Invalidate", mock.Anything).Return(nil)

	return &result
}

func MockedSecurityServiceError(err error) security.Service {
	result := fakeTokenService{}

	result.On("Create", mock.Anything).Return(nil, err)
	result.On("Find", mock.Anything).Return(nil, err)
	result.On("Validate", mock.Anything).Return(err)
	result.On("Invalidate", mock.Anything).Return(nil, err)

	return &result
}

type fakeTokenService struct {
	mock.Mock
}

func (s fakeTokenService) Create(userID primitive.ObjectID) (*security.Token, error) {
	res := s.Called(userID)
	t, _ := res.Get(0).(*security.Token)
	err, _ := res.Get(1).(error)
	return t, err
}

func (s fakeTokenService) Find(tokenID string) (*security.Token, error) {
	res := s.Called(tokenID)
	t, _ := res.Get(0).(*security.Token)
	err, _ := res.Get(1).(error)
	return t, err
}

func (s fakeTokenService) Validate(tokenString string) (*security.Token, error) {
	res := s.Called(tokenString)
	t, _ := res.Get(0).(*security.Token)
	err, _ := res.Get(1).(error)
	return t, err
}

func (s fakeTokenService) Invalidate(tokenString string) error {
	res := s.Called(tokenString)
	err, _ := res.Get(0).(error)
	return err
}
