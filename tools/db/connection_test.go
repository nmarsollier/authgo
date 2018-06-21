package db

import (
	"testing"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/stretchr/testify/assert"
)

func TestIsUniqueKey(t *testing.T) {

	var MongoUnique = mongo.WriteErrors{
		mongo.WriteError{
			Index:   1,
			Code:    11000,
			Message: "Index",
		},
	}
	var MongoError = mongo.WriteErrors{
		mongo.WriteError{
			Index:   1,
			Code:    11001,
			Message: "Other",
		},
	}

	unique := IsUniqueKeyError(MongoUnique)
	assert.Equal(t, unique, true)

	notunique := IsUniqueKeyError(MongoError)
	assert.Equal(t, notunique, false)
}
