package db

import (
	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/test/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mock Data
var TestIsUniqueError = mongo.WriteException{
	WriteErrors: []mongo.WriteError{
		{
			Code: 11000,
		},
	},
}

var TestOtherDbError = mongo.WriteException{
	WriteErrors: []mongo.WriteError{
		{
			Code: 1,
		},
	},
}

// Espect common functions
func ExpectFindOneError(coll *mock.MockCollection, err error, times int) {
	coll.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, params interface{}, update interface{}) error {
			return err
		},
	).Times(times)
}
func ExpectInsertOneError(coll *mock.MockCollection, err error, times int) {
	coll.EXPECT().InsertOne(gomock.Any(), gomock.Any()).Return("", err).Times(times)
}

func ExpectUpdateOneError(coll *mock.MockCollection, err error, times int) {
	coll.EXPECT().UpdateOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), err).Times(times)
}

func ExpectInsertOne(coll *mock.MockCollection, times int) {
	coll.EXPECT().InsertOne(gomock.Any(), gomock.Any()).Return("", nil).Times(times)
}
