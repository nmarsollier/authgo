package rest

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/rest/engine"
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/app_errors"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/tests"
	"github.com/nmarsollier/authgo/user"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetUserSignOutHappyPath(t *testing.T) {
	tokenData, tokenString := tests.TestToken()

	// User Dao Mocks
	ctrl := gomock.NewController(t)
	userCollection := db.NewMockMongoCollection(ctrl)

	// Token Dao Mocks
	tokenCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectFindOneForToken(t, tokenCollection, tokenData)

	// Database Mocks
	tokenCollection.EXPECT().UpdateOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter primitive.M, update primitive.M) (int64, error) {
			assert.Equal(t, tokenData.ID, filter["_id"].(primitive.ObjectID))

			assert.Equal(t, false, update["$set"].(primitive.M)["enabled"])

			return 1, nil
		},
	).Times(1)

	// REQUEST
	r := engine.TestRouter(token.NewProps(tokenCollection), user.NewProps(userCollection))
	InitRoutes()

	req, w := tests.TestGetRequest("/v1/user/signout", tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUserSignOutDbUpdateError(t *testing.T) {
	tokenData, tokenString := tests.TestToken()

	// User Dao Mocks
	ctrl := gomock.NewController(t)
	userCollection := db.NewMockMongoCollection(ctrl)

	// Token Dao Mocks
	tokenCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectFindOneForToken(t, tokenCollection, tokenData)

	// Database Mocks
	tests.ExpectUpdateOneError(tokenCollection, app_errors.NotFound, 1)

	// REQUEST
	r := engine.TestRouter(token.NewProps(tokenCollection), user.NewProps(userCollection))
	InitRoutes()

	req, w := tests.TestGetRequest("/v1/user/signout", tokenString)
	r.ServeHTTP(w, req)

	tests.AssertDocumentNotFound(t, w)
}

func TestGetUserSignOutInvalidToken(t *testing.T) {

	// REQUEST
	r := engine.TestRouter()
	InitRoutes()

	req, w := tests.TestGetRequest("/v1/user/signout", "123")
	r.ServeHTTP(w, req)

	tests.AssertUnauthorized(t, w)
}

func TestGetUserSignOutDbFindError(t *testing.T) {
	_, tokenString := tests.TestToken()

	// User Dao Mocks
	ctrl := gomock.NewController(t)
	userCollection := db.NewMockMongoCollection(ctrl)

	// Token Dao Mocks
	tokenCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectFindOneError(tokenCollection, app_errors.NotFound, 1)

	// REQUEST
	r := engine.TestRouter(token.NewProps(tokenCollection), user.NewProps(userCollection))
	InitRoutes()

	req, w := tests.TestGetRequest("/v1/user/signout", tokenString)
	r.ServeHTTP(w, req)

	tests.AssertUnauthorized(t, w)
}
