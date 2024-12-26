package rest

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/internal/engine/errs"
	"github.com/nmarsollier/authgo/internal/rest"
	"github.com/nmarsollier/authgo/internal/user"
	"github.com/nmarsollier/authgo/test/engine/db"
	"github.com/nmarsollier/authgo/test/engine/di"
	"github.com/nmarsollier/authgo/test/mock"
	"github.com/nmarsollier/authgo/test/mockgen"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
)

func TestGetUserCurrentHappyPath(t *testing.T) {
	userData, _ := mock.TestUser()
	tokenData, tokenString := mock.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := mockgen.NewMockCollection(ctrl)

	mock.ExpectTokenAuthFindOne(t, mongo, tokenData)
	mongo.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter user.DbUserIdFilter, updated *user.User) error {
			assert.Equal(t, tokenData.UserID, filter.ID)

			*updated = *userData
			return nil
		},
	).Times(1)

	// REQUEST
	deps := di.NewTestInjector(ctrl, 2, 0, 1, 0, 0, 0)
	deps.SetTokenCollection(mongo)
	deps.SetUserCollection(mongo)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := TestGetRequest("/users/current", tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "Login", result["login"])
	assert.Equal(t, "Name", result["name"])
	assert.NotEmpty(t, result["id"])
}

func TestGetUserCurrentErrorDisabledToken(t *testing.T) {
	tokenData, tokenString := mock.TestToken()
	tokenData.Enabled = false

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := mockgen.NewMockCollection(ctrl)

	mock.ExpectTokenAuthFindOne(t, mongo, tokenData)

	// REQUEST
	deps := di.NewTestInjector(ctrl, 1, 1, 1, 0, 0, 0)
	deps.SetTokenCollection(mongo)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := TestGetRequest("/users/current", tokenString)
	r.ServeHTTP(w, req)

	AssertUnauthorized(t, w)
}

func TestGetUserCurrentErrorDisabledUser(t *testing.T) {
	userData, _ := mock.TestUser()
	userData.Enabled = false
	tokenData, tokenString := mock.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := mockgen.NewMockCollection(ctrl)
	mock.ExpectTokenAuthFindOne(t, mongo, tokenData)
	mock.ExpectUserFindOne(mongo, userData, 1)

	// REQUEST
	deps := di.NewTestInjector(ctrl, 2, 0, 1, 0, 0, 0)
	deps.SetTokenCollection(mongo)
	deps.SetUserCollection(mongo)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := TestGetRequest("/users/current", tokenString)
	r.ServeHTTP(w, req)

	AssertDocumentNotFound(t, w)
}

func TestGetUserCurrentErrorTokenNotFound(t *testing.T) {
	_, tokenString := mock.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := mockgen.NewMockCollection(ctrl)
	mock.ExpectFindOneError(mongo, errs.Internal, 1)

	// REQUEST
	deps := di.NewTestInjector(ctrl, 1, 2, 1, 0, 0, 0)
	deps.SetUserCollection(mongo)
	deps.SetTokenCollection(mongo)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := TestGetRequest("/users/current", tokenString)
	r.ServeHTTP(w, req)

	AssertUnauthorized(t, w)
}

func TestGetUserCurrentErrorUserNotFound(t *testing.T) {
	tokenData, tokenString := mock.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := mockgen.NewMockCollection(ctrl)

	mock.ExpectTokenFindOne(mongo, tokenData, 1)

	db.ExpectFindOneError(mongo, topology.ErrServerSelectionTimeout, 1)

	// REQUEST
	deps := di.NewTestInjector(ctrl, 2, 1, 1, 0, 0, 0)
	deps.SetUserCollection(mongo)
	deps.SetTokenCollection(mongo)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := TestGetRequest("/users/current", tokenString)
	r.ServeHTTP(w, req)

	AssertInternalServerError(t, w)
}
