package token

import (
	"testing"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/nmarsollier/authgo/token"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	tokenTest := token.NewToken()
	assert.Equal(t, tokenTest.Enabled, true)
	assert.NotEqual(t, tokenTest.ID.Hex(), "000000000000000000000000")
	assert.Equal(t, tokenTest.UserID.Hex(), "000000000000000000000000")

	tokenID, err := objectid.FromHex("5b2a6b7d893dc92de5a8b833")
	tokenTest.ID = tokenID
	assert.Nil(t, err)
	doc, err := bson.NewDocumentEncoder().EncodeDocument(tokenTest)
	assert.Nil(t, err)
	jsonDoc := doc.ToExtJSON(false)
	assert.Equal(t, "{\"_id\":{\"$oid\":\"5b2a6b7d893dc92de5a8b833\"},\"userId\":{\"$oid\":\"000000000000000000000000\"},\"enabled\":true}", jsonDoc)
}
