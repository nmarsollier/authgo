package rest

import (
	"encoding/json"
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

func TestGetUserCurrentHappyPath(t *testing.T) {
	userData, _ := tests.TestUser()
	tokenData, tokenString := tests.TestToken()

	// User Dao Mocks
	ctrl := gomock.NewController(t)
	userCollection := db.NewMockMongoCollection(ctrl)
	userCollection.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, params primitive.M, updated *user.User) error {
			assert.Equal(t, tokenData.UserID, params["_id"])

			*updated = *userData
			return nil
		},
	).Times(1)

	// Token Dao Mocks
	tokenCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectFindOneForToken(t, tokenCollection, tokenData)

	// REQUEST
	r := engine.TestRouter(token.NewTokenOption(tokenCollection), user.NewOptions(userCollection))
	InitRoutes()

	req, w := tests.TestGetRequest("/v1/users/current", tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "Login", result["login"])
	assert.Equal(t, "Name", result["name"])
	assert.NotEmpty(t, result["id"])
}

func TestGetUserCurrentErrorDisabledToken(t *testing.T) {
	tokenData, tokenString := tests.TestToken()
	tokenData.Enabled = false

	// Dao Mocks
	ctrl := gomock.NewController(t)
	tokenCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectTokenFinOne(tokenCollection, tokenData, 1)

	// REQUEST
	r := engine.TestRouter(token.NewTokenOption(tokenCollection))
	InitRoutes()

	req, w := tests.TestGetRequest("/v1/users/current", tokenString)
	r.ServeHTTP(w, req)

	tests.AssertUnauthorized(t, w)
}

func TestGetUserCurrentErrorDisabledUser(t *testing.T) {
	userData, _ := tests.TestUser()
	userData.Enabled = false
	tokenData, tokenString := tests.TestToken()

	// User Dao Mocks
	ctrl := gomock.NewController(t)
	userCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectUserFindOne(userCollection, userData, 1)

	// Token Dao Mocks
	tokenCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectTokenFinOne(tokenCollection, tokenData, 1)

	// REQUEST
	r := engine.TestRouter(token.NewTokenOption(tokenCollection), user.NewOptions(userCollection))
	InitRoutes()

	req, w := tests.TestGetRequest("/v1/users/current", tokenString)
	r.ServeHTTP(w, req)

	tests.AssertDocumentNotFound(t, w)
}

func TestGetUserCurrentErrorTokenNotFound(t *testing.T) {
	_, tokenString := tests.TestToken()

	// Dao Mocks
	ctrl := gomock.NewController(t)
	tokenCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectFindOneError(tokenCollection, app_errors.Internal, 1)

	// REQUEST
	r := engine.TestRouter(token.NewTokenOption(tokenCollection))
	InitRoutes()

	req, w := tests.TestGetRequest("/v1/users/current", tokenString)
	r.ServeHTTP(w, req)

	tests.AssertUnauthorized(t, w)
}

func TestGetUserCurrentErrorUserNotFound(t *testing.T) {
	tokenData, tokenString := tests.TestToken()

	// User Dao Mocks
	ctrl := gomock.NewController(t)
	userCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectFindOneError(userCollection, app_errors.Internal, 1)

	// Token Dao Mocks
	tokenCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectTokenFinOne(tokenCollection, tokenData, 1)

	// REQUEST
	r := engine.TestRouter(token.NewTokenOption(tokenCollection), user.NewOptions(userCollection))
	InitRoutes()

	req, w := tests.TestGetRequest("/v1/users/current", tokenString)
	r.ServeHTTP(w, req)

	tests.AssertInternalServerError(t, w)
}
