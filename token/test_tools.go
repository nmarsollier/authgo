package token

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestToken() (*Token, string) {
	tokenData := &Token{
		ID:      primitive.NewObjectID(),
		UserID:  primitive.NewObjectID(),
		Enabled: true,
	}

	tokenString, _ := Encode(tokenData)

	return tokenData, tokenString
}

func ExpectTokenFindOne(coll *db.MockMongoCollection, tokenData *Token, times int) {
	coll.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, params DbTokenIdFilter, token *Token) error {
			// Asign return values
			*token = *tokenData
			return nil
		},
	).Times(times)
}

func ExpectTokenAuthFindOne(t *testing.T, coll *db.MockMongoCollection, tokenData *Token) {
	coll.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, tokenIdUpdate DbTokenIdFilter, token *Token) error {
			assert.Equal(t, tokenData.ID, tokenIdUpdate.ID)

			*token = *tokenData
			return nil
		},
	).Times(1)
}
