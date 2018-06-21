package token

import (
	"testing"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/db"
)

func TestCreate(t *testing.T) {
	srv := token.NewTestingService(newFakeDao())

	tokenID, _ := objectid.FromHex("5b2a6b7d893dc92de5a8b833")
	token, err := srv.Create(tokenID)
	assert.NotNil(t, token)
	assert.Nil(t, err)

	payload, err := srv.ExtractPayload(token)
	assert.Nil(t, err)
	assert.NotNil(t, payload.TokenID)
	assert.Equal(t, payload.UserID, "112a6b7d893dc92de5a8b811")
}

func TestValidate(t *testing.T) {
	srv := token.NewTestingService(newFakeDao())

	token := "__invalid__"

	payload, err := srv.Validate(token)
	assert.NotNil(t, err)
	assert.Nil(t, payload)

	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbklEIjoiNWIyYWZkMjA2MDFlZDljNzQ0NDVhYjU3IiwidXNlcklEIjoiNWIyYTZiN2Q4OTNkYzkyZGU1YThiODMzIn0.RBcB_B5D6uL3JXRbi2xe-V9LytIOxxLSnXv0_-rFAVU"

	payload, err = srv.Validate(token)
	assert.Nil(t, err)
	assert.NotNil(t, payload.TokenID)
	assert.Equal(t, payload.UserID, "5b2a6b7d893dc92de5a8b833")
}

func TestInvalidate(t *testing.T) {
	srv := token.NewTestingService(newFakeDao())

	token := "__invalid__"

	err := srv.Invalidate(token)
	assert.NotNil(t, err)

	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbklEIjoiNWIyYWZkMjA2MDFlZDljNzQ0NDVhYjU3IiwidXNlcklEIjoiNWIyYTZiN2Q4OTNkYzkyZGU1YThiODMzIn0.RBcB_B5D6uL3JXRbi2xe-V9LytIOxxLSnXv0_-rFAVU"

	err = srv.Invalidate(token)
	assert.Nil(t, err)
}

func newFakeDao() token.Dao {
	result := fakeDao{}

	token := token.NewToken()
	token.ID, _ = objectid.FromHex("992a6b7d893dc92de5a8b899")
	token.UserID, _ = objectid.FromHex("112a6b7d893dc92de5a8b811")

	result.On("Collection").Return(nil, nil)
	result.On("Insert", mock.Anything).Return(token, nil)
	result.On("Update", mock.Anything).Return(token, nil)
	result.On("FindByID", mock.Anything).Return(token, nil)
	result.On("FindByUserID", mock.Anything).Return(token, nil)
	result.On("Delete", mock.Anything).Return(nil)
	result.On("GetID", mock.Anything).Return(token.ID, nil)

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

func (mc *fakeDao) Insert(tokenStr *token.Token) (*token.Token, error) {
	res := mc.Called(tokenStr)
	t, _ := res.Get(0).(*token.Token)
	err, _ := res.Get(1).(error)
	return t, err
}

func (mc *fakeDao) Update(tokenStr *token.Token) (*token.Token, error) {
	res := mc.Called(tokenStr)
	t, _ := res.Get(0).(*token.Token)
	err, _ := res.Get(1).(error)
	return t, err
}

func (mc *fakeDao) FindByID(tokenID string) (*token.Token, error) {
	res := mc.Called(tokenID)
	t, _ := res.Get(0).(*token.Token)
	err, _ := res.Get(1).(error)
	return t, err
}

func (mc *fakeDao) FindByUserID(tokenID string) (*token.Token, error) {
	res := mc.Called(tokenID)
	t, _ := res.Get(0).(*token.Token)
	err, _ := res.Get(1).(error)
	return t, err
}

func (mc *fakeDao) Delete(tokenID string) error {
	res := mc.Called(tokenID)
	err, _ := res.Get(0).(error)
	return err
}

func (mc *fakeDao) GetID(ID string) (*objectid.ObjectID, error) {
	res := mc.Called(ID)
	t, _ := res.Get(0).(*objectid.ObjectID)
	err, _ := res.Get(1).(error)
	return t, err
}
