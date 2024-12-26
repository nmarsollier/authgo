package mock

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/internal/token"
	"github.com/nmarsollier/authgo/test/mockgen"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestToken() (*token.Token, string) {
	tokenData := &token.Token{
		ID:      primitive.NewObjectID(),
		UserID:  primitive.NewObjectID(),
		Enabled: true,
	}

	tokenString, _ := token.Encode(tokenData)

	return tokenData, tokenString
}

func ExpectTokenFindOne(coll *mockgen.MockCollection, tokenData *token.Token, times int) {
	coll.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, params token.DbTokenIdFilter, token *token.Token) error {
			// Asign return values
			*token = *tokenData
			return nil
		},
	).Times(times)
}

func ExpectTokenAuthFindOne(t *testing.T, coll *mockgen.MockCollection, tokenData *token.Token) {
	coll.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, tokenIdUpdate token.DbTokenIdFilter, token *token.Token) error {
			assert.Equal(t, tokenData.ID, tokenIdUpdate.ID)

			*token = *tokenData
			return nil
		},
	).Times(1)
}
