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
)

func TestPostUserPasswordHappyPath(t *testing.T) {
	userData, _ := user.TestUser()
	tokenData, tokenString := token.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	token.ExpectTokenAuthFindOne(t, mongodb, tokenData)

	mongodb.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter user.DbUserIdFilter, updated *user.User) error {
			// Check parameters
			assert.Equal(t, tokenData.UserID, filter.ID)

			// Asign return values
			*updated = *userData
			return nil
		},
	).Times(1)

	mongodb.EXPECT().UpdateOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter user.DbUserIdFilter, update user.DbUserUpdateDocument) (int64, error) {
			// Check parameters
			assert.Equal(t, userData.ID, filter.ID)

			assert.Equal(t, true, update.Set.Enabled)
			assert.Equal(t, "Name", update.Set.Name)
			assert.NotEmpty(t, update.Set.Password)
			assert.Contains(t, update.Set.Permissions, "user")
			assert.Contains(t, update.Set.Permissions, "other")

			// Asign return values
			return 1, nil
		},
	).Times(1)

	// REQUEST
	r := server.TestRouter(mongodb, log.NewTestLogger(ctrl, 6, 0, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestPostRequest("/v1/user/password", changePasswordBody{Current: "123", New: "456"}, tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostUserPasswordMissingCurrent(t *testing.T) {
	tokenData, tokenString := token.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	token.ExpectTokenAuthFindOne(t, mongodb, tokenData)

	// REQUEST
	r := server.TestRouter(mongodb, log.NewTestLogger(ctrl, 6, 0, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestPostRequest("/v1/user/password", changePasswordBody{New: "456"}, tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	result := w.Body.String()
	assert.Contains(t, result, "current")
	assert.Contains(t, result, "required")
}

func TestPostUserPasswordMissingNew(t *testing.T) {
	tokenData, tokenString := token.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	token.ExpectTokenAuthFindOne(t, mongodb, tokenData)

	// REQUEST
	r := server.TestRouter(mongodb, log.NewTestLogger(ctrl, 6, 0, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestPostRequest("/v1/user/password", changePasswordBody{Current: "123"}, tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	result := w.Body.String()

	assert.Contains(t, result, "new")
	assert.Contains(t, result, "required")
}

func TestPostUserPasswordWrongCurrent(t *testing.T) {
	userData, _ := user.TestUser()
	tokenData, tokenString := token.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	token.ExpectTokenAuthFindOne(t, mongodb, tokenData)

	mongodb.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter user.DbUserIdFilter, updated *user.User) error {
			// Check parameters
			assert.Equal(t, tokenData.UserID, filter.ID)

			// Asign return values
			*updated = *userData
			return nil
		},
	).Times(1)

	// REQUEST
	r := server.TestRouter(mongodb, log.NewTestLogger(ctrl, 6, 0, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestPostRequest("/v1/user/password", changePasswordBody{Current: "456", New: "456"}, tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	assert.NotEmpty(t, result["error"])
	assert.Contains(t, result["error"], "password")
	assert.Contains(t, result["error"], "invalid")
}

func TestPostUserPasswordUserNotFound(t *testing.T) {
	tokenData, tokenString := token.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	token.ExpectTokenAuthFindOne(t, mongodb, tokenData)

	db.ExpectFindOneError(mongodb, errs.NotFound, 1)

	// REQUEST
	r := server.TestRouter(mongodb, log.NewTestLogger(ctrl, 6, 1, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestPostRequest("/v1/user/password", changePasswordBody{Current: "123", New: "456"}, tokenString)
	r.ServeHTTP(w, req)

	server.AssertDocumentNotFound(t, w)
}

func TestPostUserPasswordUpdateFails(t *testing.T) {
	userData, _ := user.TestUser()
	tokenData, tokenString := token.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	token.ExpectTokenAuthFindOne(t, mongodb, tokenData)

	mongodb.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter user.DbUserIdFilter, updated *user.User) error {
			// Check parameters
			assert.Equal(t, tokenData.UserID, filter.ID)

			// Asign return values
			*updated = *userData
			return nil
		},
	).Times(1)

	db.ExpectUpdateOneError(mongodb, user.ErrID, 1)

	// REQUEST
	r := server.TestRouter(mongodb, log.NewTestLogger(ctrl, 6, 1, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestPostRequest("/v1/user/password", changePasswordBody{Current: "123", New: "456"}, tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	assert.NotEmpty(t, result["error"])
	assert.Contains(t, result["error"], "id")
	assert.Contains(t, result["error"], "Invalid")
}
