package rest

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/rest/engine"
	"github.com/nmarsollier/authgo/tools/apperr"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/tests"
	"github.com/nmarsollier/authgo/user"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
)

func TestGetUserCurrentHappyPath(t *testing.T) {
	userData, _ := tests.TestUser()
	tokenData, tokenString := tests.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := db.NewMockMongoCollection(ctrl)

	tests.ExpectFindOneForToken(t, mongo, tokenData)

	mongo.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, params primitive.M, updated *user.User) error {
			assert.Equal(t, tokenData.UserID, params["_id"])

			*updated = *userData
			return nil
		},
	).Times(1)

	// REQUEST
	r := engine.TestRouter(mongo)
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

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := db.NewMockMongoCollection(ctrl)

	tests.ExpectTokenFinOne(mongo, tokenData, 1)

	// REQUEST
	r := engine.TestRouter(mongo)
	InitRoutes()

	req, w := tests.TestGetRequest("/v1/users/current", tokenString)
	r.ServeHTTP(w, req)

	tests.AssertUnauthorized(t, w)
}

func TestGetUserCurrentErrorDisabledUser(t *testing.T) {
	userData, _ := tests.TestUser()
	userData.Enabled = false
	tokenData, tokenString := tests.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := db.NewMockMongoCollection(ctrl)

	tests.ExpectTokenFinOne(mongo, tokenData, 1)

	tests.ExpectUserFindOne(mongo, userData, 1)

	// REQUEST
	r := engine.TestRouter(mongo)
	InitRoutes()

	req, w := tests.TestGetRequest("/v1/users/current", tokenString)
	r.ServeHTTP(w, req)

	tests.AssertDocumentNotFound(t, w)
}

func TestGetUserCurrentErrorTokenNotFound(t *testing.T) {
	_, tokenString := tests.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := db.NewMockMongoCollection(ctrl)

	tests.ExpectFindOneError(mongo, apperr.Internal, 1)

	// REQUEST
	r := engine.TestRouter(mongo)
	InitRoutes()

	req, w := tests.TestGetRequest("/v1/users/current", tokenString)
	r.ServeHTTP(w, req)

	tests.AssertUnauthorized(t, w)
}

func TestGetUserCurrentErrorUserNotFound(t *testing.T) {
	tokenData, tokenString := tests.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := db.NewMockMongoCollection(ctrl)

	tests.ExpectTokenFinOne(mongo, tokenData, 1)

	tests.ExpectFindOneError(mongo, topology.ErrServerSelectionTimeout, 1)

	// REQUEST
	r := engine.TestRouter(mongo)
	InitRoutes()

	req, w := tests.TestGetRequest("/v1/users/current", tokenString)
	r.ServeHTTP(w, req)

	tests.AssertInternalServerError(t, w)
}
