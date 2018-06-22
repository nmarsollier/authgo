package test

import (
	"testing"

	"github.com/mongodb/mongo-go-driver/bson/objectid"

	"github.com/nmarsollier/authgo/mocks"
	"github.com/nmarsollier/authgo/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDaoCreate(t *testing.T) {
	mConn := new(mocks.Collection)
	dao := security.MockedDao(mConn)

	ID, _ := objectid.FromHex("5b2a6b7d893dc92de5a8b833")

	mConn.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).Return(ID, nil)

	token, err := dao.Create(ID)

	assert.Nil(t, err)
	assert.NotNil(t, token)
}

func TestDaoFind(t *testing.T) {
	mConn := new(mocks.Collection)
	dao := security.MockedDao(mConn)

	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		mocks.Decoder(func(v interface{}) error {
			return nil
		}),
	)

	token, err := dao.FindByID("5b2a6b7d893dc92de5a8b833")
	assert.Nil(t, err)
	assert.NotNil(t, token)
}
