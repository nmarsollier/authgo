package db

import (
	"testing"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/nmarsollier/authgo/tools/db"
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

	unique := db.IsUniqueKeyError(MongoUnique)
	assert.Equal(t, unique, true)

	notunique := db.IsUniqueKeyError(MongoError)
	assert.Equal(t, notunique, false)
}
