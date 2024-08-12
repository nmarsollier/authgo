package rest

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/rest/engine"
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/tests"
	"github.com/nmarsollier/authgo/user"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestGetUsersHappyPath(t *testing.T) {
	userData, _ := tests.TestAdminUser()
	tokenData, tokenString := tests.TestToken()

	// Token Dao Mocks
	ctrl := gomock.NewController(t)
	tokenCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectFindOneForToken(t, tokenCollection, tokenData)

	// User Dao Mocks
	userCollection := db.NewMockMongoCollection(ctrl)
	userCollection.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, params primitive.M, updated *user.User) error {
			// Check parameters
			assert.Equal(t, tokenData.UserID, params["_id"])

			// Asign return values
			*updated = *userData
			return nil
		},
	).Times(1)

	userCollection.EXPECT().Find(gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter primitive.D) (db.Cursor, error) {
			data := db.NewMockCursor(ctrl)
			data.EXPECT().Next(gomock.Any()).Return(true).Times(2)
			data.EXPECT().Next(gomock.Any()).Return(false).Times(1)

			data.EXPECT().Decode(gomock.Any()).DoAndReturn(
				func(updated *user.User) error {
					testUser, _ := tests.TestUser()

					*updated = *testUser

					return nil
				},
			).Times(2)

			data.EXPECT().Close(gomock.Any()).Return(nil).Times(1)

			// Asign return values
			return data, nil
		},
	).Times(1)

	// REQUEST
	r := engine.TestRouter(token.NewTokenOption(tokenCollection), user.NewOptions(userCollection))
	InitRoutes()

	req, w := tests.TestGetRequest("/v1/users", tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	result := w.Body.String()
	assert.NotEmpty(t, result)
}

func TestGetUsersFindError(t *testing.T) {
	userData, _ := tests.TestAdminUser()
	tokenData, tokenString := tests.TestToken()

	// Token Dao Mocks
	ctrl := gomock.NewController(t)
	tokenCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectFindOneForToken(t, tokenCollection, tokenData)

	// User Dao Mocks
	userCollection := db.NewMockMongoCollection(ctrl)
	userCollection.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, params primitive.M, updated *user.User) error {
			// Check parameters
			assert.Equal(t, tokenData.UserID, params["_id"])

			// Asign return values
			*updated = *userData
			return nil
		},
	).Times(1)

	userCollection.EXPECT().Find(gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter primitive.D) (db.Cursor, error) {
			// Asign return values
			return nil, mongo.ErrNoDocuments
		},
	).Times(1)

	// REQUEST
	r := engine.TestRouter(token.NewTokenOption(tokenCollection), user.NewOptions(userCollection))
	InitRoutes()

	req, w := tests.TestGetRequest("/v1/users", tokenString)
	r.ServeHTTP(w, req)

	tests.AssertDocumentNotFound(t, w)
}
