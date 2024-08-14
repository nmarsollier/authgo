package rest

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/rest/engine"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/errs"
	"github.com/nmarsollier/authgo/tools/tests"
	"github.com/nmarsollier/authgo/user"
	"github.com/stretchr/testify/assert"
)

func TestPostUserInHappyPath(t *testing.T) {
	userData, password := tests.TestUser()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	tests.ExpectTokenInsertOne(mongodb, 1)

	tests.ExpectUserInsertOne(mongodb, 1)

	// REQUEST
	r := engine.TestRouter(mongodb)
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

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	tests.ExpectInsertOneError(mongodb, tests.TestOtherDbError, 1)

	// REQUEST
	r := engine.TestRouter(mongodb)
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/user", user.SignUpRequest{Login: userData.Login, Password: password, Name: userData.Name}, "")
	r.ServeHTTP(w, req)

	tests.AssertInternalServerError(t, w)
}

func TestPostUserAlreayExist(t *testing.T) {
	userData, password := tests.TestUser()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	tests.ExpectInsertOneError(mongodb, tests.TestIsUniqueError, 1)

	// REQUEST
	r := engine.TestRouter(mongodb)
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/user", user.SignUpRequest{Login: userData.Login, Password: password, Name: userData.Name}, "")
	r.ServeHTTP(w, req)

	tests.AssertBadRequestError(t, w)
}

func TestPostTokenDatabaseError(t *testing.T) {
	userData, password := tests.TestUser()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	tests.ExpectUserInsertOne(mongodb, 1)

	tests.ExpectInsertOneError(mongodb, errs.Internal, 1)

	// REQUEST
	r := engine.TestRouter(mongodb)
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/user", user.SignUpRequest{Login: userData.Login, Password: password, Name: userData.Name}, "")
	r.ServeHTTP(w, req)

	tests.AssertInternalServerError(t, w)
}
