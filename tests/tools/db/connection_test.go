package db

import (
	"testing"

	"github.com/nmarsollier/authgo/tools/db"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestIsUniqueKey(t *testing.T) {

	MongoUnique := mongo.WriteException{
		WriteErrors: []mongo.WriteError{
			{
				Index:   1,
				Code:    11000,
				Message: "Index",
			},
		},
	}
	MongoError := mongo.WriteException{
		WriteErrors: []mongo.WriteError{
			{
				Index:   1,
				Code:    11001,
				Message: "Other",
			},
		},
	}

	unique := db.IsUniqueKeyError(MongoUnique)
	assert.Equal(t, unique, true)

	notunique := db.IsUniqueKeyError(MongoError)
	assert.Equal(t, notunique, false)
}
