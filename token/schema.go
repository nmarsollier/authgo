package token

import (
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/nmarsollier/authgo/tools/lookup"
)

// Token data structure
type Token struct {
	_id     string
	UserID  string `bson:"userId"`
	Enabled bool   `bson:"enabled"`
}

// New creates new User
func newToken() *Token {
	return &Token{Enabled: true}
}

// SetID sets the user id based on ObjectID
func (e *Token) SetID(ID objectid.ObjectID) {
	e._id = ID.Hex()
}

// ID get the token ID
func (e *Token) ID() string {
	return e._id
}

// NewFromBson creates new User form bson.Document
func newTokenFromBson(document bson.Document) *Token {
	return &Token{
		_id:     lookup.ObjectID(document, "_id"),
		UserID:  lookup.String(document, "userId"),
		Enabled: lookup.Bool(document, "enabled"),
	}
}
