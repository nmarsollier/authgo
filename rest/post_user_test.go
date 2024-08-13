package rest

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/rest/engine"
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/apperr"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/tests"
	"github.com/nmarsollier/authgo/user"
	"github.com/stretchr/testify/assert"
)

func TestPostUserInHappyPath(t *testing.T) {
	userData, password := tests.TestUser()

	// User Dao Mocks
	ctrl := gomock.NewController(t)
	userCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectUserInsertOne(userCollection, 1)

	// Token Dao Mocks
	tokenCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectTokenInsertOne(tokenCollection, 1)

	// REQUEST
	r := engine.TestRouter(token.TokenCollection(tokenCollection), user.UserCollection(userCollection))
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/user", user.SignUpRequest{Login: userData.Login, Password: password, Name: userData.Name}, "")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.NotEmpty(t, result["token"], "")
}

func TestPostUserMissingLogin(t *testing.T) {
	_, password := tests.TestUser()

	// REQUEST
	r := engine.TestRouter()
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/user", user.SignUpRequest{Password: password}, "")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	result := w.Body.String()
	assert.Contains(t, result, "login")
	assert.Contains(t, result, "name")
	assert.Contains(t, result, "required")
}

func TestPostUserInvalidLoginMinRule(t *testing.T) {
	// REQUEST
	r := engine.TestRouter()
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/user", user.SignUpRequest{Login: "a", Name: "b", Password: "c"}, "")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	result := w.Body.String()
	assert.Contains(t, result, "login")
	assert.Contains(t, result, "min")
}

func TestPostUserIvalidPassword(t *testing.T) {
	userData, _ := tests.TestUser()

	// REQUEST
	r := engine.TestRouter()
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/user", user.SignUpRequest{Name: userData.Name, Login: userData.Login}, "")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	result := w.Body.String()
	assert.Contains(t, result, "password")
	assert.Contains(t, result, "required")
}

func TestPostUserDatabaseError(t *testing.T) {
	userData, password := tests.TestUser()

	// User Dao Mocks
	ctrl := gomock.NewController(t)
	userCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectInsertOneError(userCollection, tests.TestOtherDbError, 1)

	// REQUEST
	r := engine.TestRouter(user.UserCollection(userCollection))
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/user", user.SignUpRequest{Login: userData.Login, Password: password, Name: userData.Name}, "")
	r.ServeHTTP(w, req)

	tests.AssertInternalServerError(t, w)
}

func TestPostUserAlreayExist(t *testing.T) {
	userData, password := tests.TestUser()

	// User Dao Mocks
	ctrl := gomock.NewController(t)
	userCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectInsertOneError(userCollection, tests.TestIsUniqueError, 1)

	// REQUEST
	r := engine.TestRouter(user.UserCollection(userCollection))
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/user", user.SignUpRequest{Login: userData.Login, Password: password, Name: userData.Name}, "")
	r.ServeHTTP(w, req)

	tests.AssertBadRequestError(t, w)
}

func TestPostTokenDatabaseError(t *testing.T) {
	userData, password := tests.TestUser()

	// User Dao Mocks
	ctrl := gomock.NewController(t)
	userCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectUserInsertOne(userCollection, 1)

	tokenCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectInsertOneError(tokenCollection, apperr.Internal, 1)

	// REQUEST
	r := engine.TestRouter(user.UserCollection(userCollection), token.TokenCollection(tokenCollection))
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/user", user.SignUpRequest{Login: userData.Login, Password: password, Name: userData.Name}, "")
	r.ServeHTTP(w, req)

	tests.AssertInternalServerError(t, w)
}
