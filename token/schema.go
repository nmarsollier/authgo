package token

import (
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// Token data structure
type Token struct {
	ID      objectid.ObjectID `bson:"_id"`
	UserID  objectid.ObjectID `bson:"userId"`
	Enabled bool              `bson:"enabled"`
}

// New creates new User
func newToken() *Token {
	return &Token{
		ID:      objectid.New(),
		Enabled: true,
	}
}
