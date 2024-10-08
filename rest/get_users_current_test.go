package rest

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/rest/server"
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/errs"
	"github.com/nmarsollier/authgo/tools/log"
	"github.com/nmarsollier/authgo/user"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
)

func TestGetUserCurrentHappyPath(t *testing.T) {
	userData, _ := user.TestUser()
	tokenData, tokenString := token.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := db.NewMockMongoCollection(ctrl)

	token.ExpectTokenAuthFindOne(t, mongo, tokenData)

	mongo.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter user.DbUserIdFilter, updated *user.User) error {
			assert.Equal(t, tokenData.UserID, filter.ID)

			*updated = *userData
			return nil
		},
	).Times(1)

	// REQUEST
	r := server.TestRouter(mongo, log.NewTestLogger(ctrl, 6, 0, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestGetRequest("/v1/users/current", tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "Login", result["login"])
	assert.Equal(t, "Name", result["name"])
	assert.NotEmpty(t, result["id"])
}

func TestGetUserCurrentErrorDisabledToken(t *testing.T) {
	tokenData, tokenString := token.TestToken()
	tokenData.Enabled = false

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := db.NewMockMongoCollection(ctrl)

	token.ExpectTokenFindOne(mongo, tokenData, 1)

	// REQUEST
	r := server.TestRouter(mongo, log.NewTestLogger(ctrl, 5, 1, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestGetRequest("/v1/users/current", tokenString)
	r.ServeHTTP(w, req)

	server.AssertUnauthorized(t, w)
}

func TestGetUserCurrentErrorDisabledUser(t *testing.T) {
	userData, _ := user.TestUser()
	userData.Enabled = false
	tokenData, tokenString := token.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := db.NewMockMongoCollection(ctrl)

	token.ExpectTokenFindOne(mongo, tokenData, 1)

	user.ExpectUserFindOne(mongo, userData, 1)

	// REQUEST
	r := server.TestRouter(mongo, log.NewTestLogger(ctrl, 6, 0, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestGetRequest("/v1/users/current", tokenString)
	r.ServeHTTP(w, req)

	server.AssertDocumentNotFound(t, w)
}

func TestGetUserCurrentErrorTokenNotFound(t *testing.T) {
	_, tokenString := token.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := db.NewMockMongoCollection(ctrl)

	db.ExpectFindOneError(mongo, errs.Internal, 1)

	// REQUEST
	r := server.TestRouter(mongo, log.NewTestLogger(ctrl, 5, 2, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestGetRequest("/v1/users/current", tokenString)
	r.ServeHTTP(w, req)

	server.AssertUnauthorized(t, w)
}

func TestGetUserCurrentErrorUserNotFound(t *testing.T) {
	tokenData, tokenString := token.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := db.NewMockMongoCollection(ctrl)

	token.ExpectTokenFindOne(mongo, tokenData, 1)

	db.ExpectFindOneError(mongo, topology.ErrServerSelectionTimeout, 1)

	// REQUEST
	r := server.TestRouter(mongo, log.NewTestLogger(ctrl, 6, 1, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestGetRequest("/v1/users/current", tokenString)
	r.ServeHTTP(w, req)

	server.AssertInternalServerError(t, w)
}
