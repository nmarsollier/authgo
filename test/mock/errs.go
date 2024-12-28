package mock

import "go.mongodb.org/mongo-driver/mongo"

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
