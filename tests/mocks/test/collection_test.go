package test

import (
	"context"
	"testing"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/nmarsollier/authgo/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testMongoUnique = mongo.WriteErrors{
	mongo.WriteError{
		Index:   1,
		Code:    11000,
		Message: "Index",
	},
}
var testMongoError = mongo.WriteErrors{
	mongo.WriteError{
		Index:   1,
		Code:    11001,
		Message: "Other",
	},
}

func TestFindOneOk(t *testing.T) {
	mConn := new(mocks.Collection)
	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		mocks.Decoder(func(v interface{}) error {
			return nil
		}),
	)
	err := mConn.FindOne(context.Background(), nil).Decode(nil)
	assert.Nil(t, err)
}
func TestFindOneError(t *testing.T) {
	mConn := new(mocks.Collection)
	mConn.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(
		mocks.Decoder(func(v interface{}) error {
			return testMongoError
		}),
	)
	err := mConn.FindOne(context.Background(), nil).Decode(nil)
	assert.Equal(t, err, testMongoError)
}

func TestInsertOk(t *testing.T) {
	mConn := new(mocks.Collection)

	ID, _ := objectid.FromHex("5b2a6b7d893dc92de5a8b833")
	mConn.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).Return(ID, nil)
	inserted, err := mConn.InsertOne(context.Background(), nil)
	assert.Nil(t, err)
	assert.NotNil(t, ID, inserted.InsertedID)
}

func TestInsertUniqueError(t *testing.T) {
	mConn := new(mocks.Collection)

	mConn.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, testMongoUnique)

	_, err := mConn.InsertOne(context.Background(), nil)
	assert.Equal(t, err, testMongoUnique)
}

func TestInsertOtherError(t *testing.T) {
	mConn := new(mocks.Collection)

	mConn.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, testMongoError)

	_, err := mConn.InsertOne(context.Background(), nil)
	assert.Equal(t, err, testMongoError)
}

func TestUpdateOk(t *testing.T) {
	mConn := new(mocks.Collection)

	mConn.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(1, 1, 1, nil)

	result, err := mConn.UpdateOne(context.Background(), nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

func TestUpdateUniqueError(t *testing.T) {
	mConn := new(mocks.Collection)

	mConn.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(1, 1, 1, testMongoUnique)

	_, err := mConn.UpdateOne(context.Background(), nil, nil)
	assert.Equal(t, err, testMongoUnique)
}

func TestUpdateOtherError(t *testing.T) {
	mConn := new(mocks.Collection)

	mConn.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(1, 1, 1, testMongoError)

	_, err := mConn.UpdateOne(context.Background(), nil, nil)
	assert.Equal(t, err, testMongoError)
}
