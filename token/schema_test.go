package token

import (
	"testing"

	"github.com/mongodb/mongo-go-driver/bson/objectid"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	token := newToken()
	assert.Equal(t, token.Enabled, true)
	assert.NotEqual(t, token.ID.Hex(), "000000000000000000000000")
	assert.Equal(t, token.UserID.Hex(), "000000000000000000000000")

	tokenID, err := objectid.FromHex("5b2a6b7d893dc92de5a8b833")
	token.ID = tokenID
	assert.Nil(t, err)
	doc, err := bson.NewDocumentEncoder().EncodeDocument(token)
	assert.Nil(t, err)
	jsonDoc := doc.ToExtJSON(false)
	assert.Equal(t, "{\"_id\":{\"$oid\":\"5b2a6b7d893dc92de5a8b833\"},\"userId\":{\"$oid\":\"000000000000000000000000\"},\"enabled\":true}", jsonDoc)
}
