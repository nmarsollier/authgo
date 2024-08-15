package rest

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/rest/server"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/errs"
	"github.com/nmarsollier/authgo/user"
	"github.com/stretchr/testify/assert"
)

func TestPostUserInHappyPath(t *testing.T) {
	userData, password := user.TestUser()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	db.ExpectInsertOne(mongodb, 1)

	db.ExpectInsertOne(mongodb, 1)

	// REQUEST
	r := server.TestRouter(mongodb)
	InitRoutes()

	req, w := server.TestPostRequest("/v1/user", user.SignUpRequest{Login: userData.Login, Password: password, Name: userData.Name}, "")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.NotEmpty(t, result["token"], "")
}

func TestPostUserMissingLogin(t *testing.T) {
	_, password := user.TestUser()

	// REQUEST
	r := server.TestRouter()
	InitRoutes()

	req, w := server.TestPostRequest("/v1/user", user.SignUpRequest{Password: password}, "")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	result := w.Body.String()
	assert.Contains(t, result, "login")
	assert.Contains(t, result, "name")
	assert.Contains(t, result, "required")
}

func TestPostUserInvalidLoginMinRule(t *testing.T) {
	// REQUEST
	r := server.TestRouter()
	InitRoutes()

	req, w := server.TestPostRequest("/v1/user", user.SignUpRequest{Login: "a", Name: "b", Password: "c"}, "")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	result := w.Body.String()
	assert.Contains(t, result, "login")
	assert.Contains(t, result, "min")
}

func TestPostUserIvalidPassword(t *testing.T) {
	userData, _ := user.TestUser()

	// REQUEST
	r := server.TestRouter()
	InitRoutes()

	req, w := server.TestPostRequest("/v1/user", user.SignUpRequest{Name: userData.Name, Login: userData.Login}, "")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	result := w.Body.String()
	assert.Contains(t, result, "password")
	assert.Contains(t, result, "required")
}

func TestPostUserDatabaseError(t *testing.T) {
	userData, password := user.TestUser()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	db.ExpectInsertOneError(mongodb, db.TestOtherDbError, 1)

	// REQUEST
	r := server.TestRouter(mongodb)
	InitRoutes()

	req, w := server.TestPostRequest("/v1/user", user.SignUpRequest{Login: userData.Login, Password: password, Name: userData.Name}, "")
	r.ServeHTTP(w, req)

	server.AssertInternalServerError(t, w)
}

func TestPostUserAlreayExist(t *testing.T) {
	userData, password := user.TestUser()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	db.ExpectInsertOneError(mongodb, db.TestIsUniqueError, 1)

	// REQUEST
	r := server.TestRouter(mongodb)
	InitRoutes()

	req, w := server.TestPostRequest("/v1/user", user.SignUpRequest{Login: userData.Login, Password: password, Name: userData.Name}, "")
	r.ServeHTTP(w, req)

	server.AssertBadRequestError(t, w)
}

func TestPostTokenDatabaseError(t *testing.T) {
	userData, password := user.TestUser()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	db.ExpectInsertOne(mongodb, 1)

	db.ExpectInsertOneError(mongodb, errs.Internal, 1)

	// REQUEST
	r := server.TestRouter(mongodb)
	InitRoutes()

	req, w := server.TestPostRequest("/v1/user", user.SignUpRequest{Login: userData.Login, Password: password, Name: userData.Name}, "")
	r.ServeHTTP(w, req)

	server.AssertInternalServerError(t, w)
}
